# ---- build stage ----
FROM golang:1.22-alpine AS build
WORKDIR /app

# copy source ก่อน แล้วค่อย tidy (สร้าง go.sum) + build
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# ---- run stage ----
FROM alpine:3.20
WORKDIR /
COPY --from=build /server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
