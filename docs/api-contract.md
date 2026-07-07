# API Contract - OpenBench v1

Dokumen ini merangkum kontrak API (API Contract) untuk Web API OpenBench yang mengatur operasional pencatatan servis ponsel (HP).

## Base URL
Semua endpoint di bawah ini diasumsikan berjalan pada server lokal/produksi. Base path untuk seluruh rute API adalah `/api/v1`.

```text
http://localhost:3000
```

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
  "cost": 1500000
}
```

* **Response (201 Created)**:
```json
{
  "data": {
    "ticket_id": "d290f1ee-6c54-4b01-90e6-d701748f0851",
    "status": "RECEIVED",
    "customer": {
      "id": "c1482f34-31fa-4f24-9b24-2795db2ab961",
      "name": "Budi Santoso",
      "phone": "081234567890"
    },
    "device_brand": "Samsung",
    "device_model": "Galaxy S23",
    "device_passcode": "pola-letter-L",
    "issue_description": "Layar pecah dan tidak menampilkan gambar setelah terjatuh",
    "repair_action": "Ganti LCD Set Full",
    "cost": 1500000,
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
  * `status`: Filter status servis (`RECEIVED`, `REPAIRING`, `PENDING_CONFIRMATION`, `FIXED`, `COMPLETED`, `CANCELLED`). *Opsional*.
  * `search`: Pencarian nama pelanggan, nomor HP, atau nomor tiket. *Opsional*.
  * `limit`: Jumlah data per halaman. *Opsional, default 10*.
  * `offset`: Offset halaman. *Opsional, default 0*.

* **Response (200 OK)**:
```json
{
  "data": [
    {
      "ticket_id": "d290f1ee-6c54-4b01-90e6-d701748f0851",
      "status": "RECEIVED",
      "customer_name": "Budi Santoso",
      "customer_phone": "081234567890",
      "device_brand": "Samsung",
      "device_model": "Galaxy S23",
      "issue_description": "Layar pecah",
      "cost": 1500000,
      "created_at": "2026-07-07T12:30:00Z"
    }
  ]
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
    "customer": {
      "id": "c1482f34-31fa-4f24-9b24-2795db2ab961",
      "name": "Budi Santoso",
      "phone": "081234567890"
    },
    "device_brand": "Samsung",
    "device_model": "Galaxy S23",
    "device_passcode": "pola-letter-L",
    "issue_description": "Layar pecah dan tidak menampilkan gambar setelah terjatuh",
    "repair_action": "Ganti LCD Set Full",
    "cost": 1500000,
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
  "issue_description": "Layar pecah dan frame bengkok (Tambahan temuan baru)",
  "repair_action": "Ganti LCD Set Full + Luruskan Frame",
  "cost": 1700000,
  "notes": "Layar diganti dengan LCD original Samsung, frame sudah diluruskan agar rapat"
}
```

* **Response (200 OK)**:
```json
{
  "data": {
    "ticket_id": "d290f1ee-6c54-4b01-90e6-d701748f0851",
    "status": "PENDING_CONFIRMATION",
    "issue_description": "Layar pecah dan frame bengkok (Tambahan temuan baru)",
    "repair_action": "Ganti LCD Set Full + Luruskan Frame",
    "cost": 1700000,
    "notes": "Layar diganti dengan LCD original Samsung, frame sudah diluruskan agar rapat",
    "updated_at": "2026-07-07T12:30:00Z"
  }
}
```

---

## 2. Customers (Manajemen Pelanggan)

### A. Mendapatkan Daftar Pelanggan
Digunakan untuk melihat data seluruh kontak pelanggan yang pernah servis.

* **URL**: `/api/v1/admin/customers`
* **Method**: `GET`
* **Query Parameters**:
  * `search`: Pencarian nama atau nomor HP. *Opsional*.

* **Response (200 OK)**:
```json
{
  "data": [
    {
      "id": "c1482f34-31fa-4f24-9b24-2795db2ab961",
      "name": "Budi Santoso",
      "phone": "081234567890",
      "total_services": 3,
      "created_at": "2026-06-01T03:00:00Z"
    }
  ]
}
```

---

### B. Mendapatkan Riwayat Servis Pelanggan
Melihat seluruh riwayat HP yang pernah diservis oleh satu pelanggan tertentu.

* **URL**: `/api/v1/admin/customers/:customer_id/services`
* **Method**: `GET`

* **Response (200 OK)**:
```json
{
  "data": [
    {
      "ticket_id": "d290f1ee-6c54-4b01-90e6-d701748f0851",
      "status": "REPAIRING",
      "device_brand": "Samsung",
      "device_model": "Galaxy S23",
      "issue_description": "Layar pecah",
      "created_at": "2026-07-07T12:30:00Z"
    },
    {
      "ticket_id": "71a1795c-59bc-4a41-b062-8ff5e902b79e",
      "status": "COMPLETED",
      "device_brand": "Apple",
      "device_model": "iPhone 11",
      "issue_description": "Ganti Baterai Health 70%",
      "created_at": "2026-06-01T03:00:00Z"
    }
  ]
}
```
