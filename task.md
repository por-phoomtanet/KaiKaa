# KaiKaa — Task List (Phase 1: Mobile + API)

> Web Admin เป็น Phase ถัดไป ยังไม่รวมใน scope นี้

---

## Monorepo Structure

```
KaiKaa/
├── apps/
│   ├── mobile/                          # React Native (Phase 1)
│   │   ├── src/
│   │   │   ├── screens/
│   │   │   │   ├── auth/
│   │   │   │   │   ├── LoginScreen.tsx
│   │   │   │   │   └── RegisterScreen.tsx
│   │   │   │   ├── DashboardScreen.tsx
│   │   │   │   ├── SalesScreen.tsx
│   │   │   │   ├── ProductsScreen.tsx
│   │   │   │   ├── ReportsScreen.tsx
│   │   │   │   └── AiScreen.tsx
│   │   │   ├── components/
│   │   │   ├── navigation/
│   │   │   ├── services/                # axios API calls
│   │   │   ├── store/                   # state management
│   │   │   └── theme/                   # colors, fonts
│   │   ├── package.json
│   │   └── app.json
│   │
│   └── web/                             # Next.js (Phase 2 — ยังไม่ทำ)
│
├── services/
│   └── api/                             # Go + Gin
│       ├── cmd/
│       │   └── server/
│       │       └── main.go
│       ├── internal/
│       │   ├── auth/
│       │   │   ├── handler.go
│       │   │   ├── service.go
│       │   │   └── model.go
│       │   ├── product/
│       │   │   ├── handler.go
│       │   │   ├── service.go
│       │   │   └── model.go
│       │   ├── sale/
│       │   │   ├── handler.go
│       │   │   ├── service.go
│       │   │   └── model.go
│       │   ├── report/
│       │   │   ├── handler.go
│       │   │   └── service.go
│       │   ├── ai/
│       │   │   ├── handler.go
│       │   │   └── service.go
│       │   └── middleware/
│       │       └── jwt.go
│       ├── config/
│       │   └── config.go
│       ├── Dockerfile
│       └── go.mod
│
├── docker-compose.yml
├── .env.example
├── requirement.md
└── task.md
```

---

## CONVENTIONS — Pattern ที่ใช้ร่วมกัน (ทำก่อนเขียน feature)

### T00a · API Response Contract
> กำหนดให้ทุก endpoint คืน JSON รูปแบบเดียวกันเสมอ

- [ ] กำหนด Standard Response format ใน `internal/response/response.go`:
  ```json
  // Success
  { "success": true,  "data": <any> }

  // Error
  { "success": false, "message": "error description" }
  ```
- [ ] เขียน helper functions:
  ```go
  response.OK(c, data)
  response.Error(c, http.StatusBadRequest, "message")
  response.Unauthorized(c)
  response.NotFound(c)
  ```
- [ ] ทุก handler ต้องใช้ helper เหล่านี้เท่านั้น ห้าม `c.JSON` โดยตรง

### T00b · API Handler Pattern
> ทุก domain (auth, product, sale, report, ai) ต้องใช้โครงสร้างเดียวกัน

- [ ] แต่ละ domain มี 3 ไฟล์:
  ```
  model.go    — struct MongoDB + request/response DTO
  service.go  — business logic, รับ/คืน (data, error)
  handler.go  — รับ request → validate → เรียก service → response.OK/Error
  ```
- [ ] Handler pattern ทุกตัว:
  ```go
  func (h *Handler) CreateProduct(c *gin.Context) {
      var req CreateProductRequest
      if err := c.ShouldBindJSON(&req); err != nil {
          response.Error(c, 400, err.Error()); return
      }
      data, err := h.service.Create(c, req)
      if err != nil {
          response.Error(c, 500, err.Error()); return
      }
      response.OK(c, data)
  }
  ```
- [ ] Router grouping ใน `cmd/server/main.go`:
  ```
  /api/auth/...       — public
  /api/v1/...         — ต้อง JWT middleware
  ```

