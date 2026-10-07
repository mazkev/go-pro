# ==========================================
# STAGE 1: BUILDER (Tempat Memasak Kode)
# ==========================================
FROM golang:1.23-alpine AS builder

# Set folder kerja di dalam container
WORKDIR /app

# Salin file dependensi dulu agar memanfaatkan cache Docker
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh source code proyek
COPY . .

# Kompilasi kode menjadi binary Linux statis
# CGO_ENABLED=0 membuat binary 100% mandiri tanpa ketergantungan library C eksternal
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api-service ./cmd/api

# ==========================================
# STAGE 2: FINAL RUNNER (Hanya Wadah Eksekusi Super Ringan)
# ==========================================
FROM alpine:latest

WORKDIR /app

# Salin HANYA file binary dari stage builder di atas
COPY --from=builder /app/api-service .

# Port yang dibuka
EXPOSE 8080

# Perintah untuk menyalakan aplikasi
CMD ["./api-service"]
