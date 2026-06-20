# KaiKaa ☕

แอปพลิเคชัน **POS (Point of Sale)** สำหรับร้านค้าเล็ก เช่น ร้านกาแฟและร้านเครื่องดื่ม
ใช้งานบนมือถือเป็นหลัก พร้อมระบบ **AI สรุปยอดขาย** และวางแผนต่อยอดเป็น Web Admin

> **Phase 1 (ปัจจุบัน):** Mobile + Backend API เสร็จสมบูรณ์
> **Phase 2 (อนาคต):** Web Admin, หลายสาขา, inventory, loyalty program

📄 รายละเอียด: [requirement.md](requirement.md) · แผนงาน: [task.md](task.md)

---

## ✨ Features

| | ฟีเจอร์ | รายละเอียด |
|---|---|---|
| 🔐 | **Authentication** | สมัคร/เข้าสู่ระบบด้วย JWT · รองรับหลายร้าน (multi-tenant) |
| 🧋 | **จัดการสินค้า** | เพิ่ม / แก้ไข / ลบ สินค้า (ชื่อ, ราคา, emoji) |
| 🛒 | **ขายสินค้า (POS)** | แตะสินค้า → เลือกเงินสด/โอน → บันทึกทันที |
| 🏠 | **Dashboard** | ยอดขายวันนี้, สัดส่วนเงินสด/โอน, สินค้าขายดี |
| 📊 | **รายงาน** | รายวัน (กราฟ 7 วัน) / รายเดือน (4 สัปดาห์ + สถิติ) |
| ✨ | **AI สรุป** | วิเคราะห์ยอดขาย + คำแนะนำ ด้วย openrouter |

---

## 🛠 Tech Stack

| ส่วน | เทคโนโลยี |
|---|---|
| **Mobile** | React Native 0.75, React Navigation, axios, react-native-svg |
| **Backend** | Go 1.22, Gin, MongoDB driver, JWT (golang-jwt), bcrypt |
| **Database** | MongoDB 7 |
| **AI** | openrouter (chat completions) |
| **Infra** | Docker, Docker Compose |
| **CI** | GitHub Actions (go vet + tests + MongoDB service) |
| **Web Admin** | Next.js + Tailwind + Ant Design *(Phase 2)* |

---

## 📁 โครงสร้าง (Monorepo)

```
KaiKaa/
├── apps/
│   ├── mobile/                 React Native (Phase 1)
│   │   └── src/
│   │       ├── api/            API layer (axios + JWT + domain services)
│   │       ├── components/     UI components (Card, BottomSheet, BarChart, …)
│   │       ├── context/        AuthContext
│   │       ├── hooks/          useApi, useToast
│   │       ├── navigation/     Bottom tabs + custom TabBar
│   │       ├── screens/        Dashboard / Sales / Products / Reports / AI / Auth
│   │       └── theme/          Design tokens
│   └── web/                    Next.js (Phase 2 — ยังไม่เริ่ม)
│
├── services/
│   └── api/                    Go + Gin
│       ├── cmd/server/         main.go (entrypoint)
│       ├── internal/
│       │   ├── auth/           register / login / JWT
│       │   ├── product/        CRUD สินค้า
│       │   ├── sale/           บันทึกการขาย
│       │   ├── report/         daily / monthly aggregation
│       │   ├── ai/             สรุปด้วย openrouter
│       │   ├── middleware/     JWT auth
│       │   ├── response/       รูปแบบ response กลาง
│       │   └── token/          JWT generate/parse
│       ├── config/             โหลด env + เชื่อม MongoDB
│       └── Dockerfile
│
├── docker-compose.yml          API + MongoDB
├── .env.example
├── requirement.md
└── task.md
```

---

## 🚀 เริ่มใช้งาน

### สิ่งที่ต้องมี
- **Docker** + Docker Compose (สำหรับ Backend — ไม่ต้องติดตั้ง Go เอง)
- **Node 18+** (สำหรับ Mobile)
- Android SDK / Xcode (สำหรับรันแอปบนอุปกรณ์จริง)

### 1) Backend + Database

```bash
cp .env.example .env          # เติม JWT_SECRET และ OPENROUTER_KEY
docker compose up -d --build

curl localhost:8080/health    # ควรได้ {"status":"ok"}
```

