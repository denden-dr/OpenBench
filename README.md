# OpenBench

OpenBench adalah aplikasi administrasi yang dirancang khusus untuk mempermudah operasional bisnis reparasi dan servis ponsel (HP). Aplikasi ini membantu pelaku bisnis dalam melakukan pencatatan servis masuk secara digital, terstruktur, dan mudah dilacak. 

Ke depannya, OpenBench akan terus dikembangkan untuk mencakup berbagai utilitas pendukung bisnis lainnya, menjadikannya sistem administrasi terpadu (All-in-One) untuk usaha Anda.

---

## Fitur Utama (Saat Ini & Perencanaan)

- 📝 **Pencatatan Servis Masuk**: Mencatat detail perangkat, keluhan pelanggan, dan status perbaikan secara langsung.
- 🔍 **Pelacakan Status**: Memantau perangkat mana yang sedang dikerjakan, menunggu suku cadang, atau sudah selesai.
- 🚀 **Utilitas Tambahan (Segera Hadir)**: Pengembangan lebih lanjut untuk mencakup manajemen inventaris suku cadang, pencatatan keuangan dasar, dan manajemen data pelanggan.

## Arsitektur Teknis

Meskipun ditujukan untuk kemudahan bisnis, OpenBench dibangun di atas fondasi teknologi modern yang menjamin kecepatan, keamanan, dan keandalan tingkat tinggi (Enterprise-Grade):

- **Backend**: [Golang](https://go.dev/) dengan framework [Fiber v3](https://gofiber.io/) (Kinerja sangat cepat dan efisien).
- **Database**: [PostgreSQL 16](https://www.postgresql.org/) (Sistem database relasional paling tangguh).
- **Driver Database**: pgxpool (Untuk manajemen koneksi *pooling* berkinerja tinggi).
- **Infrastruktur**: Containerized menggunakan Docker/Podman untuk kemudahan *deployment*.
- **Konfigurasi**: 12-Factor App methodology dengan *strict environment variables loading*.

---

## Panduan Instalasi (Untuk Developer / Teknisi)

### Persyaratan Sistem
- [Go](https://go.dev/doc/install) versi 1.21 atau lebih baru.
- [Docker](https://docs.docker.com/get-docker/) atau [Podman](https://podman.io/) (beserta docker-compose/podman-compose).
- Make utility.

### Langkah-langkah Menjalankan Aplikasi Lokal

1. **Persiapan Konfigurasi**
   Salin template konfigurasi dan sesuaikan jika perlu:
   ```bash
   cp .env.example .env
   ```
   *(Pastikan tidak mengubah nama-nama variabel, karena sistem menggunakan validasi ketat/strict mode).*

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
   Server akan berjalan dan secara otomatis mencoba terhubung ke database (dilengkapi dengan sistem ketahanan *Exponential Backoff Retry* jika database terlambat merespon).

4. **Mematikan Sistem**
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
