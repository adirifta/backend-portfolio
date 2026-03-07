# Visitor Tracking API Documentation

Fitur untuk mencatat dan menampilkan data pengunjung website portfolio secara otomatis.
Setiap kunjungan akan tercatat dan dapat dilihat statistiknya per **hari, minggu, bulan, dan tahun**.

---

## Table of Contents

- [Visitor Tracking API Documentation](#visitor-tracking-api-documentation)
  - [Table of Contents](#table-of-contents)
  - [Database Schema](#database-schema)
    - [Tabel `visitors`](#tabel-visitors)
    - [Indexes](#indexes)
  - [API Endpoints](#api-endpoints)
    - [1. Catat Kunjungan](#1-catat-kunjungan)
    - [2. Statistik Pengunjung (Public)](#2-statistik-pengunjung-public)
    - [3. Statistik Pengunjung (Admin)](#3-statistik-pengunjung-admin)
  - [Response Models](#response-models)
    - [VisitorStats](#visitorstats)
    - [ChartDataPoint](#chartdatapoint)
    - [PageVisitCount](#pagevisitcount)
  - [Integrasi React](#integrasi-react)
    - [Track Visitor](#track-visitor)
    - [Tampilkan Statistik](#tampilkan-statistik)
    - [Contoh Dashboard Component](#contoh-dashboard-component)
  - [Arsitektur](#arsitektur)
    - [File Structure](#file-structure)
    - [Statistik yang Tersedia](#statistik-yang-tersedia)

---

## Database Schema

### Tabel `visitors`

| Kolom       | Tipe                     | Deskripsi                          |
| ----------- | ------------------------ | ---------------------------------- |
| id          | BIGSERIAL (PK)           | Primary key                        |
| ip_address  | VARCHAR(45), NOT NULL    | IP address pengunjung              |
| user_agent  | TEXT                     | Browser / device info              |
| path        | VARCHAR(255), NOT NULL   | Halaman yang dikunjungi            |
| referer     | VARCHAR(512)             | URL asal pengunjung                |
| country     | VARCHAR(100)             | Negara (opsional, untuk GeoIP)     |
| city        | VARCHAR(100)             | Kota (opsional, untuk GeoIP)       |
| visited_at  | TIMESTAMP WITH TIME ZONE | Waktu kunjungan                    |
| created_at  | TIMESTAMP WITH TIME ZONE | Waktu record dibuat                |

### Indexes

| Index                        | Kolom                    |
| ---------------------------- | ------------------------ |
| idx_visitors_visited_at      | visited_at               |
| idx_visitors_ip_address      | ip_address               |
| idx_visitors_path            | path                     |
| idx_visitors_ip_visited      | ip_address, visited_at   |

> Tabel dibuat otomatis via GORM AutoMigrate. SQL migration manual tersedia di `scripts/migrations/004_create_visitors_table.sql`.

---

## API Endpoints

### 1. Catat Kunjungan

Mencatat satu kunjungan pengunjung. Dipanggil oleh frontend React setiap kali halaman dimuat.

```
POST /api/visitors/track
```

**Auth:** Tidak diperlukan (public)

**Request Body** (opsional):

```json
{
  "path": "/about"
}
```

| Field | Tipe   | Wajib | Default | Deskripsi                |
| ----- | ------ | ----- | ------- | ------------------------ |
| path  | string | Tidak | `"/"`   | Path halaman yang dikunjungi |

**Response** `200 OK`:

```json
{
  "message": "Visitor recorded"
}
```

**Response** `500 Internal Server Error`:

```json
{
  "error": "Failed to record visitor"
}
```

> Data yang otomatis dicatat dari request: `ip_address`, `user_agent`, `referer`.

---

### 2. Statistik Pengunjung (Public)

Menampilkan statistik pengunjung tanpa data sensitif (tanpa IP address dan user agent).

```
GET /api/visitors/stats
```

**Auth:** Tidak diperlukan (public)

**Response** `200 OK`:

```json
{
  "today": 42,
  "this_week": 285,
  "this_month": 1200,
  "this_year": 15000,
  "total": 50000,
  "unique_today": 30,
  "unique_this_week": 180,
  "unique_this_month": 800,
  "unique_this_year": 10000,
  "unique_total": 35000,
  "daily_chart": [
    { "date": "2026-03-01", "count": 38 },
    { "date": "2026-03-02", "count": 45 },
    { "date": "2026-03-03", "count": 42 }
  ],
  "monthly_chart": [
    { "date": "2026-01", "count": 1100 },
    { "date": "2026-02", "count": 1350 },
    { "date": "2026-03", "count": 1200 }
  ],
  "top_pages": [
    { "path": "/", "count": 5000 },
    { "path": "/about", "count": 3200 },
    { "path": "/portfolio", "count": 2800 }
  ]
}
```

---

### 3. Statistik Pengunjung (Admin)

Menampilkan statistik lengkap termasuk data recent visitors (IP, user agent, dll).

```
GET /api/admin/visitors/stats
```

**Auth:** Diperlukan (Cookie JWT + CSRF Token)

**Headers:**

```
Cookie: access_token=<jwt_token>
X-XSRF-TOKEN: <csrf_token>
```

**Response** `200 OK`:

```json
{
  "today": 42,
  "this_week": 285,
  "this_month": 1200,
  "this_year": 15000,
  "total": 50000,
  "unique_today": 30,
  "unique_this_week": 180,
  "unique_this_month": 800,
  "unique_this_year": 10000,
  "unique_total": 35000,
  "daily_chart": [
    { "date": "2026-03-01", "count": 38 },
    { "date": "2026-03-02", "count": 45 }
  ],
  "monthly_chart": [
    { "date": "2026-01", "count": 1100 },
    { "date": "2026-02", "count": 1350 }
  ],
  "top_pages": [
    { "path": "/", "count": 5000 },
    { "path": "/about", "count": 3200 }
  ],
  "recent_visitors": [
    {
      "id": 1234,
      "ip_address": "103.123.45.67",
      "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)...",
      "path": "/portfolio",
      "referer": "https://google.com",
      "country": "",
      "city": "",
      "visited_at": "2026-03-03T14:30:00Z",
      "created_at": "2026-03-03T14:30:00Z"
    }
  ]
}
```

---

## Response Models

### VisitorStats

| Field             | Tipe             | Deskripsi                              |
| ----------------- | ---------------- | -------------------------------------- |
| today             | int64            | Total kunjungan hari ini               |
| this_week         | int64            | Total kunjungan minggu ini             |
| this_month        | int64            | Total kunjungan bulan ini              |
| this_year         | int64            | Total kunjungan tahun ini              |
| total             | int64            | Total kunjungan seluruhnya             |
| unique_today      | int64            | Pengunjung unik hari ini (by IP)       |
| unique_this_week  | int64            | Pengunjung unik minggu ini             |
| unique_this_month | int64            | Pengunjung unik bulan ini              |
| unique_this_year  | int64            | Pengunjung unik tahun ini              |
| unique_total      | int64            | Pengunjung unik seluruhnya             |
| daily_chart       | ChartDataPoint[] | Data chart harian (30 hari terakhir)   |
| monthly_chart     | ChartDataPoint[] | Data chart bulanan (12 bulan terakhir) |
| top_pages         | PageVisitCount[] | Halaman paling banyak dikunjungi       |
| recent_visitors   | Visitor[]        | 20 pengunjung terakhir (admin only)    |

### ChartDataPoint

| Field | Tipe   | Deskripsi                                 |
| ----- | ------ | ----------------------------------------- |
| date  | string | Tanggal (`YYYY-MM-DD`) atau bulan (`YYYY-MM`) |
| count | int64  | Jumlah kunjungan                          |

### PageVisitCount

| Field | Tipe   | Deskripsi           |
| ----- | ------ | ------------------- |
| path  | string | Path halaman        |
| count | int64  | Jumlah kunjungan    |

---

## Integrasi React

### Track Visitor

Letakkan di root component (`App.jsx`) atau layout utama agar semua halaman tercatat:

```jsx
import { useEffect } from 'react';
import { useLocation } from 'react-router-dom';

const API_URL = import.meta.env.VITE_API_URL;

function useVisitorTracker() {
  const location = useLocation();

  useEffect(() => {
    fetch(`${API_URL}/api/visitors/track`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ path: location.pathname }),
    }).catch(() => {
      // Silently ignore tracking errors — jangan ganggu UX
    });
  }, [location.pathname]);
}

// Gunakan di App.jsx:
function App() {
  useVisitorTracker();

  return (
    <Routes>
      {/* ... routes ... */}
    </Routes>
  );
}
```

### Tampilkan Statistik

```jsx
import { useState, useEffect } from 'react';

const API_URL = import.meta.env.VITE_API_URL;

function useVisitorStats() {
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`${API_URL}/api/visitors/stats`)
      .then((res) => res.json())
      .then(setStats)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  return { stats, loading };
}
```

### Contoh Dashboard Component

```jsx
import { useVisitorStats } from '../hooks/useVisitorStats';

function VisitorDashboard() {
  const { stats, loading } = useVisitorStats();

  if (loading) return <p>Loading...</p>;
  if (!stats) return <p>Gagal memuat statistik</p>;

  return (
    <div className="visitor-dashboard">
      <h2>Statistik Pengunjung</h2>

      <div className="stats-grid">
        <StatCard label="Hari Ini" value={stats.today} unique={stats.unique_today} />
        <StatCard label="Minggu Ini" value={stats.this_week} unique={stats.unique_this_week} />
        <StatCard label="Bulan Ini" value={stats.this_month} unique={stats.unique_this_month} />
        <StatCard label="Tahun Ini" value={stats.this_year} unique={stats.unique_this_year} />
        <StatCard label="Total" value={stats.total} unique={stats.unique_total} />
      </div>

      {/* Gunakan library chart seperti recharts / chart.js */}
      <h3>Kunjungan 30 Hari Terakhir</h3>
      {/* stats.daily_chart → [{date, count}, ...] */}

      <h3>Kunjungan per Bulan</h3>
      {/* stats.monthly_chart → [{date, count}, ...] */}

      <h3>Halaman Terpopuler</h3>
      <ul>
        {stats.top_pages?.map((page) => (
          <li key={page.path}>
            {page.path} — {page.count} kunjungan
          </li>
        ))}
      </ul>
    </div>
  );
}

function StatCard({ label, value, unique }) {
  return (
    <div className="stat-card">
      <h4>{label}</h4>
      <p className="total">{value.toLocaleString()}</p>
      <p className="unique">{unique.toLocaleString()} unik</p>
    </div>
  );
}
```

---

## Arsitektur

```
Frontend (React)                    Backend (Go + Gin)
┌─────────────────┐               ┌──────────────────────────────┐
│                  │  POST /track  │  handler.TrackVisitor()      │
│  useEffect() ───┼──────────────►│    └─► repo.Record()         │
│  (setiap page   │               │          └─► INSERT visitors  │
│   navigation)   │               │                               │
│                  │  GET /stats   │  handler.GetVisitorStatsPublic│
│  Dashboard ◄────┼──────────────►│    └─► repo.GetStats()       │
│  Component      │               │          └─► SELECT + GROUP   │
└─────────────────┘               └──────────────────────────────┘
                                            │
                                  ┌─────────▼─────────┐
                                  │   PostgreSQL       │
                                  │   visitors table   │
                                  └───────────────────┘
```

### File Structure

```
models/
  visitor.go              # Model: Visitor, VisitorStats, ChartDataPoint, PageVisitCount

internal/
  repository/
    interfaces.go         # VisitorRepository interface
    visitor.go            # GORM implementation

  handler/
    visitor.go            # TrackVisitor, GetVisitorStats, GetVisitorStatsPublic

  router/
    router.go             # Route registration

scripts/
  migrations/
    004_create_visitors_table.sql   # Manual SQL migration
```

### Statistik yang Tersedia

| Metrik            | Periode                    | Tipe          |
| ----------------- | -------------------------- | ------------- |
| Total kunjungan   | Hari / Minggu / Bulan / Tahun / All-time | Semua visits |
| Unique visitors   | Hari / Minggu / Bulan / Tahun / All-time | Berdasarkan IP |
| Daily chart       | 30 hari terakhir           | Per hari      |
| Monthly chart     | 12 bulan terakhir          | Per bulan     |
| Top pages         | All-time                   | Top 10        |
| Recent visitors   | Terbaru                    | 20 terakhir (admin only) |
