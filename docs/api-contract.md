# API Contract - OpenBench v1

Dokumen ini merangkum kontrak API (API Contract) untuk Web API OpenBench yang mengatur operasional pencatatan servis ponsel (HP).

## Base URL
Semua endpoint di bawah ini diasumsikan berjalan pada server lokal/produksi. Base path untuk seluruh rute API adalah `/api/v1`.

```text
http://localhost:3000
```

## Keamanan & Otorisasi (Security)
Seluruh rute yang berada di bawah prefix `/api/v1/admin/*` adalah **rute terproteksi (Protected Routes)**. 
Untuk dapat mengakses rute-rute ini, klien **wajib** menyertakan *cookie* otentikasi `access_token` yang valid (dihasilkan dari endpoint Login). `access_token` juga dikembalikan di *response body* Login untuk fleksibilitas klien. 

* Jika *cookie* tidak dikirimkan, kedaluwarsa, atau tidak valid, server akan menolak permintaan dan mengembalikan status HTTP `401 Unauthorized` (sesuai format *Problem Details*).
* Rute publik yang tidak terproteksi di bawah `/api/v1/admin/*` hanyalah rute untuk *Login*.

## Format Respon Umum (Standard Response)
Semua respon sukses akan mengembalikan data di dalam objek `"data"`, sedangkan kegagalan akan mengembalikan detail kesalahan sesuai standar RFC 7807 (Problem Details for HTTP APIs) dengan tipe konten `application/problem+json`.

**Format Sukses:**
```json
{
  "data": {}
}
```

**Format Error:**
```json
{
  "type": "https://openbench.local/errors/bad-request",
  "title": "Bad Request",
  "status": 400,
  "detail": "Keterangan detail kesalahan (contoh: format email tidak valid).",
  "instance": "/api/v1/admin/services"
}
```

---

## 1. Service Tickets (Pencatatan Servis)

### A. Membuat Tiket Servis Baru
Endpoint ini digunakan ketika pelanggan membawa HP-nya untuk diservis. Data pelanggan dan perangkat akan didaftarkan/dihubungkan secara otomatis. Cost (Total Biaya) dan Tindakan Perbaikan ditentukan di awal berdasarkan diagnosa.

* **URL**: `/api/v1/admin/services`
* **Method**: `POST`
* **Request Body**:
```json
{
  "customer_name": "Budi Santoso",
  "customer_phone": "081234567890",
  "device_brand": "Samsung",
  "device_model": "Galaxy S23",
  "device_passcode": "pola-letter-L",
  "issue_description": "Layar pecah dan tidak menampilkan gambar setelah terjatuh",
  "repair_action": "Ganti LCD Set Full",
  "cost": 1500000,
  "warranty_days": 30
}
```

* **Response (201 Created)**:
```json
{
  "data": {
    "ticket_id": "d290f1ee-6c54-4b01-90e6-d701748f0851",
    "status": "RECEIVED",
    "customer_name": "Budi Santoso",
    "customer_phone": "081234567890",
    "device_brand": "Samsung",
    "device_model": "Galaxy S23",
    "device_passcode": "pola-letter-L",
    "issue_description": "Layar pecah dan tidak menampilkan gambar setelah terjatuh",
    "repair_action": "Ganti LCD Set Full",
    "cost": 1500000,
    "warranty_days": 30,
    "notes": "",
    "created_at": "2026-07-07T12:30:00Z",
    "updated_at": "2026-07-07T12:30:00Z"
  }
}
```

---

### B. Mendapatkan Daftar Tiket Servis
Mendapatkan semua tiket servis, mendukung filter berdasarkan status dan pencarian berdasarkan Nomor Tiket / Nama Pelanggan / Nomor HP Pelanggan.

