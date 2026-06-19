# KaiKaa — Requirements

## 1. Project Overview

**KaiKaa** คือแอปพลิเคชัน POS (Point of Sale) สำหรับร้านค้าเล็ก เช่น ร้านกาแฟและร้านเครื่องดื่ม ออกแบบมาเพื่อใช้งานบนมือถือเป็นหลัก พร้อม Web Admin สำหรับจัดการข้อมูลจากหลังบ้าน และระบบ AI ที่วิเคราะห์ยอดขายและให้คำแนะนำประจำวัน

**เฟสแรก (Phase 1):** มุ่งเน้นร้านค้าเล็กที่มีเจ้าของคนเดียว ใช้งานง่าย ไม่ซับซ้อน

**อนาคต (Future Phase):** วางแผนขยายรองรับร้านค้าขนาดกลาง เช่น หลายสาขา, หลาย cashier, ระบบ inventory และ loyalty program

---

## 2. Technology Stack

### Mobile
- React Native
- WebView

### Web Admin
- Next.js
- Tailwind CSS
- Ant Design

### Backend
- Go
- Gin Framework

### Database
- MongoDB

### Authentication
- JWT

### Infrastructure
- Docker
- Docker Compose

### AI
- openrouter

---

## 3. Features

### 3.1 Authentication
- เจ้าของร้านลงทะเบียนด้วย email + password
- เข้าสู่ระบบ → ได้รับ JWT token
- JWT ใช้สำหรับ authorize ทุก API call
- รองรับหลายร้าน (multi-tenant) แต่ละร้านมีข้อมูลแยกกัน

---

### 3.2 Product Management (จัดการสินค้า)

**สินค้าแต่ละรายการประกอบด้วย:**
- ชื่อสินค้า
- ราคา (บาท)
- Emoji icon

**ฟังก์ชัน:**
- ดูรายการสินค้าทั้งหมด
- เพิ่มสินค้าใหม่
- แก้ไขชื่อ / ราคา / emoji
- ลบสินค้า

---

### 3.3 Sales / POS (บันทึกการขาย)

- แสดงสินค้าเป็น Grid 2 คอลัมน์ (emoji + ชื่อ + ราคา)
- กดสินค้า → เปิด Payment Sheet เลือกวิธีชำระเงิน
  - **เงินสด (Cash)**
  - **โอนเงิน (Transfer)**
- บันทึก Sale record: สินค้า, ราคา, วิธีชำระ, เวลา
- แสดง Toast notification หลังบันทึกสำเร็จ

**Sale record fields:**
```
product_id, name, emoji, total, method (cash|transfer), timestamp
```

---

### 3.4 Dashboard (หน้าหลัก)

- **ยอดขายวันนี้** (รวมทุก transaction)
- **จำนวน transaction** วันนี้
- **สัดส่วนเงินสด vs โอนเงิน** (%, progress bar)
- **สินค้าขายดี (Top 4)** — เรียงตามจำนวนแก้ว พร้อม progress bar และยอดรายได้
- แสดง % การเปลี่ยนแปลงเทียบค่าเฉลี่ย

---

### 3.5 Reports (รายงาน)

#### รายวัน (Daily)
- ยอดขายวันนี้
- Bar chart ย้อนหลัง 7 วัน
- รายการล่าสุด 6 รายการ (emoji, ชื่อสินค้า, เวลา, ยอด, วิธีชำระ)

#### รายเดือน (Monthly)
- ยอดขายเดือนนี้
- Bar chart แบ่งรายสัปดาห์ (4 สัปดาห์)
- สถิติสรุป:
  - ยอดเฉลี่ยต่อวัน
  - จำนวนบิลรวม
  - ยอดเฉลี่ยต่อบิล

---

### 3.6 AI Summary (AI สรุปประจำวัน)

เรียก AI (openrouter) วิเคราะห์ข้อมูลยอดขายและส่งออก 3 ส่วน:

1. **สรุปยอดขายวันนี้** — ยอดรวม, จำนวน transaction, % เปลี่ยนแปลง, ช่วงเวลาขายดี, วิธีชำระที่นิยม
2. **สินค้าขายดี** — Top 3 พร้อมจำนวนที่ขาย
3. **คำแนะนำสำหรับร้าน** — 3 ข้อ เช่น จัดโปรโมชัน, เพิ่มสต็อก, กระตุ้นการโอนเงิน

---

### 3.7 Navigation

Bottom Navigation Bar (5 ปุ่ม):

| ไอคอน | หน้า |
|-------|------|
| 🏠 | หน้าหลัก (Dashboard) |
| 📊 | รายงาน (Reports) |
| ＋ (FAB กลาง) | ขายสินค้า (Sales/POS) |
| ✨ | AI สรุป |
| ⊞ | เมนูรวม (Hub) |

---

### 3.8 Web Admin

จัดการจาก browser สำหรับเจ้าของร้าน:
- ดูรายงานยอดขาย (รายวัน / รายเดือน)
- จัดการสินค้าทั้งหมด
- ดู transaction history
- ตั้งค่าร้าน (ชื่อร้าน, ข้อมูลทั่วไป)

---

## 4. Data Models

### Shop
```
_id, owner_id, name, created_at
```

### User
```
_id, email, password_hash, shop_id, created_at
```

### Product
```
_id, shop_id, name, price, emoji, is_active, created_at, updated_at
```

### Sale
```
_id, shop_id, product_id, product_name, product_emoji,
total, method (cash|transfer), sold_at
```

---

## 5. API Endpoints (ภาพรวม)

### Auth
```
POST /api/auth/register
POST /api/auth/login
```

### Products
```
GET    /api/products
POST   /api/products
PUT    /api/products/:id
DELETE /api/products/:id
```

### Sales
```
POST /api/sales
GET  /api/sales?date=YYYY-MM-DD
```

### Reports
```
GET /api/reports/daily?date=YYYY-MM-DD
GET /api/reports/monthly?month=YYYY-MM
```

### AI
```
POST /api/ai/summary
```

---

## 6. UI Theme

- สีหลัก: `#b06a43` (กาแฟน้ำตาลอุ่น)
- สีพื้นหลัง: `#f7efe4` (ครีมอุ่น)
- ฟอนต์: Anuphan (รองรับภาษาไทย)
- มุมโค้งมน (border-radius สูง)
- Bottom Sheet animation สำหรับ Payment และ Product Editor
