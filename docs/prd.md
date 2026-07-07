# Product Requirements Document (PRD) - OpenBench

**Versi:** 1.0  
**Status:** Draft / Perencanaan Aktif  
**Fokus:** Minimum Viable Product (MVP) - Sistem Pencatatan Servis Internal

---

## 1. Ringkasan Eksekutif (Executive Summary)
**OpenBench** adalah sistem administrasi dan pencatatan (SaaS internal) yang dirancang khusus untuk pelaku bisnis perbaikan ponsel (Servis HP). Produk ini berfokus untuk menggantikan proses pencatatan manual (kertas/buku) dengan sistem digital terpusat yang aman, cepat, dan mudah dilacak. MVP saat ini difokuskan penuh pada pengelolaan tiket servis dan data pelanggan untuk admin toko dan teknisi.

## 2. Tujuan & Sasaran (Objectives & Goals)
- **Digitalisasi Pencatatan**: Menghindari hilangnya data servis akibat nota kertas yang hilang atau rusak.
- **Transparansi Alur Servis**: Menyediakan pelacakan status perbaikan yang jelas (Diterima -> Dikerjakan -> Selesai).
- **Efisiensi Admin**: Mempercepat proses pendaftaran barang masuk (sekali input untuk data pelanggan, perangkat, dan keluhan).
- **Rekam Jejak Berharga**: Menyimpan riwayat servis setiap pelanggan untuk membangun loyalitas pelanggan dan referensi garansi.

## 3. Target Pengguna (Target Audience)
- **Pemilik Bisnis (Business Owner)**: Memantau total HP yang sedang masuk, proses yang berjalan, dan pendapatan.
- **Admin Toko (Frontdesk)**: Menerima barang dari pelanggan, membuat tiket servis, dan menyerahkan barang saat selesai.
- **Teknisi (Technician)**: Melihat daftar HP yang harus diperbaiki, mengupdate status pengerjaan, dan memasukkan catatan teknis beserta biaya final.

---

## 4. Ruang Lingkup MVP (In-Scope for v1)

### 4.1. Manajemen Tiket Servis (Service Ticketing)
- Pembuatan tiket servis baru dengan sistem *bundling* (menginput nama pelanggan, detail HP, dan keluhan sekaligus).
- Melacak dan memperbarui status tiket:
  - `RECEIVED`: HP diterima oleh admin dari pelanggan.
  - `REPAIRING`: Teknisi sedang melakukan pengecekan atau perbaikan.
  - `PENDING_CONFIRMATION`: Perbaikan tertunda karena butuh konfirmasi pelanggan (misal: ada isu tambahan, salah diagnosa, atau perubahan biaya).
  - `FIXED`: Perbaikan selesai, perangkat siap diambil oleh pelanggan.
  - `COMPLETED`: Perangkat sudah diserahkan kembali ke pelanggan dan transaksi selesai (ditutup).
  - `CANCELLED`: Servis dibatalkan (karena mesin tidak bisa diperbaiki, biaya tidak cocok, dsb).
- Menetapkan Total Biaya (`cost`) dan Tindakan Perbaikan (`repair_action`) di awal berdasarkan diagnosa. Jika di tengah jalan ada kerusakan tambahan, admin cukup meng-update biaya dan tindakan tersebut, lalu menunggu persetujuan pelanggan.
- Memberikan catatan internal teknisi untuk referensi di masa depan.

### 4.2. Manajemen Pelanggan & Perangkat
- Sistem secara otomatis mencatat dan menyimpan kontak pelanggan baru jika belum pernah terdaftar.
- Melihat daftar pelanggan dan riwayat akumulasi transaksi mereka (berpotensi berguna jika pelanggan mengklaim garansi).

---

## 5. Alur Pengguna (User Flow)

1. **Alur Penerimaan Barang (Frontdesk)**
   - Pelanggan datang membawa HP rusak.
   - Admin bertanya nama dan nomor WA.
   - Admin membuat tiket servis baru (Input: Merek HP, Keluhan, Pola Kunci Layar, Nama, WA).
   - Tiket tercipta dengan status `RECEIVED`.

2. **Alur Pengerjaan (Teknisi & Admin)**
   - Teknisi melihat daftar servis yang berstatus `RECEIVED`.
   - Teknisi mulai mengecek/mengerjakan HP, dan merubah status menjadi `REPAIRING`.
   - *Kasus Khusus*: Jika ditemukan kerusakan tambahan atau salah diagnosa awal, teknisi melapor ke admin. Admin kemudian memperbarui `Tindakan Perbaikan` dan `Total Biaya` di sistem, lalu mengubah status ke `PENDING_CONFIRMATION` hingga ada persetujuan dari pelanggan.
   - Setelah perbaikan selesai dan dites normal, teknisi merubah status menjadi `FIXED` (siap diambil) dan menambahkan Catatan (jika ada).

3. **Alur Pengambilan Barang (Frontdesk)**
   - Pelanggan datang mengambil HP yang berstatus `FIXED` (atau `CANCELLED`).
   - Admin mencari nomor tiket atau nama pelanggan di sistem.
   - Admin memverifikasi pembayaran, menyerahkan HP, lalu mengubah status akhir menjadi `COMPLETED` (transaksi ditutup).

---

## 6. Persyaratan Teknis (Technical Requirements)

- **Backend**: Golang dengan framework Fiber v3 untuk menunjang konkurensi dan respons ultra-cepat.
- **Database**: PostgreSQL 16 dengan proteksi integritas data dan penggunaan `pgxpool`.
- **Infrastruktur**: Containerization menggunakan Docker, dengan pola konfigurasi *fail-fast* melalui variabel lingkungan (`.env`).
- **Pola Desain API**: RESTful JSON API. Setiap endpoint admin berjalan di bawah rute `/api/v1/admin/*`.

---

## 7. Di Luar Ruang Lingkup Saat Ini (Out of Scope for v1)
Fitur berikut tidak akan dikerjakan pada rilis awal (MVP) namun dipertimbangkan untuk versi masa depan (v2/v3):
1. **Notifikasi WhatsApp Otomatis**: Mengirim WA ke pelanggan secara otomatis saat status HP berubah menjadi `DONE`.
2. **Portal Pengecekan Publik**: Halaman web statis (dengan rute `/api/v1/public/*`) di mana pelanggan bisa melacak progres servis dengan memasukkan nomor nota/tiket.
3. **Manajemen Inventaris Sparepart**: Memotong stok layar/baterai secara otomatis ketika tiket servis `DONE`.
4. **Modul Akuntansi/Keuangan**: Laporan laba rugi, kasir otomatis, cetak struk (thermal printer integration).