### T00c · Mobile API Layer Pattern
> ทุก screen ดึงข้อมูลผ่าน layer เดียวกัน ห้ามเรียก axios ตรงจาก component

- [ ] สร้าง `src/services/apiClient.ts` — axios instance:
  ```ts
  // base URL จาก env, attach JWT header ทุก request อัตโนมัติ
  // interceptor: ถ้า 401 → logout + redirect Login
  ```
- [ ] แต่ละ domain มีไฟล์ service ของตัวเอง:
  ```
  src/services/auth.service.ts
  src/services/product.service.ts
  src/services/sale.service.ts
  src/services/report.service.ts
  src/services/ai.service.ts
  ```
- [ ] ทุก service function มีรูปแบบ:
  ```ts
  export const getProducts = async (): Promise<Product[]> => {
      const res = await apiClient.get('/v1/products')
      return res.data.data
  }
  ```

### T00d · Mobile Screen Pattern
> ทุก screen ใช้โครงสร้าง Loading → Error → Data เหมือนกัน

- [ ] สร้าง custom hook `src/hooks/useApi.ts`:
  ```ts
  const { data, loading, error, refetch } = useApi(() => getProducts())
  // จัดการ loading / error / retry ใน hook เดียว
  ```
- [ ] ทุก screen มีโครงสร้างเดิม:
  ```tsx
  const { data, loading, error, refetch } = useApi(...)

  if (loading) return <LoadingScreen />
  if (error)   return <ErrorScreen onRetry={refetch} />
  return <ScreenContent data={data} />
  ```
- [ ] สร้าง shared components:
  - `src/components/LoadingScreen.tsx` — spinner กลางจอ
  - `src/components/ErrorScreen.tsx` — ข้อความ error + ปุ่ม retry

---

## API CONTRACT — สัญญากลาง (Backend ↔ Mobile ต้องตรงกัน 100%)

> **กฎ:** JSON key ทุกตัวเป็น `snake_case` · เงินเป็น integer (บาท) · เวลาเป็น `HH:mm` · วันที่เป็น `YYYY-MM-DD`
> ทุก response ห่อด้วย `{ "success": ..., "data": ... }` ตาม T00a เสมอ (ตัวอย่างด้านล่างแสดงเฉพาะส่วน `data`)

### Auth

**POST /api/auth/register**
```jsonc
// request
{ "shop_name": "ร้านกาแฟ Brewtiful", "email": "owner@shop.com", "password": "secret123" }
// data
{ "token": "<jwt>", "shop": { "id": "...", "name": "ร้านกาแฟ Brewtiful" } }
```

**POST /api/auth/login**
```jsonc
// request
{ "email": "owner@shop.com", "password": "secret123" }
// data
{ "token": "<jwt>", "shop": { "id": "...", "name": "ร้านกาแฟ Brewtiful" } }
```

### Products

**GET /api/v1/products** → `data` = array
```jsonc
[
  { "id": "...", "name": "ลาเต้", "price": 55, "emoji": "☕" }
]
```

**POST /api/v1/products**
```jsonc
// request
{ "name": "ลาเต้เย็น", "price": 60, "emoji": "🧋" }
// data — product ที่สร้างเสร็จ
{ "id": "...", "name": "ลาเต้เย็น", "price": 60, "emoji": "🧋" }
```

**PUT /api/v1/products/:id** — request เหมือน POST, `data` = product หลังแก้ไข

**DELETE /api/v1/products/:id**
```jsonc
// data
{ "id": "...", "deleted": true }
```

### Sales

**POST /api/v1/sales**
```jsonc
// request
{ "product_id": "...", "method": "cash" }   // method: "cash" | "transfer"
// data — sale record ที่บันทึก (เวลา/ราคา server เป็นคนเติม)
{
  "id": "...", "product_id": "...", "name": "ลาเต้", "emoji": "☕",
  "total": 55, "method": "cash", "time": "09:30", "sold_at": "2026-06-19T09:30:00Z"
}
```

