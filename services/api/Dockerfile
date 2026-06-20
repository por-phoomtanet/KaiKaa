# ---- build stage ----
FROM golang:1.22-alpine AS build
WORKDIR /app

# download deps ก่อน (cache layer นี้ตราบใดที่ go.mod/go.sum ไม่เปลี่ยน)
COPY go.mod go.sum ./
RUN go mod download

# copy source แล้ว build (go.sum มีอยู่แล้ว → reproducible)
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# ---- run stage ----
FROM alpine:3.20
WORKDIR /
COPY --from=build /server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
