# ⚡️ [APLIKASI-PENGAJUAN-PINJAMAN] Backend Service

## 1\. Setup & Run

### 1.1. Prasyarat

  * **Go** (1.25.4)
  * **PostgreSQL** (18.1.1)
  * **Postman** (11.75.3)

### 1.2. Instalasi

```bash
git clone https://github.com/zoelfikars/pinjam-app-golang.git
cd [APLIKASI-PENGAJUAN-PINJAMAN]
go mod tidy
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### 1.3. Konfigurasi

1.  Buat file `.env`: `cp .env.example .env`
2.  Isi kredensial DB dan set `SERVER_PORT=8095`.
3.  *Generate* `JWT_SECRET`: `openssl rand -hex 32` (bisa juga generate jwt dari website online)

### 1.4. Database & Run

1.  Migrasi DB: `goose up`
2.  Jalankan Service: `go run cmd/main.go` (otomatis seeder)

## 2\. Testing API (Postman)

1.  **Import** *Collection* Postman.
2.  Di **"Variables"** *Collection*:
      * Set `domain` menjadi `http://localhost:8095`.
      * *Login* Admin/Nasabah, *copy* token, dan *paste* ke variabel `admin` dan `nasabah`. sesuaikan dengan role (untuk mempermudah testing)
3.  Jalankan seluruh *Collection* untuk *testing* API.