# OpenBench

OpenBench adalah aplikasi administrasi yang dirancang khusus untuk mempermudah operasional bisnis reparasi dan servis ponsel (HP) maupun alat elektronik lainnya. Aplikasi ini membantu pelaku bisnis dalam melakukan pencatatan servis masuk secara digital, terstruktur, dan mudah dilacak. 

Seiring perkembangannya, OpenBench kini telah dilengkapi berbagai utilitas pendukung bisnis yang menjadikannya sistem administrasi terpadu (All-in-One) untuk usaha Anda.

---

## Fitur Utama

- 🔐 **Autentikasi & Keamanan (Auth)**: Sistem login menggunakan JWT (Access & Refresh token) yang aman dan berbasis *cookie*.
- 📝 **Manajemen Tiket Servis (Ticket)**: Mencatat detail perangkat, keluhan pelanggan, dan melacak status perbaikan secara sistematis.
- 📦 **Manajemen Inventaris (Inventory)**: Mengelola stok produk, ketersediaan suku cadang, dan penyesuaian (*adjustment*) stok barang.
- 💰 **Point of Sale (POS)**: Sistem kasir terpadu yang terhubung dengan tiket dan inventaris (mendukung transaksi *atomic* yang konsisten).
- 🛡️ **Klaim Garansi (Warranty)**: Pendataan dan pengecekan masa garansi untuk setiap perbaikan yang telah diselesaikan.

## Arsitektur Teknis

Meskipun ditujukan untuk kemudahan bisnis, OpenBench dibangun di atas fondasi teknologi modern yang menjamin kecepatan, keamanan, dan keandalan tingkat tinggi (Enterprise-Grade):

- **Fullstack Architecture**: Menggunakan pendekatan **GOTTH Stack** (Go, Templ, Tailwind, HTMX) untuk pengembangan web interaktif tanpa penderitaan SPA kompleks.
- **Backend**: [Golang](https://go.dev/) dengan framework [Fiber v3](https://gofiber.io/) (Kinerja sangat cepat dan efisien).
- **Frontend Engine**: Komponen HTML reaktif menggunakan [Templ](https://templ.guide/), penataan gaya lewat Tailwind CSS (Binary CLI), serta interaktivitas dari [HTMX](https://htmx.org/) dan Alpine.js.
- **Database**: [PostgreSQL 16](https://www.postgresql.org/) dengan akses melalui `sqlx` dan `pgx/v5` (eksekusi SQL efisien & mendukung *atomic transactions*).
- **Infrastruktur**: Containerized menggunakan Docker/Podman (dilengkapi dengan *Testcontainers* untuk *integration test*).
- **UI Design System**: Mengimplementasikan estetika *Glassmorphism*, *Self-hosted custom fonts*, dan standar Lucide Icons.

---

## Panduan Instalasi (Untuk Developer / Teknisi)

### Persyaratan Sistem
- [Go](https://go.dev/doc/install) versi 1.26.1 atau lebih baru.
- [Docker](https://docs.docker.com/get-docker/) atau [Podman](https://podman.io/) (beserta docker-compose/podman-compose).
- Make utility.

### Langkah-langkah Menjalankan Aplikasi Lokal

1. **Persiapan Konfigurasi**
   Salin template konfigurasi dan sesuaikan jika perlu:
   ```bash
   cp .env.example .env
   ```
   *(Pastikan tidak mengubah nama-nama variabel secara sembarangan, karena sistem menggunakan validasi ketat/strict mode).*

2. **Menjalankan Database**
   Nyalakan container PostgreSQL di latar belakang:
   ```bash
   make up
   ```

3. **Menjalankan Aplikasi (Mode Development)**
   Untuk pengembangan dengan pengalaman terbaik yang dilengkapi fitur *hot-reload* (otomatis me-rebuild Go, Templ, dan Tailwind CSS setiap ada perubahan file), jalankan:
   ```bash
   make dev
   ```
   *(Pastikan Anda telah menginstal utilitas Air melalui `go install github.com/air-verse/air@latest`)*.
   
   Jika Anda hanya ingin menjalankan aplikasi sekali jalan (tanpa *hot-reload*):
   ```bash
   make run
   ```

4. **Menjalankan Pengujian (Testing)**
   Proyek ini dilengkapi dengan modul pengujian dan integrasi otomatis yang mendalam. Jalankan perintah berikut untuk menguji seluruh logika:
   ```bash
   go test ./...
   ```

5. **Mematikan Sistem**
   Untuk mematikan container database:
   ```bash
   make down
   ```

## Verifikasi Kesehatan Sistem

Setelah aplikasi berjalan, Anda dapat memverifikasi bahwa server API dan koneksi Database aktif melalui *endpoint health-check*:

```bash
curl http://localhost:3000/health
```

**Respon sukses:**
```json
{
  "database": "up",
  "status": "up",
  "timestamp": "2026-07-07T19:00:00+07:00"
}
```

---
*Dibuat untuk efisiensi bisnis, direkayasa untuk masa depan.*
