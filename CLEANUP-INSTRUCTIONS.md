# Files to Delete - Cleanup Instructions

This document lists files that should be deleted from the repository after CSRF + Cookie migration.

## Files to Remove

```
internal/auth/cookie.go                    # Deprecated: Cookie management removed
internal/middleware/csrf.go                # Deprecated: CSRF protection removed
internal/middleware/csrf_deprecated.go     # Temporary placeholder file
```

## How to Delete

### Option 1: Using Git (Recommended)
```bash
cd "d:\--- Project ---\-- Golang --\backend-portfolio"
git rm internal/auth/cookie.go internal/middleware/csrf.go internal/middleware/csrf_deprecated.go
git commit -m "chore: Remove deprecated cookie and CSRF files

Files removed:
- internal/auth/cookie.go (cookie management deprecated)
- internal/middleware/csrf.go (CSRF protection no longer needed)
- internal/middleware/csrf_deprecated.go (temporary file)

Migration to JWT Bearer tokens complete.

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

### Option 2: Manual Delete + Git
```bash
# Windows Command Prompt
cd "d:\--- Project ---\-- Golang --\backend-portfolio"
del internal\auth\cookie.go
del internal\middleware\csrf.go
del internal\middleware\csrf_deprecated.go
git add -A
git commit -m "chore: Remove deprecated cookie and CSRF files"
```

### Option 3: Using PowerShell
```powershell
cd "d:\--- Project ---\-- Golang --\backend-portfolio"
Remove-Item internal/auth/cookie.go -Force
Remove-Item internal/middleware/csrf.go -Force
Remove-Item internal/middleware/csrf_deprecated.go -Force
git add -A
git commit -m "chore: Remove deprecated cookie and CSRF files"
```

## Verification

After deletion, verify files are removed:
```bash
git status   # Should show 3 files deleted
git log      # Should show new commit
```

## Why These Files?

✅ **JWT Bearer tokens** now used for authentication
✅ **CSRF protection** not needed (Bearer tokens immune to CSRF)
✅ **No cookies** for auth tokens anymore
✅ **Deprecated stubs** were only placeholders during migration

---

**Status:** Ready for cleanup
**Delete Date:** Now
**Impact:** Code cleanup only - no functional changes
