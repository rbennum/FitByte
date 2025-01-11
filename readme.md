# Go PS-3-T
Cara run?
1. `go run main.go`
2. kalau mau tes--misalnya--. Buka `localhost:port/swagger/index.html`

# Installation
```
# clone the repository
git clone [this_git_url]

# set up environment
cp .env.example .env

# go run
go run main.go
```

# Configuration
Create a `.env` file in the root directory with the following variables:
```
DB_HOST=localhost
DB_USER=user
DB_PASSWORD=pw
DB_NAME=ps3t
DB_PORT=5432 #Example: postgresql
PORT=3000
MODE=DEBUG
PROD_HOST=#Your production host
DEBUG_HOST=0.0.0.0
```

# Running the App

In Go, there are two ways to run the app

## Build

```go
# For build, run this command
go build -o .build/<name-of-build.extension>

# NOTE: it is important to put the build inside of the .build folder
# to ensure the gitignore caught up with the files

# After build go application
.build/<name-of-build.extension>
```

## Go Run
```
# just do this
go run main.go

# then your operating system asking for firewall permission
```

## On Docker

```
docker-compose up -d
# Dengan flag -d untuk menjalankan container di background (detached mode).
## Dengan menggunakan -d, terminal akan langsung kembali ke prompt tanpa menampilkan log container di terminal.
## Container akan terus berjalan di background setelah perintah ini dijalankan.

# Kalau mau tambah flag --build jika ada perubahan pada Dockerfile
docker-compose up -d --build
## Flag --build digunakan untuk memaksa Docker Compose membangun ulang image sebelum menjalankan container.
## Biasanya digunakan ketika ada perubahan pada Dockerfile atau file yang terkait dengan image.
## Perintah ini akan melakukan build ulang image dan kemudian menjalankan container di background.

# Jika ingin menggunakan container database postgres di local bisa menggunakan argument berikut: --profile local
contoh: docker compose --profile local up -d --build
```

# API Documentation
One the applicatio is running, you can access the Swagger API documentation at:
```
http://localhost:3000/swagger/index.html
```

# Use of Dependency Injection
Caranya adalah
1. Setup dari dependensi dasar sebuah service yang sekiranya tidak membutuhkan dependensi service lain, bisa cek pada `di/injector.go`
2. Siapkan constructor dari sebuah service
3. Siapkan injector dengan memanggilkan `MustInvoke` pada servis yang di-dependensi
4. Kembali pada `di/injector.go` dan lakukan dengan `do.Provide` diurutan terbawah setelah service yang di-dependensikan
5. Jika service tersebut dibutukan pada `router` lakukan dengan cara seperti ini, contoh 
```go
authHandler := do.MustInvoke[authHandler.AuthorizationHandler](di.Injector)
```

# Use Migrations
 cara melakukan migrasi database menggunakan `golang-migrate` dengan opsi manual atau otomatis melalui konfigurasi di file `.env`.


## Prasyarat
1. Pastikan Anda telah menginstall CLI `golang-migrate`.
2. File `.env` memiliki konfigurasi berikut:

```shell
ENABLE_AUTO_MIGRATION=FALSE
```

## Manual

Jika `ENABLE_AUTO_MIGRATION` diset ke `FALSE`, migrasi perlu dijalankan secara manual menggunakan CLI.

### Mode Manual Migration

1. Migration Up
        
    Gunakan perintah berikut untuk menjalankan semua migrasi ke versi terbaru:

    ```shell
    migrate -database "postgres://username:password@localhost:5432/dbname?sslmode=disable" -path database/migrations up
    ```

    Gantilah `username`, `password`, `localhost`, `5432`, dan `dbname `dengan informasi koneksi database Anda.

2. Migration Down
        
    Gunakan perintah berikut untuk menurunkan versi migrasi satu langkah ke bawah:

    ```shell
    migrate -database "postgres://username:password@localhost:5432/dbname?sslmode=disable" -path database/migrations down
    ```

    Gantilah `username`, `password`, `localhost`, `5432`, dan `dbname `dengan informasi koneksi database Anda.

3. Membuat File Migration Baru
        
    Untuk membuat file migrasi baru, gunakan perintah berikut:

    ```shell
    migrate create -ext sql -dir database/migrations <NameMigrationFile>   
    ```

## Mode Auto Migration
Jika `ENABLE_AUTO_MIGRATION` diset ke `TRUE`, migrasi akan dijalankan secara otomatis ketika aplikasi dijalankan.