**GET /api/v1/sales?date=YYYY-MM-DD** → `data` = array ของ sale record (รูปแบบเดียวกับด้านบน)

### Reports

**GET /api/v1/reports/daily?date=YYYY-MM-DD**
```jsonc
{
  "total": 23450,
  "count": 22,
  "cash": 12100,
  "transfer": 11350,
  "cash_pct": 52,
  "transfer_pct": 48,
  "change_pct": 12,
  "best_sellers": [
    { "rank": 1, "name": "ลาเต้", "emoji": "☕", "qty": 6, "revenue": 330, "pct_width": 100 }
  ],
  "week_bars": [
    { "day": "จ", "amount": 1820 }, { "day": "อ", "amount": 2140 }
    /* ครบ 7 วัน, ตัวสุดท้าย = วันนี้ */
  ],
  "recent": [
    { "name": "ลาเต้", "emoji": "☕", "time": "09:30", "total": 55, "method": "cash" }
    /* สูงสุด 6 รายการ ล่าสุดก่อน */
  ]
}
```

**GET /api/v1/reports/monthly?month=YYYY-MM**
```jsonc
{
  "month_total": 58460,
  "week_bars": [
    { "label": "สัปดาห์ 1", "amount": 14200 }
    /* ครบ 4 สัปดาห์ */
  ],
  "avg_per_day": 1945,
  "bill_count": 1082,
  "avg_per_bill": 54
}
```

### AI

**POST /api/v1/ai/summary**
```jsonc
// request — ว่าง {} (server ดึงยอดวันนี้จาก DB เอง) หรือระบุ { "date": "YYYY-MM-DD" }
{}
// data
{
  "summary": "วันนี้ทำยอดขายได้ ฿23,450 จาก 22 รายการ เพิ่มขึ้น 12%...",
  "best_sellers": [
    { "name": "ลาเต้", "emoji": "☕", "qty": 6 }
    /* Top 3 */
  ],
  "recommendations": [
    "เพิ่มสต็อกลาเต้และชาไทยเย็นช่วงเช้า",
    "ออกโปรเซ็ตช่วงบ่าย 13:00–15:00",
    "กระตุ้นการโอนด้วยส่วนลด 5 บาท"
  ]
}
```

---

## BACKEND — Go + Gin + MongoDB

### T01 · Project Setup ✅ (ยังไม่ verify runtime)
- [x] สร้าง Monorepo folder structure ตามด้านบน
- [x] Init Go module ใน `services/api/` (`go.mod` เขียนมือ เพราะเครื่องไม่มี Go local)
- [x] ตั้งค่า MongoDB connection ใน `config/config.go` + `config/mongo.go`
- [x] สร้าง `docker-compose.yml` (Go API + MongoDB)
- [x] สร้าง `services/api/Dockerfile` (multi-stage)
- [x] สร้าง `.env.example` (PORT, MONGO_URI, MONGO_DB, JWT_SECRET, OPENROUTER_KEY)
- [x] Scaffold React Native ใน `apps/mobile/` + theme tokens
- **DoD:** `docker compose up` ขึ้นทั้ง API + Mongo · `GET /health` คืน 200 · `cd apps/mobile && npm start` รันได้
  - ⚠️ **ยังไม่ verify** — เครื่องนี้ดึง image จาก Docker Hub ไม่ได้ (network EOF). ต้องรัน `docker compose up -d --build` แล้ว `curl localhost:8080/health` ยืนยันอีกครั้งเมื่อ network ใช้ได้
  - ⚠️ Mobile: เครื่องเป็น Node 16 ต้องอัปเป็น **Node 18+** ก่อน `npm install`

