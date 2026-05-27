# Low Issues Implementation Plan

## Backend

### L-1. Bug — Short ID collision potential
- **File:** `apps/backend/internal/model/ticket.go`
- **Deskripsi:** Short IDs dihasilkan tanpa deteksi collision. Error handling untuk unique constraint violation mungkin tidak mengembalikan pesan user-friendly.
- **Solusi:** Tambah retry logic atau tingkatkan panjang short ID.

### L-2. Bug — Warranty ticket dibuat dengan status/payment yang tidak biasa
- **File:** `apps/backend/internal/service/warranty_claim_service.go` L134-135
- **Deskripsi:** Saat approve warranty, tiket baru dibuat dengan `StatusOnProcess` + `PaymentPaid`. Tiket `on_process` tapi `paid` bisa membingungkan logika bisnis (misalnya kalkulasi revenue).

### L-3. Bad Practice — Validator instance tidak di-share
- **File:** `apps/backend/internal/service/ticket_service.go` L33, `warranty_claim_service.go` L36
- **Deskripsi:** Setiap service membuat `validator.New()` sendiri. Validator thread-safe dan seharusnya di-share.
- **Solusi:** Buat singleton atau inject instance.

### L-4. Bad Practice — `MapRepositoryError` selalu map ke `ErrTicketNotFound`
- **File:** `apps/backend/internal/service/errors.go` L90-91
- **Deskripsi:** `ErrNotFound` dari repo selalu dimapping ke `ErrTicketNotFound`. Untuk warranty claims, seharusnya `ErrWarrantyClaimNotFound`.

### L-5. Refactor — SQL query duplication masif
- **File:** `apps/backend/internal/repository/ticket_repo.go`
- **Deskripsi:** SELECT column list duplikat 4 kali. UPDATE query duplikat antara `Update()` dan `UpdateTx()`. INSERT query duplikat di `Create/CreateTx`.
- **Solusi:** Ekstrak column lists ke constants. Merge `Update` untuk delegate ke `UpdateTx`.

### L-6. Refactor — Logika duplikat pembuatan ticket di `ApproveClaim` dan `VoidClaim`
- **File:** `apps/backend/internal/service/warranty_claim_service.go`
- **Deskripsi:** Keduanya membangun `model.Ticket` dengan field mappings yang hampir identik.
- **Solusi:** Ekstrak helper `newClaimTicket(parentTicket, claim, opts)`.

### L-7. Refactor — Masking/normalization helpers di file service
- **File:** `apps/backend/internal/service/ticket_service.go` L250-295
- **Deskripsi:** Logika phone masking, name masking, dan normalization ada di service file. Seharusnya reusable.
- **Solusi:** Pindah ke package `internal/util`.

### L-8. Refactor — `Update()` method tidak pernah dipanggil
- **File:** `apps/backend/internal/repository/ticket_repo.go` L176-214
- **Deskripsi:** Non-transactional `Update()` didefinisikan tapi `UpdateTicket` selalu menggunakan `UpdateTx()`. Dead code.
- **Solusi:** Hapus atau dokumentasikan alasan dipertahankan.

## Frontend

### L-9. Bad Practice — Mock short_id format tidak match backend
- **File:** `apps/frontend/src/lib/mocks/mockData.ts`
- **Deskripsi:** Mock `short_id` seperti `"TK-001"` mungkin tidak sesuai format dari backend.

### L-10. Bad Practice — `handlers.ts` event typed as `any`
- **File:** `apps/frontend/src/lib/mocks/handlers.ts` L4
- **Solusi:** Import `RequestEvent` dari `@sveltejs/kit`.

### L-11. Bad Practice — Multiple a11y ignore comments
- **File:** `apps/frontend/src/routes/+page.svelte` L1169-1176
- **Deskripsi:** Accessibility warnings disuppress pada edit drawer backdrop. Keyboard users tidak bisa dismiss dengan Escape.

### L-12. UX — Bahasa campuran di step descriptions track page
- **File:** `apps/frontend/src/routes/track/[id]/+page.svelte` L13
- **Deskripsi:** Step 3 desc in English sementara lainnya Indonesian.