* **URL**: `/api/v1/admin/services`
* **Method**: `GET`
* **Query Parameters**:
  * `status`: Filter status servis (`RECEIVED`, `REPAIRING`, `PENDING_CONFIRMATION`, `FIXED`, `COMPLETED`, `CANCELLED`, `RETURNED`). *Opsional*.
  * `search`: Pencarian nama pelanggan, nomor HP, atau nomor tiket. *Opsional*.
  * `limit`: Jumlah data per halaman. *Opsional, default 10*.
  * `offset`: Offset halaman. *Opsional, default 0*.

* **Response (200 OK)**:
```json
{
  "data": [
    {
      "ticket_id": "d290f1ee-6c54-4b01-90e6-d701748f0851",
      "ticket_number": "TKT-20260707-1234",
      "status": "RECEIVED",
      "customer_name": "Budi Santoso",
      "device_brand": "Samsung",
      "device_model": "Galaxy S23",
      "created_at": "2026-07-07T12:30:00Z"
    }
  ],
  "meta": {
    "total_data": 45,
    "limit": 10,
    "offset": 0,
    "total_pages": 5
  }
}
```

---

### C. Mendapatkan Detail Tiket Servis
Melihat detail penuh sebuah tiket servis spesifik.

* **URL**: `/api/v1/admin/services/:ticket_id`
* **Method**: `GET`

* **Response (200 OK)**:
```json
{
  "data": {
    "ticket_id": "d290f1ee-6c54-4b01-90e6-d701748f0851",
    "status": "RECEIVED",
    "customer_name": "Budi Santoso",
    "customer_phone": "081234567890",
    "device_brand": "Samsung",
    "device_model": "Galaxy S23",
    "device_passcode": "pola-letter-L",
    "issue_description": "Layar pecah dan tidak menampilkan gambar setelah terjatuh",
    "repair_action": "Ganti LCD Set Full",
    "cost": 1500000,
    "warranty_days": 30,
    "notes": "",
    "created_at": "2026-07-07T12:30:00Z",
    "updated_at": "2026-07-07T12:30:00Z"
  }
}
```

---

### D. Mengubah Status Servis
Mengubah alur pengerjaan servis ponsel (misalnya ketika teknisi mulai mengerjakannya, atau admin memindahkan ke status butuh konfirmasi).

* **URL**: `/api/v1/admin/services/:ticket_id/status`
* **Method**: `PATCH`
* **Request Body**:
```json
{
  "status": "PENDING_CONFIRMATION"
}
```

* **Response (200 OK)**:
```json
{
  "data": {
    "ticket_id": "d290f1ee-6c54-4b01-90e6-d701748f0851",
    "status": "PENDING_CONFIRMATION",
    "updated_at": "2026-07-07T12:30:00Z"
  }
}
```

---

### E. Memperbarui Detail Tiket Servis
Digunakan ketika ada perubahan diagnosa/kerusakan, penyesuaian tindakan perbaikan (`repair_action`), perubahan biaya, atau menambahkan catatan teknis kerusakan.

* **URL**: `/api/v1/admin/services/:ticket_id`
* **Method**: `PUT`
* **Request Body**:
```json
{
  "customer_name": "Budi Santoso",
  "customer_phone": "081234567891",
  "issue_description": "Layar pecah dan frame bengkok (Tambahan temuan baru)",
  "repair_action": "Ganti LCD Set Full + Luruskan Frame",
  "cost": 1700000,
  "warranty_days": 30,
  "notes": "Layar diganti dengan LCD original Samsung, frame sudah diluruskan agar rapat"
}
```

* **Response (200 OK)**:
```json
{
  "data": {
    "ticket_id": "d290f1ee-6c54-4b01-90e6-d701748f0851",
    "status": "PENDING_CONFIRMATION",
    "customer_name": "Budi Santoso",
    "customer_phone": "081234567891",
    "issue_description": "Layar pecah dan frame bengkok (Tambahan temuan baru)",
    "repair_action": "Ganti LCD Set Full + Luruskan Frame",
    "cost": 1700000,
    "warranty_days": 30,
    "notes": "Layar diganti dengan LCD original Samsung, frame sudah diluruskan agar rapat",
    "updated_at": "2026-07-07T12:30:00Z"
  }
}
```

