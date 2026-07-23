# OpenBench

OpenBench adalah aplikasi administrasi yang dirancang khusus untuk mempermudah operasional bisnis reparasi dan servis ponsel (HP) maupun alat elektronik lainnya. Aplikasi ini membantu pelaku bisnis dalam melakukan pencatatan servis masuk secara digital, terstruktur, dan mudah dilacak. 

Seiring perkembangannya, OpenBench kini telah dilengkapi berbagai utilitas pendukung bisnis yang menjadikannya sistem administrasi terpadu (All-in-One) untuk usaha Anda.

---

## Fitur Utama

- 🔐 **Autentikasi & Keamanan (Auth)**: Sistem keamanan dengan verifikasi token Bearer JWT yang aman.
- 📝 **Manajemen Tiket Servis (Ticket)**: Mencatat detail perangkat, keluhan pelanggan, dan melacak status perbaikan secara sistematis.
- 📦 **Manajemen Inventaris (Inventory)**: Mengelola stok produk, ketersediaan suku cadang, dan penyesuaian (*adjustment*) stok barang.
- 💰 **Point of Sale (POS)**: Sistem kasir terpadu yang terhubung dengan tiket dan inventaris (mendukung transaksi *atomic* yang konsisten).
- 🛡️ **Klaim Garansi (Warranty)**: Pendataan dan pengecekan masa garansi untuk setiap perbaikan yang telah diselesaikan.

## Arsitektur Teknis

OpenBench dibangun di atas fondasi teknologi *Micro Frontends* dan *Web API* modern yang menjamin kecepatan, skalabilitas, dan keandalan tingkat tinggi (Enterprise-Grade):

- **Struktur Repositori**: *Monorepo* yang mengisolasi backend dan multi-frontend.
- **Backend (Web API)**: [Golang](https://go.dev/) dengan framework [Fiber v3](https://gofiber.io/) yang berjalan sangat efisien dan menyajikan JSON API.
- **Frontend (Micro Frontends)**: Dua aplikasi terpisah untuk publik (`web-user`) dan internal admin (`web-admin`), dibangun dengan [React 19](https://react.dev/), [TypeScript](https://www.typescriptlang.org/), dan dikompilasi super cepat menggunakan [Vite](https://vitejs.dev/).
- **UI Design System**: Menggunakan [Tailwind CSS v4](https://tailwindcss.com/) dengan palet estetika *Glassmorphism*, font kustom mandiri (*self-hosted*), dan standar ikon Lucide.
- **Database**: [PostgreSQL 16](https://www.postgresql.org/) dengan akses melalui `sqlx` dan `pgx/v5` (mendukung *atomic transactions*).
- **Infrastruktur**: Database dikelola menggunakan Podman/Docker Compose, dan *Integration Test* didukung oleh ekosistem *Testcontainers*.
- **E2E Testing**: Aplikasi standalone `apps/e2e-testing` berbasis [Playwright](https://playwright.dev/), dijalankan terhadap stack terisolasi menggunakan `docker-compose.test.yml` khusus pengujian.

---

## Panduan Instalasi (Untuk Developer / Teknisi)

### Persyaratan Sistem
- [Go](https://go.dev/doc/install) versi 1.26.1 atau lebih baru.
- [Node.js](https://nodejs.org/) versi terbaru dan package manager `pnpm`.
- [Docker](https://docs.docker.com/get-docker/) atau [Podman](https://podman.io/) (beserta docker-compose/podman-compose).
- Make utility.

### Langkah-langkah Menjalankan Aplikasi Lokal

1. **Persiapan Konfigurasi**
   Salin template konfigurasi untuk Web API dan sesuaikan jika perlu:
   ```bash
   cp apps/webapi/.env.example apps/webapi/.env
   ```
   *(Catatan: pastikan nilai rahasia dalam `.env` valid, khususnya enkripsi key yang diwajibkan 32-karakter).*

2. **Instalasi Dependensi**
   Untuk mengunduh modul Go dan *package node_modules* untuk kedua aplikasi React:
   ```bash
   make install-api
   make install-user
   make install-admin
   ```

3. **Menjalankan Database (PostgreSQL)**
   ```bash
   make up
   ```

4. **Menjalankan Server (Development)**
   Buka 3 tab terminal terpisah di root repositori untuk menjalankan layanan dengan fitur *hot-reload* bawaan:

   - **Backend API Server**:
     ```bash
     make dev-api
     ```
     *(Membutuhkan Air terpasang global: `go install github.com/air-verse/air@latest`)*

   - **User Portal Frontend**:
     ```bash
     make dev-user
     ```

   - **Admin Dashboard Frontend**:
     ```bash
     make dev-admin
     ```

5. **Menjalankan Build & Pengujian**

   Untuk melakukan proses pengujian *unit* pada backend:
   ```bash
   make test-api
   ```

   Untuk pengujian *unit* pada frontend admin (Vitest):
   ```bash
   make test-admin
   ```

   Untuk menjalankan **E2E Test** secara lengkap (otomatis spin up stack, jalankan Playwright, lalu teardown):
   ```bash
   make test-e2e
   ```
   *(Membutuhkan Podman berjalan. Stack pengujian `docker-compose.test.yml` dikelola sepenuhnya oleh `scripts/test-e2e.sh`.)*

   Untuk mem-*build* file binary Go dan proses kompilasi Vite untuk kedua frontend:
   ```bash
   make build-all
   ```

6. **Mematikan Sistem**
   Untuk mematikan container database:
   ```bash
   make down
   ```

## Verifikasi Kesehatan Sistem

Setelah API backend berjalan (biasanya di port 3000), Anda dapat memverifikasi koneksi database melalui *endpoint health-check*:

```bash
curl http://localhost:3000/health
```

**Respon sukses:**
```json
{
  "database": "up",
  "status": "up",
  "timestamp": "2026-07-19T19:00:00+07:00"
}
```

---
*Dibuat untuk efisiensi bisnis, direkayasa untuk masa depan.*
