### Install dependencies
```bash
go mod tidy
```

## Jalankan aplikasi:
```bash
go run main.go
```

## Masuk database postgresql
```bash
docker exec -it backend-portfolio-db-1 psql -U portfolio_user -d portfolio_db
```

## Lihat semua tabel
```bash
\dt
```

## Lihat data users
```bash
SELECT * FROM users;
```

## Lihat data about
```bash
SELECT * FROM abouts;
```

## Keluar dari psql
```bash
\q
```

### Cara Menambah/Update kolom database postgresql
- Buat File Migration SQL (Buat file baru scripts/migrations/001_add_new_columns_to_abouts.sql)
- Update Model About dengan kolom baru
- Update Docker Compose untuk Menjalankan Migration (Pada baris volumes => ./scripts/migrations:/docker-entrypoint-initdb.d/migrations)
- Cara Menjalankan Migration (Copy File Migration dan Restart):
# Copy file migration ke container
```bash
docker cp scripts/migrations/001_add_new_columns_to_abouts.sql backend-portfolio-db-1:/tmp/
```

# Jalankan migration
```bash
docker exec -it backend-portfolio-db-1 psql -U portfolio_user -d portfolio_db -f /tmp/001_add_new_columns_to_abouts.sql
```

### Cara Delete kolom database postgresql
- Masuk ke `docker exec -it backend-portfolio-db-1 psql -U portfolio_user -d portfolio_db`
- Lalu bisa `ALTER TABLE nama_tabel DROP COLUMN nama_kolom;`
- Jika mau drop beberapa kolom sekaligus:
`
ALTER TABLE nama_tabel 
DROP COLUMN kolom1,
DROP COLUMN kolom2;
`

## Deploy Backend Ke Google Cloud Console
- Cloud Build 
- Artifacts Registry
- Cloud SQL (Postgresql)
- Cloud Run