### T02 · Auth API ✅ (verified)
- [x] Data model: `User` + `Shop` (`internal/auth/model.go`)
- [x] `POST /api/auth/register` — สร้าง user + shop
- [x] `POST /api/auth/login` — ตรวจ password, คืน JWT
- [x] JWT middleware (`internal/middleware/jwt.go` — set user_id/shop_id ลง context)
- [x] เพิ่ม: `internal/response` (T00a), `internal/token` (jwt gen/parse), `GET /api/v1/me` (ทดสอบ auth)
- **DoD:** register แล้วได้ token · login ด้วย password ผิดได้ 401 · เรียก `/api/v1/...` โดยไม่มี token ได้ 401 · password ถูก hash (bcrypt) ไม่เก็บ plain text
  - ✅ **ผ่านครบ** — ทดสอบจริง: register 200, email ซ้ำ 409, login ผิด 401, login ถูก 200, /v1/me ไม่มี token 401, มี token 200, hash = `$2a$10$...` (bcrypt)

### T03 · Products API ✅ (verified)
- [x] Data model: `Product` (`internal/product/model.go`)
- [x] `GET /api/v1/products` — ดึงสินค้าทั้งหมดของร้าน
- [x] `POST /api/v1/products` — เพิ่มสินค้าใหม่
- [x] `PUT /api/v1/products/:id` — แก้ไขสินค้า
- [x] `DELETE /api/v1/products/:id` — ลบสินค้า
- [x] Go tests: CRUD + isolation per shop + cross-shop ErrNotFound (8 funcs)
- **DoD:** CRUD ครบ คืน JSON ตรงตาม API Contract · เห็นเฉพาะสินค้าของ shop ตัวเอง (filter ด้วย shop_id จาก token) · แก้/ลบสินค้าร้านอื่นได้ 404
  - ✅ **ผ่านครบ** — HTTP จริง: create/list/update/delete 200, list ร้านอื่น `[]`, update/delete ร้านอื่น 404, name ว่าง 400, ไม่มี token 401

### T04 · Sales API ✅ (verified)
- [x] Data model: `Sale` (`internal/sale/model.go` — snapshot ชื่อ/ราคา/emoji)
- [x] `POST /api/v1/sales` — บันทึก transaction (product_id, method)
- [x] `GET /api/v1/sales?date=YYYY-MM-DD` — ดึง transaction ของวันที่ระบุ
- [x] Go tests: snapshot, cross-shop product, date filter, sort desc, invalid date (7 funcs)
- **DoD:** ขายแล้ว server snapshot ราคา/ชื่อ/emoji + เติม `sold_at`/`time` เอง · `method` รับเฉพาะ `cash`/`transfer` (อื่นได้ 400) · query by date คืนเฉพาะวันนั้น
  - ✅ **ผ่านครบ** — HTTP จริง: ขาย 200 (time เวลาไทย +7), method ผิด 400, product มั่ว 400, date default=วันนี้, วันไม่มีขาย `[]`, format ผิด 400, ไม่มี token 401

### T05 · Reports API ✅ (verified)
- [x] `GET /api/v1/reports/daily?date=YYYY-MM-DD`
  - ยอดรวม, จำนวน transaction, cash/transfer แยก + change_pct (เทียบเฉลี่ย 6 วัน)
  - Top 4 สินค้าขายดี (rank, qty, revenue, pct_width)
  - ยอดขายย้อนหลัง 7 วัน (เติม 0 วันไม่มีขาย, label วันไทย)
  - รายการล่าสุด 6 รายการ
- [x] `GET /api/v1/reports/monthly?month=YYYY-MM`
  - ยอดรวมเดือน, แบ่งรายสัปดาห์ (4 สัปดาห์)
  - ยอดเฉลี่ยต่อวัน, จำนวนบิลรวม, ยอดเฉลี่ยต่อบิล
- [x] Go tests: aggregation, week-fill, pct=100, recent cap 6, monthly buckets, empty, invalid (6 funcs)
- **DoD:** response ตรงตาม API Contract ทุก key · `week_bars` ครบ 7 วัน (เติม 0 วันที่ไม่มีขาย) · `cash_pct + transfer_pct = 100` · วันไม่มีข้อมูลคืนค่า 0 ไม่ใช่ error
  - ✅ **ผ่านครบ** — HTTP จริง: daily/monthly คืนทุก key, pct รวม 100, week_bars 7/4, วันว่าง zeros, date ผิด 400, no-token 401

