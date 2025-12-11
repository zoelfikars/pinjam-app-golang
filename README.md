Tentu, berikut adalah kerangka *README.md* yang disesuaikan dengan fokus pada **Goose Migration** dan konfigurasi PostgreSQL serta JWT Secret, berdasarkan detail yang Anda berikan.

-----

# 🚀 [NAMA PROYEK ANDA] Backend Service (Go & PostgreSQL)

[Tambahkan deskripsi singkat proyek Go Anda di sini. Jelaskan tujuan utama service ini (misalnya, API untuk Pengajuan Pinjaman, dll.).]

## 📋 Daftar Isi

1.  [Prasyarat](https://www.google.com/search?q=%231-prasyarat)
2.  [Instalasi dan Setup](https://www.google.com/search?q=%232-instalasi-dan-setup)
3.  [Konfigurasi `.env`](https://www.google.com/search?q=%233-konfigurasi-env)
4.  [Manajemen Database (Goose)](https://www.google.com/search?q=%234-manajemen-database-goose)
5.  [Menjalankan Proyek](https://www.google.com/search?q=%235-menjalankan-proyek)
6.  [Dokumentasi API (Postman)](https://www.google.com/search?q=%236-dokumentasi-api-postman)

-----

## 1\. Prasyarat

Pastikan Anda telah menginstal perangkat lunak berikut:

  * **Go:** Versi Go 1.21 atau lebih tinggi.
  * **PostgreSQL:** Server database PostgreSQL yang sedang berjalan.
  * **Git:** Untuk mengkloning repositori.

## 2\. Instalasi dan Setup

### 2.1. Clone Repositori

Buka Terminal atau Command Prompt Anda dan jalankan:

```bash
git clone https://www.andarepository.com/
cd [NAMA DIREKTORI PROYEK]
```

### 2.2. Install Goose (Database Migration Tool)

Kita menggunakan Goose untuk mengelola skema database. Instal Goose secara global:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### 2.3. Install Dependencies Go

Unduh semua modul Go yang dibutuhkan proyek:

```bash
go mod tidy
```

## 3\. Konfigurasi `.env`

Aplikasi memerlukan konfigurasi melalui file variabel lingkungan `.env`.

### 3.1. Buat File `.env`

Buat salinan dari file contoh (`.env.example`) dan beri nama `.env` di *root* proyek:

```bash
cp .env.example .env
```

### 3.2. Contoh Isi File `.env`

Pastikan Anda menyesuaikan nilai-nilai berikut dengan lingkungan lokal atau *staging* Anda:

```ini
DB_HOST=localhost
DB_PORT=5432
DB_USER=db_user
DB_PASSWORD='password'
DB_NAME=db_name
SERVER_PORT=8095
JWT_SECRET=c30cd2b98764e05c2236bb1548210d99
REDIS_ADDR=localhost:6379
REDIS_DB=0

# Konfigurasi Goose
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://db_user:password@localhost:5432/db_name
GOOSE_MIGRATION_DIR=./database/migration
```

### 3.3. Membuat `JWT_SECRET` Baru

Kunci JWT (JSON Web Token) harus unik dan rahasia. Anda dapat membuatnya secara acak menggunakan perintah *shell* (di Linux/macOS):

```bash
# Menghasilkan 32 karakter hexadecimal acak
openssl rand -hex 32
```

Salin hasil *output* dari perintah di atas dan gunakan untuk mengisi variabel `JWT_SECRET` di file `.env` Anda.

## 4\. Manajemen Database (Goose)

Sebelum menjalankan aplikasi, Anda harus menyiapkan dan memigrasikan skema database.

### 4.1. Migrasi Database

Asumsikan Anda telah membuat database bernama `pengajuan_app_db` di PostgreSQL:

Gunakan `goose up` untuk menjalankan semua migrasi yang tertunda yang ada di direktori `GOOSE_MIGRATION_DIR`:

```bash
goose up
```

Perintah ini akan membaca konfigurasi `GOOSE_DRIVER` dan `GOOSE_DBSTRING` dari lingkungan shell Anda (atau file `.env` jika dikonfigurasi) dan menerapkan skema terbaru.

### 4.2. Perintah Goose Lain yang Berguna

| Perintah | Deskripsi |
| :--- | :--- |
| `goose status` | Melihat status semua migrasi (sudah atau belum dijalankan). |
| `goose down` | Membalikkan migrasi terakhir yang diterapkan (rollback). |
| `goose create [nama_migrasi] sql` | Membuat file migrasi baru. |

## 5\. Menjalankan Proyek

Setelah database siap, Anda dapat menjalankan *service* backend Go:

```bash
go run cmd/main.go
```

(Asumsikan file utama Anda adalah `cmd/main.go`)

Aplikasi akan berjalan di port yang ditentukan di `SERVER_PORT` (default: `http://localhost:8095`).

## 6\. Dokumentasi API (Postman)

Koleksi Postman untuk menguji semua *endpoint* API tersedia di direktori `[PATH/TO/POSTMAN/COLLECTION.json]`.

### 6.1. Cara Import Koleksi Postman

Ikuti langkah-langkah ini untuk mengimpor dan menjalankan API:

1.  **Buka Postman:** Luncurkan aplikasi Postman Anda.
2.  **Pilih "Import":** Klik tombol **"Import"** yang berada di sudut kiri atas Postman.
3.  **Pilih File:** Pilih opsi **"Choose Files"** (atau tarik & lepas file) dan arahkan ke file koleksi Postman di dalam proyek Anda:
    ```
    [PATH/TO/POSTMAN/COLLECTION.json]
    ```
4.  **Atur Environment:** Setelah koleksi terimpor, Anda mungkin perlu mengatur *environment* untuk menyesuaikan URL dasar dan *token* otentikasi.
      * Cek variabel `baseURL` di *Environment* Postman dan pastikan nilainya adalah `http://localhost:[SERVER_PORT]` (misalnya: `http://localhost:8095`).

### 6.2. Penggunaan Token JWT

  * Untuk *endpoint* yang memerlukan otentikasi, Anda harus menjalankan *endpoint* login terlebih dahulu.
  * Salin `access_token` dari respons login.
  * Tempelkan *token* tersebut ke bagian **Authorization** (biasanya menggunakan tipe **Bearer Token**) pada *request* API yang Anda uji.