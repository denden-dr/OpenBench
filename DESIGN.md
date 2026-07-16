# OpenBench Design System & UI Guidelines

Dokumen ini adalah panduan referensi (*source of truth*) untuk pengembangan antarmuka pengguna (UI) di proyek OpenBench. Dokumen ini dirancang khusus agar konsisten ketika diproses oleh AI agent yang bekerja dalam *codebase* ini.

## 1. Tech Stack (GOTTH)
Proyek ini secara tegas menggunakan arsitektur **GOTTH**:
- **Go**: Logika *backend* dan *web handlers*.
- **Templ** (`a-h/templ`): *Engine* HTML *type-safe* untuk komponen UI.
- **Tailwind CSS**: *Utility-first CSS framework* (menggunakan versi 3 dengan *standalone binary* via `npx tailwindcss@3`).
- **HTMX**: Penggerak interaktivitas dinamis (AJAX, *partial reloads*) melalui atribut HTML.
- **Alpine.js**: Menangani *state* mikro sisi klien (seperti *dropdown, modal, dismissible alerts*).

## 2. Struktur Direktori
Pastikan Anda menempatkan file UI sesuai hierarki berikut:
```text
ui/
├── static/
│   ├── css/       # input.css (Tailwind entry) dan style.css (hasil build)
│   └── fonts/     # File lokal .woff2 (Plus Jakarta Sans, Outfit, JetBrains Mono)
└── views/
    ├── layouts/   # Template kerangka utama HTML (main.templ)
    ├── pages/     # Halaman utuh (dirender saat full-page load) -> misal: pages/auth/login.templ
    └── components/# Potongan UI (Fragments) yang dapat di-reuse atau dikirim sebagai respons HTMX
```

## 3. Typography (Self-Hosted)
Kita menghindari ketergantungan pada Google Fonts demi performa dan privasi. Semua font disimpan di `ui/static/fonts/`.
- **`font-sans`** (Utama): **Plus Jakarta Sans**. Digunakan untuk *body text*, paragraf, label, form, dan mayoritas antarmuka.
- **`font-display`** (Judul): **Outfit**. Digunakan KHUSUS untuk `h1`, `h2`, `h3`, atau teks *hero* yang membutuhkan kesan *bold* dan sangat modern.
- **`font-mono`** (Teknis): **JetBrains Mono**. Digunakan untuk data presisi, kode resi, ID, dan representasi angka tabular.

## 4. Aesthetics & Visual Guidelines
Desain harus selalu terasa **Premium, Elegan, dan Modern (Rich Aesthetics)**. Jangan pernah membuat UI yang terlihat seperti prototipe dasar (MVP). 

Aturan desain yang harus diikuti:
- **Warna Dasar**: Menggunakan `slate` dari Tailwind (misal: `bg-slate-50` untuk *background* utama, `text-slate-900` untuk teks primer, `text-slate-500` untuk teks sekunder).
- **Warna Aksen**: 
  - *Primary*: Deep Navy (`#324376`) -> `bg-primary`, `text-primary`.
  - *Secondary*: Muted Blue (`#586ba4`) -> `bg-secondary`, `text-secondary`.
  - *Accent*: Soft Yellow (`#f5dd90`) -> `bg-accent`, `text-accent`.
  - *Tertiary*: Coral Orange (`#f68e5f`) -> `bg-tertiary`, `text-tertiary`.
  - *Danger/Error*: Coral Red (`#f76c5e`) -> `bg-danger`, `text-danger`.
- **Glassmorphism**: Komponen *card* atau wadah utama sering menggunakan perpaduan latar belakang semi-transparan putih dengan efek *blur*: `bg-white/70 backdrop-blur-xl border border-white/40 shadow-2xl`.
- **Bentuk (Shapes)**: Gunakan sudut yang lebih membulat untuk kesan modern (`rounded-xl`, `rounded-2xl`).
- **Animasi**: Sertakan animasi halus untuk *hover states* (`transition-colors duration-200`) dan *background* bergradasi radial/animasi halus.

## 5. HTMX & Templ Workflow
### Pembuatan Komponen
- Halaman utuh harus selalu dibungkus menggunakan komponen *layout* (contoh: `@layouts.BaseLayout("Judul") { ... }`).
- Pecah komponen yang kompleks menjadi fungsi `templ` kecil di dalam folder `components/`.

### Interaksi HTMX
- Jangan me-reload halaman saat men-submit form! Gunakan `hx-post` dan `hx-target`.
- Saat memproses *request* HTMX di `web_handler.go` yang menghasilkan *error* (seperti validasi), kembalikan **hanya** komponen *fragment* (misalnya `auth_components.LoginError("Pesan")`), bukan keseluruhan halaman.
- Untuk mengarahkan pengguna ke halaman lain setelah form HTMX berhasil, server Go harus mengembalikan *header* HTTP `HX-Redirect: /url-tujuan`.
- Gunakan `.htmx-indicator` bersama elemen SVG loading (spinner) agar tombol terlihat merespon ketika HTMX sedang memanggil server.

## 6. Alpine.js Workflow
- Gunakan Alpine.js murni untuk interaksi yang tidak memerlukan data dari server.
- Contoh: Menutup pesan error: `<div x-data="{ show: true }" x-show="show"> <button @click="show = false">X</button> </div>`.

## 7. Pengembangan (Dev Commands)
Ketika mengembangkan UI dan komponen, Anda hanya perlu menjalankan satu perintah sakti di terminal:
```bash
make dev
```
(Atau `air` secara langsung jika sudah terpasang).
Perintah ini akan secara otomatis memantau file `.go`, `.templ`, dan `.css`, lalu seketika menjalankan ulang kompilasi `templ`, kompilasi `tailwind`, dan menyalakan ulang *server* Go dalam satu alur yang cepat!

## 8. Icons
Untuk ikon antarmuka, kita secara ketat menggunakan *library* **`github.com/dimmerz92/go-icons/lucide`** (Lucide Icons) yang dikemas sebagai komponen `templ`.
- **Dilarang keras** memasukkan tag `<svg>` mentah hasil *copy-paste* dari web (kecuali untuk logo spesifik/ilustrasi *custom* yang tidak ada di Lucide).
- **Format Import**: `import "github.com/dimmerz92/go-icons/lucide"`
- **Cara Penggunaan**: Panggil fungsi komponen ikon dengan menyisipkan atribut `templ.Attributes` untuk menentukan kelas CSS.
  ```go
  @lucide.Mail(templ.Attributes{"class": "w-5 h-5 text-slate-400"})
  ```
- **Kenapa Lucide?**: Memiliki garis (*stroke*) yang membulat, konsisten, dan sangat serasi dengan estetika *Glassmorphism* dan sudut bundar (`rounded-xl`) pada desain proyek ini.