---

### F. Emergency Edit Tiket Servis
Digunakan dalam kondisi darurat untuk mengubah seluruh data tiket servis secara komprehensif tanpa batasan aturan alur standar (termasuk mengganti status secara paksa, mengubah data perangkat, dll).

* **URL**: `/api/v1/admin/services/:ticket_id/emergency`
* **Method**: `PUT`
* **Request Body**:
```json
{
  "customer_name": "Budi Santoso",
  "customer_phone": "081234567891",
  "device_brand": "Samsung",
  "device_model": "Galaxy S23",
  "device_passcode": "pola-letter-L",
  "status": "FIXED",
  "issue_description": "Layar pecah dan frame bengkok (Tambahan temuan baru)",
  "repair_action": "Ganti LCD Set Full + Luruskan Frame",
  "cost": 1700000,
  "warranty_days": 30,
  "notes": "Layar diganti dengan LCD original Samsung, frame sudah diluruskan agar rapat"
}
```

* **Response (200 OK)**:
```json
{
  "data": {
    "ticket_id": "d290f1ee-6c54-4b01-90e6-d701748f0851",
    "status": "FIXED",
    "customer_name": "Budi Santoso",
    "customer_phone": "081234567891",
    "device_brand": "Samsung",
    "device_model": "Galaxy S23",
    "device_passcode": "pola-letter-L",
    "issue_description": "Layar pecah dan frame bengkok (Tambahan temuan baru)",
    "repair_action": "Ganti LCD Set Full + Luruskan Frame",
    "cost": 1700000,
    "warranty_days": 30,
    "notes": "Layar diganti dengan LCD original Samsung, frame sudah diluruskan agar rapat",
    "updated_at": "2026-07-07T12:35:00Z"
  }
}
```

---

## 3. Garansi & Klaim (Warranty & Claims)

Sistem garansi otomatis aktif (data `warranties` dibuat) ketika sebuah tiket berstatus `COMPLETED` dan memiliki `warranty_days` > 0.
Tanggal kedaluwarsa garansi dihitung dari saat tiket diubah menjadi `COMPLETED` ditambah dengan `warranty_days`.

### A. Cek Status Garansi (Berdasarkan Ticket ID)
Digunakan admin saat kustomer membawa nomor tiket lama untuk melihat apakah garansi masih valid.

* **URL**: `/api/v1/admin/warranties/by-ticket/:ticket_id`
* **Method**: `GET`
* **Response (200 OK)**:
```json
{
  "data": {
    "id": "w-5432-1098",
    "ticket_id": "d290f1ee-6c54-4b01-90e6-d701748f0851",
    "start_date": "2026-07-07T12:30:00Z",
    "end_date": "2026-08-06T12:30:00Z",
    "status": "ACTIVE"
  }
}
```

---

### B. Membuat Klaim Garansi Baru
Digunakan untuk mendaftarkan HP kustomer ke antrean perbaikan (klaim garansi). Hanya bisa dibuat jika garansi berstatus `ACTIVE`.

* **URL**: `/api/v1/admin/claims`
* **Method**: `POST`
* **Request Body**:
```json
{
  "warranty_id": "w-5432-1098",
  "issue_description": "Layar sentuh tidak responsif di bagian pojok kiri atas setelah diganti minggu lalu"
}
```
* **Response (201 Created)**:
```json
{
  "data": {
    "claim_id": "c-9876-5432",
    "claim_number": "CLM-20260714-0001",
    "warranty_id": "w-5432-1098",
    "status": "RECEIVED",
    "issue_description": "Layar sentuh tidak responsif di bagian pojok kiri atas setelah diganti minggu lalu",
    "repair_action": null,
    "notes": null,
    "created_at": "2026-07-14T09:00:00Z"
  }
}
```

---