### L-13. UX — Stats section horizontal scroll di mobile
- **File:** `apps/frontend/src/routes/+page.svelte` L485-486
- **Deskripsi:** Stats menggunakan `min-w-[960px]` memaksa horizontal scroll. 5 stat cards dalam scroll kurang optimal.
- **Solusi:** Grid 2x2+1 layout di mobile.

### L-14. Performance — `Intl.NumberFormat` di-construct setiap pemanggilan
- **File:** `apps/frontend/src/routes/+page.svelte` L451
- **Deskripsi:** `formatCurrency()` membuat instance `Intl.NumberFormat` baru setiap dipanggil.
- **Solusi:** Cache formatter sebagai module-level constant.

### L-15. Performance — Google Fonts via CSS @import (render-blocking)
- **File:** `apps/frontend/src/app.css` L1
- **Deskripsi:** `@import url('https://fonts.googleapis.com/...')` menunda first paint.
- **Solusi:** Gunakan `<link rel="preload" as="style">` di `app.html` atau self-host font.

## Infrastructure

### L-16. Makefile — Missing `.PHONY` declarations
- **File:** `Makefile` L49
- **Deskripsi:** Hanya beberapa target yang dideklarasi `.PHONY`.

### L-17. Makefile — No NAME guard untuk `migrate-create`
- **File:** `Makefile` L28-29
- **Deskripsi:** `make migrate-create` tanpa `NAME=...` membuat migration tanpa nama.
- **Solusi:** Tambah guard: `@test -n "$(NAME)" || (echo "NAME required" && exit 1)`.

### L-18. Makefile — Tidak ada `make lint` atau `make fmt`
- **File:** `Makefile`
- **Solusi:** Tambah target lint dan format untuk backend dan frontend.

### L-19. Makefile — Tidak ada `make clean`
- **File:** `Makefile`
- **Solusi:** Tambah target untuk clean build artifacts, test caches, container volumes.

### L-20. Migration — Backup table dari 000005 tidak pernah dibersihkan
- **File:** `apps/backend/migrations/000005_remove_warranty_expiry_date.up.sql`
- **Deskripsi:** `_migration_000005_backup` table tidak pernah di-drop.

### L-21. Migration — Tidak ada index pada `tickets.payment_status`
- **File:** `apps/backend/migrations/`
- **Deskripsi:** Dashboard stats query unpaid repairs by `payment_status` tanpa index.

### L-22. Docs — `feature_edge.md` di project root, bukan di `docs/`
- **File:** `feature_edge.md`
- **Solusi:** Pindah ke `docs/feature_edge.md`.

### L-23. Docs — README `up` target description incomplete
- **File:** `README.md`
- **Deskripsi:** Tabel Makefile targets tidak menyebutkan `down`/`stop` targets atau `compose-test-*`.

### L-24. Docs — README warranty_expiry_date misleading
- **File:** `README.md` L33
- **Deskripsi:** Menyatakan "computes warranty_expiry_date" yang bisa disalahartikan sebagai database storage. Column sudah dihapus di migration 000005.

### L-25. Dependencies — `@supabase/supabase-js` unused
- **File:** `apps/frontend/package.json` L27
- **Solusi:** `npm uninstall @supabase/supabase-js`.

### L-26. Project — `.vscode/` committed meskipun di-gitignore
- **File:** `apps/frontend/.vscode/`
- **Solusi:** `git rm -r --cached apps/frontend/.vscode/`.

### T-6. Testing Gap — Backend
- **Deskripsi:** `idempotency_storage.go` tidak punya unit test — hanya ditest secara implisit via integration.

### T-7. Testing Gap — Backend
- **Deskripsi:** Tidak ada benchmark tests atau load testing scripts.

### P-6. Performance — Frontend
- **File:** `+page.svelte` L451
- **Deskripsi:** `Intl.NumberFormat` di-construct setiap call.

### P-7. Performance — Frontend
- **File:** `app.css` L1
- **Deskripsi:** Google Fonts via render-blocking CSS @import.