### T06 · AI Summary API ✅ (verified)
- [x] `POST /api/v1/ai/summary`
  - รับข้อมูลยอดขายวันนี้ (จาก report service)
  - ส่ง prompt ไปยัง openrouter (ผ่าน interface `Completer` → mock ได้)
  - คืนผล 3 ส่วน: สรุปยอด (AI), สินค้าขายดี Top 3 (ข้อมูลจริง), คำแนะนำ 3 ข้อ (AI)
- [x] Go tests: parseAIJSON (plain/fenced/none), happy path, LLM error, invalid date (6 funcs, mock LLM)
- **DoD:** คืน JSON ตรง Contract (`summary`, `best_sellers`, `recommendations`) · ข้อความเป็นภาษาไทย · openrouter ล่ม/timeout คืน error ที่อ่านได้ ไม่ทำ server crash · ไม่ hardcode API key (อ่านจาก env)
  - ✅ **ผ่านครบ** — unit test (mock) ยืนยัน structure + best_sellers จากข้อมูลจริง · HTTP จริง: ไม่มี key → 502 (server ไม่ crash, /health ยัง 200), date ผิด 400, no-token 401
  - ⚠️ e2e กับ openrouter จริงต้องใส่ `OPENROUTER_KEY` ใน `.env` ก่อน (ยังไม่ได้ทดสอบกับ LLM จริง)

---

## MOBILE — React Native

### T07 · Project Setup ✅ (tsc ผ่าน)
- [x] Init React Native project (scaffold + Node 24)
- [x] ติดตั้ง dependencies (navigation, axios, async-storage, toast, chart, bottom-sheet, safe-area)
- [x] ตั้งค่า theme: สี `#b06a43` / `#f7efe4`, ฟอนต์ Anuphan (`src/theme/index.ts`)
- [x] **API layer (T00c):** `src/config.ts`, `src/api/` (types, storage, client + JWT/401 interceptor, auth/products/sales/reports/ai services)
- [x] **useApi hook (T00d):** `src/hooks/useApi.ts` + `LoadingScreen` + `ErrorScreen`
- **DoD:** `npx tsc --noEmit` ผ่าน (0 errors) · รันแอปจริงต้องมี emulator (ไว้ทำตอนมี screens)

### T07b · Frontend Template (Design System) ✅ (tsc ผ่าน)

> อิงจาก prototype `Sales Tracker AI Warm (standalone).html`

**Theme Tokens** — `src/theme/index.ts`
- [x] Colors:
  ```
  primary:    #b06a43   (กาแฟน้ำตาล — header, button, active)
  background: #f7efe4   (ครีมอุ่น — พื้นหลังหลัก)
  surface:    #ffffff   (card, sheet)
  border:     #efe3d4   (เส้นขอบ card)
  textDark:   #4a3526   (ข้อความหลัก)
  textMuted:  #ab9683   (ข้อความรอง)
  cash:       #2f9e6b   (เงินสด — สีเขียว)
  transfer:   #3b6fe0   (โอนเงิน — สีน้ำเงิน)
  ```
- [x] Typography: ฟอนต์ Anuphan, ขนาด 11 / 12 / 13 / 14 / 15 / 17 / 19 / 22 / 26 / 34
- [x] Spacing / Border Radius: radius 9 / 12 / 18 / 20 / 24 / 32

**Reusable Components** — `src/components/`
- [x] `HeaderBar` — พื้นหลัง primary, title + subtitle + optional back button
- [x] `Card` — ขอบมน, border `#efe3d4`
- [x] `ProgressBar` — rounded, รองรับสี primary / cash / transfer
- [x] `EmojiAvatar` — กล่องสี่เหลี่ยมมน แสดง emoji (ปรับ size ได้)
- [x] `PrimaryButton` — พื้นหลัง primary, กด scale
- [x] `GhostButton` — border only, secondary action
- [x] `BottomSheet` — slide up (RN Modal) + drag handle + dim overlay (generic; PaymentSheet/ProductEditorSheet ประกอบตอน T11/T12)
- [x] `Toast` — bottom center, พื้นหลัง `#4a3526`, ✓ icon (auto-dismiss คุมที่ screen)
- [x] `BarChart` — custom bar chart (View ล้วน ไม่พึ่ง library)
- [x] `TabSwitch` — สลับ tab, พื้นหลัง rgba
- ~~`StatusBar`~~ — ข้าม (เป็นแค่ mockup ใน prototype; อุปกรณ์จริงมี status bar ของระบบเอง)

