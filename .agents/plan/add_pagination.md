# Rencana: Implementasi Pagination pada Endpoint List

## Deskripsi Masalah
Endpoint untuk mengambil daftar (List) *tickets* maupun *warranty claims* saat ini mengambil seluruh baris data dari database tanpa adanya limitasi (seperti klausa `LIMIT` dan `OFFSET` pada SQL).
Seiring dengan bertambahnya volume data di sistem OpenBench, pengembalian seluruh data ini akan memicu *memory leak* dan *bottleneck* performa yang parah di sisi backend dan database.

## Detail Lokasi Kode
1. **Ticket Repository**: `apps/backend/internal/repository/ticket_repo.go` (Method `List`)
2. **Warranty Claim Repository**: `apps/backend/internal/repository/warranty_claim_repo.go` (Method `List`)
3. Serta seluruh layer *Service* dan *Handler* yang terkait (`ticket_service.go`, `ticket_handler.go`, `warranty_claim_service.go`, `warranty_claim_handler.go`).

## Rencana Solusi
1. **Pilih Metode Pagination**: Gunakan *offset-based pagination* (menyediakan parameter `page` dan `limit`) atau *cursor-based pagination* (lebih direkomendasikan untuk sorting berbasis tanggal `entry_date DESC`).
2. **Ubah DTO**: Buat struktur request DTO standar untuk parameter pagination (contoh: `ListQueryRequest` dengan field `limit` dan `page`/`offset`).
3. **Ubah Repository**: Tambahkan filter `LIMIT` dan `OFFSET` pada query SQL di metode `List()`.
4. **Update Response DTO**: Ganti response *array list* langsung menjadi *paginated response* yang mungkin menyertakan metadata seperti `total_count`, `has_next_page`, dll.
5. **Sesuaikan Handler & Service**: Parsing *query parameters* pada Fiber context, lalu passing ke layer *service* dan *repository*.
6. **Perbarui Frontend**: Sesuaikan panggilan API di aplikasi *SvelteKit* (`apps/frontend`) untuk mengkonsumsi *paginated response* tersebut (misal: penambahan tombol *Load More* atau *Pagination Controls*).