API จะรันที่ `http://localhost:8080`

### 2) Mobile

```bash
cd apps/mobile
npm install                   # ต้องใช้ Node 18+
npm start                     # Metro bundler
npm run android               # หรือ npm run ios
```

> **หมายเหตุ base URL:** แอปตั้งค่าให้ Android emulator เรียก `10.0.2.2:8080`
> และ iOS simulator เรียก `localhost:8080` โดยอัตโนมัติ (ดู [apps/mobile/src/config.ts](apps/mobile/src/config.ts))
> หากรันบนอุปกรณ์จริง ให้แก้เป็น IP ของเครื่องที่รัน backend

---

## 🔌 API Reference

ทุก response ห่อด้วยรูปแบบเดียวกัน:
```jsonc
{ "success": true,  "data": <any> }      // สำเร็จ
{ "success": false, "message": "..." }    // ผิดพลาด
```

| Method | Endpoint | Auth | คำอธิบาย |
|---|---|:---:|---|
| POST | `/api/auth/register` | – | สมัคร (สร้าง user + shop) → คืน JWT |
| POST | `/api/auth/login` | – | เข้าสู่ระบบ → คืน JWT |
| GET | `/api/v1/products` | ✅ | สินค้าทั้งหมดของร้าน |
| POST | `/api/v1/products` | ✅ | เพิ่มสินค้า |
| PUT | `/api/v1/products/:id` | ✅ | แก้ไขสินค้า |
| DELETE | `/api/v1/products/:id` | ✅ | ลบสินค้า |
| POST | `/api/v1/sales` | ✅ | บันทึกการขาย |
| GET | `/api/v1/sales?date=YYYY-MM-DD` | ✅ | รายการขายของวัน |
| GET | `/api/v1/reports/daily?date=YYYY-MM-DD` | ✅ | รายงานรายวัน |
| GET | `/api/v1/reports/monthly?month=YYYY-MM` | ✅ | รายงานรายเดือน |
| POST | `/api/v1/ai/summary` | ✅ | สรุปยอดขายด้วย AI |

> route ที่ขึ้นต้น `/api/v1/` ต้องแนบ header `Authorization: Bearer <token>`

---

## 🧪 Testing

Backend มี unit + integration tests (~35 test functions) ครอบคลุม auth, products, sales, reports, ai

```bash
# รันในเครื่องที่มี Go + MongoDB
cd services/api
MONGO_URI=mongodb://localhost:27017 go test ./...
```

หรือรันผ่าน Docker (ไม่ต้องติดตั้ง Go) — ต่อ network ของ compose เพื่อเข้าถึง MongoDB:

```bash
docker compose up -d mongo
docker run --rm --network kaikaa_default \
  -v "$PWD/services/api:/app" -w /app \
  -e MONGO_URI=mongodb://mongo:27017 \
  golang:1.22-alpine go test ./...
```

> integration tests จะ **skip อัตโนมัติ** ถ้าต่อ MongoDB ไม่ได้ — `go test` จึงรันได้ทุกที่
> CI (GitHub Actions) รัน `go vet` + tests พร้อม MongoDB service ทุก push ที่แตะ `services/api/`

---

## ⚙️ Environment Variables

| ตัวแปร | ค่าเริ่มต้น | คำอธิบาย |
|---|---|---|
| `PORT` | `8080` | พอร์ต API |
| `MONGO_URI` | `mongodb://mongo:27017` | MongoDB connection |
| `MONGO_DB` | `kaikaa` | ชื่อ database |
| `JWT_SECRET` | – | secret สำหรับเซ็น JWT (**ต้องเปลี่ยน**) |
| `OPENROUTER_KEY` | – | API key ของ openrouter ([ขอที่นี่](https://openrouter.ai/keys)) |
| `OPENROUTER_MODEL` | `openai/gpt-4o-mini` | model ที่ใช้สรุป |

---

## 📌 สถานะโปรเจกต์

| | Backend | Mobile |
|---|---|---|
| **สถานะ** | ✅ เสร็จ + verified (HTTP + tests) | ✅ เสร็จ (type-check ผ่าน) |
| **เหลือ** | – | รันบน emulator/อุปกรณ์จริง |

**ยังไม่ทำ (Phase 2):** Web Admin (Next.js), หลายสาขา, inventory, loyalty