**Layout Templates**
- [x] `ScreenWrapper` — SafeAreaView + background `#f7efe4`
- [x] `SectionTitle` — ข้อความ section
- [x] `Divider` — เส้น `#f3ebe0` ความสูง 1px

### T08 · Auth ✅ (tsc ผ่าน)
- [x] `LoginScreen` — form email + password, เรียก `/api/auth/login`
- [x] `RegisterScreen` — form ชื่อร้าน + email + password, เรียก `/api/auth/register`
- [x] เก็บ JWT (+ shop) ใน AsyncStorage
- [x] Auth guard: `RootNavigator` + `AuthContext` — ไม่มี token → Login
- [x] Logout (ลบ token+shop ผ่าน `signOut`)
- [x] กัน double-submit, แสดง error จาก API, `TextField` component
- **DoD:** login สำเร็จเข้าหน้าหลัก · ปิดแอปเปิดใหม่ยัง login อยู่ (load จาก storage) · token หมดอายุ/401 เด้งกลับ Login อัตโนมัติ (interceptor → setUnauthorizedHandler)
  - ⚠️ verify ด้วย `tsc` (0 errors) — รัน emulator จริงไว้ทำตอน screens ครบ (T14)

### T09 · Bottom Navigation ✅ (tsc ผ่าน)
- [x] Bottom Tab Navigator 5 แท็บ (custom `TabBar`): Dashboard / Reports / Sales / AI / Hub
- [x] ปุ่ม ＋ ขาย กลาง — วงกลมน้ำตาลยกสูง (translateY -14 + shadow)
- [x] Active tab สี `#b06a43`, inactive `#c3b3a2` · ไอคอน SVG (react-native-svg)
- [x] `NavigationContainer` + `MainTabs` + 5 placeholder screens · Hub มี shortcut + logout
- **DoD:** สลับครบ 5 แท็บไม่ค้าง · ปุ่ม ＋ กลางเด่นกว่าแท็บอื่น · แท็บ active สีถูกต้อง
  - ⚠️ verify ด้วย `tsc` (0 errors) — รัน emulator จริงไว้ทำตอน screens ครบ

### T10 · Dashboard Screen ✅ (tsc ผ่าน)
- [x] เรียก `reportApi.dailyReport()` ผ่าน `useApi`
- [x] แสดงยอดขายวันนี้ (ตัวเลขใหญ่ + badge ▲/▼ % เปลี่ยนแปลง + จำนวนรายการ)
- [x] Card เงินสด vs โอนเงิน (% + progress bar สีเขียว/น้ำเงิน)
- [x] รายการ Top 4 สินค้าขายดี (emoji + ชื่อ + progress + ยอด + จำนวน) + empty state
- [x] refresh-on-focus (`useFocusEffect`) — ขายเสร็จกลับมายอดอัปเดต
- **DoD:** ตัวเลขตรงกับ API · ใช้ `useApi` ครบ loading/error/retry (T00d) · ขายสินค้าเสร็จกลับมาหน้านี้ยอดอัปเดต · ร้านยังไม่มีข้อมูลแสดง empty state ไม่ค้าง spinner
  - ⚠️ verify ด้วย `tsc` (0 errors) — รัน emulator จริงไว้ทำตอน screens ครบ

### T11 · Sales / POS Screen
- [ ] เรียก `GET /api/v1/products` แสดง Grid 2 คอลัมน์
- [ ] กดสินค้า → เปิด Bottom Sheet เลือกชำระเงิน
  - ปุ่ม "เงินสด" → เรียก `POST /api/v1/sales` (method: cash)
  - ปุ่ม "โอนเงิน" → เรียก `POST /api/v1/sales` (method: transfer)
