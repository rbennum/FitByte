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
In go, there are two ways to run the app
## Build
```
# For build, use this
go build

# After build go application
./go-gin-template.exe
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

```

# API Documentation
One the applicatio is running, you can access the Swagger API documentation at:
```
http://localhost:3000/swagger/index.html
```