### C. Mengelola Klaim (List, Detail, Update Status & Info)
Klaim garansi memiliki *lifecycle* (status) yang sama dengan tiket servis reguler, namun tabel dan endpoint-nya terpisah agar pelaporannya tidak tercampur.

* **List Claims**: `GET /api/v1/admin/claims`
* **Detail Claim**: `GET /api/v1/admin/claims/:claim_id`
* **Ubah Status**: `PATCH /api/v1/admin/claims/:claim_id/status` (Payload: `{"status": "FIXED"}`)
* **Update Teknisi**: `PUT /api/v1/admin/claims/:claim_id`
  *Payload:*
  ```json
  {
    "issue_description": "Layar sentuh tidak responsif...",
    "repair_action": "Bongkar ulang dan pasang kembali konektor flexibel LCD yang kendor",
    "notes": "Tidak ada penambahan biaya, konektor hanya kendor akibat tekanan"
  }
  ```

---

## 4. Authentication (Auth)

Sistem menggunakan skema JWT (JSON Web Token) dengan strategi *Access Token* dan *Refresh Token*. Token dikembalikan melalui dua jalur: **HTTP-Only Cookies** (wajib untuk akses rute terproteksi) dan **Response Body** (untuk fleksibilitas klien).

- `access_token`: Berumur pendek (misal: 15 menit), digunakan untuk otorisasi di seluruh rute `/api/v1/admin/*`. Dikirim sebagai *cookie* dan di *response body*.
- `refresh_token`: Berumur panjang (misal: 7 hari), eksklusif hanya dapat dikirim dan dibaca oleh rute `/api/v1/auth/refresh`. Hanya dikirim sebagai *cookie*.

### A. Sign In (Admin Login)
Endpoint untuk memvalidasi kredensial admin dan menerbitkan *cookies* autentikasi serta token di *response body*.

* **URL**: `/api/v1/auth/login`
* **Method**: `POST`
* **Request Body**:
```json
{
  "email": "admin@openbench.local",
  "password": "secretpassword123"
}
```

* **Response (200 OK)**:
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_at": "2026-07-08T12:45:00Z",
    "user": {
      "id": "u-1234-5678",
      "email": "admin@openbench.local",
      "role": "ADMIN"
    }
  }
}
```
* **Response Headers (Set-Cookie)**:
```http
Set-Cookie: access_token=eyJhb...; Path=/; HttpOnly; Secure; SameSite=Strict; Max-Age=900
Set-Cookie: refresh_token=eyJhb...; Path=/api/v1/auth/refresh; HttpOnly; Secure; SameSite=Strict; Max-Age=604800
```

---

### B. Refresh Token
Mendapatkan `access_token` yang baru. Klien tidak perlu mengatur konfigurasi khusus; browser akan otomatis menyisipkan *cookie* `refresh_token` saat memanggil endpoint ini.

* **URL**: `/api/v1/auth/refresh`
* **Method**: `POST`
* **Request Headers**: *Browser otomatis mengirim cookie `refresh_token`*

* **Response (200 OK)**:
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_at": "2026-07-08T12:45:00Z"
  }
}
```
* **Response Headers (Set-Cookie)**:
```http
Set-Cookie: access_token=eyJhb...; Path=/; HttpOnly; Secure; SameSite=Strict; Max-Age=900
```

---

### C. Sign Out (Logout)
Menghapus seluruh *cookies* sesi autentikasi pada browser dengan merubah nilai Max-Age menjadi 0.

* **URL**: `/api/v1/auth/logout`
* **Method**: `POST`
* **Request Headers**: *Browser otomatis mengirim cookie `access_token`*

* **Response (200 OK)**:
```json
{
  "data": {
    "message": "Logged out successfully"
  }
}
```
* **Response Headers (Set-Cookie)**:
```http
Set-Cookie: access_token=; Path=/; HttpOnly; Secure; SameSite=Strict; Max-Age=0
Set-Cookie: refresh_token=; Path=/api/v1/auth/refresh; HttpOnly; Secure; SameSite=Strict; Max-Age=0
```