- [ ] Toast notification หลังบันทึกสำเร็จ
- [ ] Refresh Dashboard หลังขายสำเร็จ
- **DoD:** กดขายแล้ว record เข้า DB จริง · Toast เด้งหายเองใน 2s · กดรัวๆ ไม่ยิงซ้ำ (กัน double submit) · sheet ปิดหลังขายเสร็จ

### T12 · Products Screen
- [ ] เรียก `GET /api/v1/products` แสดงเป็น List (emoji + ชื่อ + ราคา)
- [ ] ปุ่ม "เพิ่มสินค้า" → Bottom Sheet กรอกชื่อ + ราคา + emoji
  - เรียก `POST /api/v1/products`
- [ ] ปุ่ม ✎ แก้ไข → Bottom Sheet แก้ไข → เรียก `PUT /api/v1/products/:id`
- [ ] ปุ่ม 🗑 ลบ → confirm แล้วเรียก `DELETE /api/v1/products/:id`
- **DoD:** เพิ่ม/แก้/ลบแล้ว list อัปเดตทันที · ลบมี confirm ก่อน · ชื่อว่าง/ราคาไม่ใช่ตัวเลข กดบันทึกไม่ได้

### T13 · Reports Screen
- [ ] Tab สลับ รายวัน / รายเดือน
- [ ] **รายวัน:**
  - ยอดขายวันนี้
  - Bar chart 7 วัน
  - รายการล่าสุด 6 รายการ (emoji, ชื่อ, เวลา, ยอด, วิธีชำระ)
- [ ] **รายเดือน:**
  - ยอดขายเดือนนี้
  - Bar chart แบ่ง 4 สัปดาห์
  - สถิติ: ยอดเฉลี่ยต่อวัน / จำนวนบิล / ยอดเฉลี่ยต่อบิล
- **DoD:** สลับ tab แล้วข้อมูลเปลี่ยนถูก · bar chart สูงสัมพันธ์กับค่าจริง · เงินสด=เขียว โอน=น้ำเงิน ตรง theme

### T14 · AI Summary Screen
- [ ] ปุ่ม "สรุปวันนี้" → เรียก `POST /api/v1/ai/summary`
- [ ] แสดง Loading state ระหว่างรอ AI
- [ ] แสดงผล 3 การ์ด:
  - 📈 สรุปยอดขายวันนี้
  - ⭐ สินค้าขายดี Top 3
  - 💡 คำแนะนำ 3 ข้อ
- **DoD:** กดแล้วเห็น loading ระหว่างรอ · ผลแสดงครบ 3 การ์ด · AI error แสดงข้อความ + ปุ่มลองใหม่ ไม่ค้าง loading

---

## สรุป

| ID | งาน | ประเภท | สถานะ |
|----|-----|--------|-------|
| T01 | Project Setup (Backend) | Backend | ✅ (รอ verify) |
| T02 | Auth API | Backend | ✅ |
| T03 | Products API | Backend | ✅ |
| T04 | Sales API | Backend | ✅ |
| T05 | Reports API | Backend | ✅ |
| T06 | AI Summary API | Backend | ✅ |
| T07 | Project Setup + Template | Mobile | ✅ |
| T08 | Auth (Login / Register) | Mobile | ✅ |
| T09 | Bottom Navigation | Mobile | ✅ |
| T10 | Dashboard Screen | Mobile | ✅ |
| T11 | Sales / POS Screen | Mobile | ⬜ |
| T12 | Products Screen | Mobile | ⬜ |
| T13 | Reports Screen | Mobile | ⬜ |
| T14 | AI Summary Screen | Mobile | ⬜ |

---

> **ลำดับแนะนำ:** T01 → T02 → T03 → T04 → T05 → T06 แล้วค่อย T07 → T08 → T09 → T10 → T11 → T12 → T13 → T14
