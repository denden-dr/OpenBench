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

- **Backend**: [Golang](https://go.dev/) dengan framework [Fiber v3](https://gofiber.io/) (Kinerja sangat cepat dan efisien).
- **Database**: [PostgreSQL 16](https://www.postgresql.org/) (Sistem database relasional paling tangguh).
- **Database Access**: `sqlx` dipadukan dengan `pgx/v5` untuk eksekusi SQL yang efisien, aman, dan mendukung manajemen transaksi basis data berkinerja tinggi.
- **Infrastruktur**: Containerized menggunakan Docker/Podman (dilengkapi dengan *Testcontainers* untuk *integration test*).
- **Konfigurasi**: 12-Factor App methodology dengan *strict environment variables loading* (menggunakan Viper).

---

## Panduan Instalasi (Untuk Developer / Teknisi)

### Persyaratan Sistem
- [Go](https://go.dev/doc/install) versi 1.25 atau lebih baru.
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

3. **Menjalankan Server Web API**
   Jalankan backend server:
   ```bash
   make run
   ```
   Server akan berjalan dan secara otomatis mencoba terhubung ke database (dilengkapi dengan sistem ketahanan *Exponential Backoff Retry* jika database lambat merespon).

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
