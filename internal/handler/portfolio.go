package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"backend-portfolio/internal/repository"
	"backend-portfolio/models"

	"github.com/gin-gonic/gin"
)

// GetAllPortfolio returns every portfolio item with its media files.
func (h *Handler) GetAllPortfolio(c *gin.Context) {
	portfolios, err := h.portfolios.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch portfolio items"})
		return
	}
	c.JSON(http.StatusOK, portfolios)
}

// GetPortfolio returns a single portfolio item by ID.
func (h *Handler) GetPortfolio(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		return
	}

	portfolio, err := h.portfolios.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio item not found"})
		return
	}
	c.JSON(http.StatusOK, portfolio)
}

// CreatePortfolio creates a new portfolio with optional file uploads (multipart).
func (h *Handler) CreatePortfolio(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32 MB max
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	portfolio := models.Portfolio{
		Title:       title,
		Description: c.PostForm("description"),
		Category:    c.PostForm("category"),
		Tags:        c.PostForm("tags"),
		ProjectURL:  c.PostForm("project_url"),
	}

	files := c.Request.MultipartForm.File["media_files"]

	err := h.portfolios.WithTransaction(func(tx repository.PortfolioRepository) error {
		if err := tx.Create(&portfolio); err != nil {
			return err
		}

		if len(files) > 0 {
			mediaFiles, err := saveUploadedFiles(files, portfolio.ID)
			if err != nil {
				return fmt.Errorf("failed to save files: %w", err)
			}
			for i := range mediaFiles {
				if err := tx.CreateMedia(&mediaFiles[i]); err != nil {
					return err
				}
			}
			portfolio.MediaFiles = mediaFiles
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create portfolio"})
		return
	}
	c.JSON(http.StatusCreated, portfolio)
}

// UpdatePortfolio updates an existing portfolio (supports both JSON and multipart).
func (h *Handler) UpdatePortfolio(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		return
	}

	portfolio, err := h.portfolios.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio item not found"})
		return
	}

	contentType := c.GetHeader("Content-Type")
	isMultipart := strings.Contains(contentType, "multipart/form-data")

	if isMultipart {
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
			return
		}
		if v := c.PostForm("title"); v != "" {
			portfolio.Title = v
		}
		if v := c.PostForm("description"); v != "" {
			portfolio.Description = v
		}
		if v := c.PostForm("category"); v != "" {
			portfolio.Category = v
		}
		if v := c.PostForm("tags"); v != "" {
			portfolio.Tags = v
		}
		if v := c.PostForm("project_url"); v != "" {
			portfolio.ProjectURL = v
		}
	} else {
		var update struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Category    string `json:"category"`
			Tags        string `json:"tags"`
			ProjectURL  string `json:"project_url"`
		}
		if err := c.ShouldBindJSON(&update); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if update.Title != "" {
			portfolio.Title = update.Title
		}
		portfolio.Description = update.Description
		portfolio.Category = update.Category
		portfolio.Tags = update.Tags
		portfolio.ProjectURL = update.ProjectURL
	}

	// Handle new file uploads if present
	if isMultipart {
		files := c.Request.MultipartForm.File["media_files"]
		if len(files) > 0 {
			mediaFiles, ferr := saveUploadedFiles(files, portfolio.ID)
			if ferr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save files: %v", ferr)})
				return
			}

			err := h.portfolios.WithTransaction(func(tx repository.PortfolioRepository) error {
				for i := range mediaFiles {
					if err := tx.CreateMedia(&mediaFiles[i]); err != nil {
						return err
					}
				}
				return tx.Save(portfolio)
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update portfolio"})
				return
			}

			// Reload with new media
			portfolio, _ = h.portfolios.FindByID(id)
			c.JSON(http.StatusOK, portfolio)
			return
		}
	}

	// Regular update without new files
	if err := h.portfolios.Save(portfolio); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update portfolio"})
		return
	}
	c.JSON(http.StatusOK, portfolio)
}

// DeletePortfolio removes a portfolio and its associated media files from storage.
func (h *Handler) DeletePortfolio(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		return
	}

	portfolio, err := h.portfolios.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio item not found"})
		return
	}

	// Delete files from disk
	for _, media := range portfolio.MediaFiles {
		removeMediaFile(media.URL)
	}

	if err := h.portfolios.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete portfolio item"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Portfolio item deleted successfully"})
}

// DeletePortfolioMedia removes a specific media file from a portfolio.
func (h *Handler) DeletePortfolioMedia(c *gin.Context) {
	portfolioID, err := parseNamedIDParam(c, "portfolio_id")
	if err != nil {
		return
	}
	mediaID, err := parseNamedIDParam(c, "media_id")
	if err != nil {
		return
	}

	media, err := h.portfolios.FindMedia(portfolioID, mediaID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media file not found"})
		return
	}

	removeMediaFile(media.URL)

	if err := h.portfolios.DeleteMedia(media); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete media file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Media file deleted successfully"})
}

// ─── helpers ────────────────────────────────────────────────────────────────

// saveUploadedFiles persists multipart files to disk and returns media models.
func saveUploadedFiles(files []*multipart.FileHeader, portfolioID uint) ([]models.PortfolioMedia, error) {
	if err := os.MkdirAll("uploads/portfolio", 0755); err != nil {
		return nil, err
	}

	media := make([]models.PortfolioMedia, 0, len(files))
	for i, fh := range files {
		src, err := fh.Open()
		if err != nil {
			return nil, err
		}

		ext := filepath.Ext(fh.Filename)
		name := fmt.Sprintf("%d_%d_%s%s",
			portfolioID,
			time.Now().UnixNano(),
			strings.ReplaceAll(fh.Filename, " ", "_"),
			ext,
		)
		dst, err := os.Create(filepath.Join("uploads/portfolio", name))
		if err != nil {
			src.Close()
			return nil, err
		}
		_, copyErr := io.Copy(dst, src)
		src.Close()
		dst.Close()
		if copyErr != nil {
			return nil, copyErr
		}

		fileType := "image"
		if strings.HasPrefix(fh.Header.Get("Content-Type"), "video/") {
			fileType = "video"
		}

		media = append(media, models.PortfolioMedia{
			PortfolioID: portfolioID,
			Type:        fileType,
			URL:         "/uploads/portfolio/" + name,
			OrderIndex:  i,
		})
	}
	return media, nil
}

// removeMediaFile deletes a media file from disk (best-effort).
func removeMediaFile(url string) {
	if url != "" {
		rel := strings.TrimPrefix(url, "/uploads/")
		os.Remove(filepath.Join("uploads", rel))
	}
}
