# e-meeting-api

## Pendahuluan
**e-meeting-api** adalah sebuah REST API untuk manajemen meeting berbasis web.

## Persyaratan
- Go (minimal versi 1.20)
- Git
- Database (sesuai konfigurasi di `config.json`)

## Cara Clone dan Menjalankan

1. Clone repository:
   ```sh
   git clone https://github.com/agung-margi/e-meeting-api.git
   ```
2. Masuk ke folder proyek:
   ```sh
   cd e-meeting-api/
   ```
3. Pastikan file konfigurasi tersedia:
   Buat file `configs/config.json` berdasarkan contoh berikut:
   ```json
   {
   "database": {
    "host": "",
    "port": "",
    "username": "",
    "password": "",
    "dbname": ""
   },
   "jwt_secret_key": "",
   "base_url": "",
   "port": "",
   "base_url_fe": "",
   "photos_path": "./public/photos/",
   "email_config": {
      "host": "",
      "port": "",
      "username": "",
      "password": ""
     }
   }
   ```
4. Jalankan aplikasi:
   ```sh
   go run cmd/main.go
   ```
5. Akses API melalui Swagger:
   Buka di browser:
   ```
   http://localhost:8080/api/v1/swagger/index.html
   ```

## Struktur Proyek
```
├── cmd/
│   ├── main.go
├── configs/
│   ├── config.json
├── internal/
│   ├── domain/
│   ├── usecase/
├── pkg/
│   ├── database/
│   ├── middleware/
│   ├── util/
├── presenter/
│   ├── handler/
│   ├── model/
├── public/
│   ├── photos/
├── go.mod
└── go.sum
```