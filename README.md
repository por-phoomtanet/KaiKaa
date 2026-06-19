# KaiKaa

แอป POS สำหรับร้านค้าเล็ก (ร้านกาแฟ/เครื่องดื่ม) — มือถือเป็นหลัก + Web Admin + AI สรุปยอดขาย

ดูรายละเอียดที่ [requirement.md](requirement.md) และแผนงานที่ [task.md](task.md)

## Tech Stack
- **Mobile:** React Native + WebView
- **Web Admin:** Next.js + Tailwind + Ant Design (Phase 2)
- **Backend:** Go + Gin
- **Database:** MongoDB
- **Auth:** JWT
- **AI:** openrouter
- **Infra:** Docker + Docker Compose

## โครงสร้าง (Monorepo)
```
apps/
  mobile/        React Native (Phase 1)
  web/           Next.js (Phase 2)
services/
  api/           Go + Gin
docker-compose.yml
```

## เริ่มใช้งาน

### Backend + Database
```bash
cp .env.example .env        # เติมค่า JWT_SECRET, OPENROUTER_KEY
docker compose up -d --build
curl localhost:8080/health  # ควรได้ {"status":"ok"}
```

### Mobile (ต้องใช้ Node 18+)
```bash
cd apps/mobile
npm install
npm start
npm run android   # หรือ npm run ios
```
