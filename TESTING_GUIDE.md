# 📋 Panduan Pengujian Mandiri & Dokumentasi Perubahan
> **Project:** Trello Kanban Suite (`mazkev/react-trello-kanban-suite`) & Clean Architecture Go API (`mazkev/go-pro`)

---

## 🚀 1. Cara Menyalakan Server Sendiri

Buka 2 jendela terminal / PowerShell terpisah:

### Terminal 1: Backend Go API (Port 8080)
```powershell
cd "d:\prog\tes only"
go run ./cmd/api
```
*Tunggu hingga muncul pesan server berjalan di port `:8080` dan koneksi database MySQL/SQLite terhubung.*

### Terminal 2: Frontend React Vite (Port 5173)
```powershell
cd "d:\prog\tes only\frontend"
npm run dev
```
*Buka browser di alamat [http://localhost:5173](http://localhost:5173).*

---

## 🔑 Kredensial Uji Coba
- **Email:** `kevin@example.com`
- **Password:** `rahasia123`

---

## 📦 2. Rangkuman Berkas yang Telah Diubah & Komit GitHub

### A. Backend Go (`mazkev/go-pro`)
Commit: `e1fe304` (*feat: add card comments and activity thread module*) & `fafefd3` (*feat: update profile endpoint*)

| Lapisan Clean Architecture | Berkas | Deskripsi Perubahan |
| :--- | :--- | :--- |
| **Domain** | `internal/domain/board.go` | Menambahkan entity `Comment`, relasi `Comments []Comment` pada `Card`, interface `CommentRepository` & `CommentUsecase`. |
| **Repository** | `internal/repository/comment_repository.go` | Implementasi GORM CRUD komentar (`Create`, `FindByCardID`, `FindByID`, `Delete`). |
| **Repository** | `internal/repository/board_repository.go` | Preload otomatis `Lists.Cards.Comments` agar data komentar ikut terbawa saat board dibuka. |
| **Usecase** | `internal/usecase/comment_usecase.go` | Logika bisnis validasi kartu, author komentar otomatis, dan proteksi hak akses delete. |
| **Delivery / HTTP** | `internal/delivery/http/comment_handler.go` | Controller handler untuk `POST /cards/:cardId/comments`, `GET /cards/:cardId/comments`, dan `DELETE /comments/:id`. |
| **Main Routing** | `cmd/api/main.go` | AutoMigrate tabel `domain.Comment` & registrasi rute API komentar. |
| **REST Testing** | `rest.http` | Test case no. 23 (Profil), 24 (Tambah Komentar), 25 (List Komentar), 26 (Hapus Komentar). |

### B. Frontend React (`mazkev/react-trello-kanban-suite`)
Commit: `611249d` (*feat: card comments & badge*), `ffc08ea` (*refactor: remove port 8080*), `eb3d082` (*feat: navbar dropdowns*), `315cf61` (*fix: smooth sidebar*), `1344e91` (*feat: settings modal*)

| Komponen / Modul | Berkas | Deskripsi Perubahan |
| :--- | :--- | :--- |
| **API Client** | `src/services/api.js` | Helper method `getCardComments`, `createCardComment`, `deleteCardComment`. |
| **Zustand Store** | `src/store/useBoardStore.js` | Action state `addComment` & `deleteComment` terhubung langsung ke backend. |
| **Detail Kartu (Modal)** | `src/components/CardModal.jsx` | Form diskusi dan list komentar nyata dengan avatar, waktu posting, dan tombol hapus komentar milik sendiri. |
| **Pratinjau Kartu** | `src/components/Card.jsx` | Indikator badge ikon pesan + counter jumlah komentar pada kartu. |
| **Sidebar Navigation** | `src/App.jsx` | Animasi toggle buka/tutup sidebar 300ms cubic-bezier yang halus tanpa teks gepeng. |
| **Navbar Utama** | `src/components/Navbar.jsx` | Menu interaktif: *Workspace*, *Recent*, *Starred*, *Templates*, *+ Create*, *Help*, *Profile*. |
| **Modal Pengaturan** | `src/components/SettingsModal.jsx` | Tab update profil (nama & password), manajemen tema board, dan status aplikasi level produksi. |

---

## 🧪 3. Skenario Pengujian Frontend (Browser)

Setelah kedua server dinyalakan dan Anda membuka `http://localhost:5173`:

### Skenario 1: Menguji Fitur Komentar & Diskusi Kartu
1. Klik salah satu kartu tugas di papan Kanban untuk membuka **Card Detail Modal**.
2. Gulir ke bagian **Activity & Comments** di sisi kanan bawah.
3. Ketik komentar di dalam textarea (misal: *"Fitur ini siap diuji sekarang!"*), lalu klik tombol **Comment**.
4. **Hasil yang diharapkan:**
   - Komentar langsung muncul di daftar percakapan dengan nama Anda dan avatar.
   - Komentar tersimpan di backend (database).
5. Tutup modal kartu.
6. **Hasil yang diharapkan:**
   - Di kartu tersebut muncul badge ikon pesan beserta jumlah komentar (misal: `💬 1`).
7. Buka kembali modal kartu, lalu klik ikon tempat sampah (🗑️) pada komentar Anda.
8. **Hasil yang diharapkan:**
   - Komentar terhapus dan badge counter di kartu berkurang otomatis.

---

### Skenario 2: Menguji Animasi Sidebar (Collapse / Expand)
1. Klik tombol panah kiri (`ChevronLeft`) di bagian atas sidebar.
2. **Hasil yang diharapkan:**
   - Sidebar menutup dengan mulus dan halus (durasi 300ms) tanpa teks yang terpotong/gepeng tiba-tiba.
   - Muncul tombol tab melayang (`ChevronRight`) di tepi kiri layar.
3. Klik tombol tab melayang tersebut untuk membuka kembali sidebar.
4. **Hasil yang diharapkan:**
   - Sidebar membuka kembali dengan transisi yang halus.

---

### Skenario 3: Menguji Menu Dropdown Navbar
1. Klik menu **Workspace**: Beralih antar board, lihat jumlah list pada board, atau buat board baru.
2. Klik menu **Starred**: Lihat papan favorit yang ditandai bintang (sinkron dengan tombol bintang di header board).
3. Klik menu **Templates**: Pilih salah satu template (misal: *Agile Sprint* atau *Design Roadmap*) dan klik **Use Template** untuk instan membuat board dengan kolom siap pakai.
4. Klik tombol **+ Create** di navbar: Input judul board dan pilih tema warna gradien.
5. Klik ikon tanda tanya **Help**: Panduan shortcut keyboard dan gesture navigasi.

---

### Skenario 4: Menguji Modal Pengaturan (Settings) & Update Profil
1. Klik ikon **Settings** (⚙️) di header atau dari menu avatar profil.
2. Pada tab **Profile**:
   - Ubah nama Anda (misal: `"Kevin Pratama"`) lalu klik **Save Profile**.
   - Sistem akan memanggil endpoint `PUT /me` di backend dan nama profil terupdate.
3. Pada tab **Preferences**:
   - Ganti wallpaper tema board aktif secara langsung.
4. Pada tab **About**:
   - Pastikan informasi yang ditampilkan bernuansa aplikasi produksi profesional tanpa teks developer port `:8080`.

---

### Skenario 5: Menguji Fitur Kalender Bulanan Interaktif (Calendar View)
1. Di sidebar kiri, klik menu **Calendar**.
2. Periksa tampilan **Grid Bulan**:
   - Grid 7 kolom (Minggu s/d Sabtu) dengan kotak-kotak tanggal.
   - Penanda lingkaran biru untuk hari ini (*Today indicator*).
   - Kartu tugas yang memiliki `dueDate` muncul di kotak tanggal yang sesuai.
3. Klik tombol `<` dan `>` untuk beralih antar bulan, atau klik **Hari Ini**.
4. Klik salah satu kartu di dalam kotak tanggal: modal kartu langsung terbuka!
5. Coba klik toggle **Agenda** untuk melihat tampilan daftar deadline urut waktu.

---

### Skenario 6: Menguji Fitur Upload & Lampiran File (Card Attachments)
1. Buka salah satu kartu di papan Kanban.
2. Klik tombol **Attachment** di sidebar kanan atau tombol **+ Tambah Lampiran**.
3. Pilih file gambar (PNG/JPG) atau dokumen PDF.
4. Periksa:
   - File terunggah ke backend Go (`/uploads/attachments/...`).
   - Thumbnail gambar, nama file, ukuran KB, dan tanggal unggah langsung tampil.
   - Tombol unduh membuka file asli dari server backend.
5. Tutup modal: kartu di board menampilkan badge ikon paperclip (`📎 1`).
6. Buka kembali modal kartu dan klik ikon tempat sampah (🗑️) untuk menghapus lampiran.

---

### Skenario 7: Kolaborasi Multi-User, Role-Based Access Control (RBAC) & Card Assignee
1. **Mengundang Anggota Board (Share Modal):**
   - Pada header board Kanban, perhatikan avatar anggota dan tombol **Share** (ikon `UserPlus`).
   - Klik tombol **Share**. Modal manajemen kolaborator akan terbuka.
   - Masukkan email pengguna terdaftar (misal: `user2@example.com`).
   - Pilih peran (*Role*): `Admin`, `Member`, atau `Viewer`.
   - Klik **Undang**. Anggota baru akan langsung muncul di daftar dengan badge peran.
2. **Mengubah Peran & Menghapus Anggota:**
   - Sebagai pemilik/admin board, Anda dapat mengubah peran anggota via dropdown role (`Admin`, `Member`, `Viewer`).
   - Klik tombol hapus (🗑️) untuk mengeluarkan anggota dari board.
3. **Menugaskan Anggota ke Kartu (Card Assignee):**
   - Buka modal salah satu kartu tugas.
   - Klik tombol **Members** di sidebar aksi kartu.
   - Popover pemilih anggota akan muncul menampilkan seluruh kolaborator board.
   - Klik nama anggota untuk menugaskannya ke kartu tersebut (atau klik *Tanpa Penugasan* untuk melepas).
   - Tutup modal: kartu pada kolom Kanban akan menampilkan avatar assignee di pojok kanan bawah.

---

## ⚡ 4. Skenario Pengujian Backend API (via `rest.http`)

Gunakan ekstensi **REST Client** di IDE pada berkas [`rest.http`](file:///d:/prog/tes%20only/rest.http):

1. **Langkah 1: Login untuk Mendapatkan Token JWT Baru (Request #10)**
   - Jalankan request baris 77 (`POST http://localhost:8080/login`).
   - Salin nilai token dari respon JSON jika token lama sudah kedaluwarsa.
2. **Langkah 2: Tambah Komentar (Request #24)**
   - Jalankan request baris 280 (`POST http://localhost:8080/cards/1/comments`).
   - Respon yang diharapkan: HTTP `201 Created` dengan objek komentar.
3. **Langkah 3: Ambil Seluruh Komentar Kartu (Request #25)**
   - Jalankan request baris 290 (`GET http://localhost:8080/cards/1/comments`).
   - Respon yang diharapkan: HTTP `200 OK` dengan array komentar.
4. **Langkah 4: Hapus Komentar (Request #26)**
   - Jalankan request baris 295 (`DELETE http://localhost:8080/comments/:id`).
   - Respon yang diharapkan: HTTP `200 OK` dengan pesan konfirmasi penghapusan.
5. **Langkah 5: Update Profil (Request #23)**
   - Jalankan request baris 267 (`PUT http://localhost:8080/me`).
   - Respon yang diharapkan: HTTP `200 OK` dengan user profil yang sudah diperbarui.
6. **Langkah 6: Manajemen Kolaborator & RBAC (Request #29 - #33)**
   - **Request #29**: `POST /boards/1/members` (Undang anggota baru dengan role `member`/`viewer`).
   - **Request #30**: `GET /boards/1/members` (Daftar semua anggota dan peran mereka).
   - **Request #31**: `PUT /boards/1/members/:userId` (Ubah role menjadi `admin`, `member`, atau `viewer`).
   - **Request #32**: `DELETE /boards/1/members/:userId` (Keluarkan anggota dari board).
   - **Request #33**: `PUT /cards/1/assign` (Tugaskan user ID tertentu ke kartu).

---

## 🧪 5. Menjalankan Unit Testing Mandiri

Seluruh logika bisnis (Usecase) telah dilengkapi dengan unit test berbasis Table-Driven Tests dan Mock Repository:

### Perintah Menjalankan Unit Test (Terminal/PowerShell):
```powershell
cd "d:\prog\tes only"
go test -v ./internal/usecase/...
```

### Cakupan Unit Test yang Telah Diimplementasikan:
1. **Board Usecase** (`board_usecase_test.go`):
   - Validasi pembuatan board & otomatisasi 3 kolom bawaan (*To Do*, *In Progress*, *Done*).
   - Validasi update board dan proteksi hak milik user.
   - Hapus board.
2. **Board Member & RBAC Usecase** (`board_member_usecase_test.go`):
   - Validasi undang anggota baru dengan email terdaftar dan role valid.
   - Proteksi hak admin (hanya admin/owner yang bisa mengundang dan mengganti role).
   - Proteksi role owner tidak dapat diturunkan/dihapus.
   - Pencegahan duplikasi anggota dalam board yang sama.
3. **Card Usecase** (`card_usecase_test.go`):
   - Validasi pembuatan kartu & default warna label.
   - Perpindahan/geser kartu antar kolom (*move card*).
   - Pembaruan isi, label, cover, dan status checklist kartu.
   - Penugasan anggota ke kartu (*AssignCard*).
4. **Checklist Usecase** (`checklist_usecase_test.go`):
   - Pembuatan sub-tugas dan validasi judul wajib diisi.
   - Toggle centang selesai/belum selesai.
   - Penghapusan sub-tugas.
5. **Comment Usecase** (`comment_usecase_test.go`):
   - Validasi komentar tidak boleh kosong dan kartu harus ada.
   - Resolusi nama author komentar.
   - Proteksi hak akses delete (hanya pemilik yang dapat menghapus).
6. **User Usecase** (`user_usecase_test.go`):
   - Validasi registrasi, deteksi email duplikat, dan enkripsi bcrypt password.
   - Update profil pengguna.

---
*Semua perubahan dan unit test telah terintegrasi dan siap untuk diuji secara mandiri!*

