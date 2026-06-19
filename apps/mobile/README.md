# KaiKaa Mobile

React Native app (Phase 1).

## ข้อกำหนด
- **Node >= 18** (เครื่องนี้ตอน scaffold เป็น Node 16 — ต้องอัปเกรดก่อนรันจริง)
- JDK 17 + Android SDK (สำหรับ Android)
- Xcode (สำหรับ iOS)

## เริ่มใช้งาน
```bash
cd apps/mobile
npm install        # ต้องใช้ Node 18+
npm start          # metro bundler
npm run android    # หรือ npm run ios
```

## โครงสร้าง
```
src/
├── screens/      # หน้าจอแต่ละหน้า (T08, T10–T14)
├── components/   # reusable components (T07b)
├── navigation/   # bottom tabs (T09)
├── services/     # API layer (T00c)
├── hooks/        # useApi ฯลฯ (T00d)
├── store/        # state
└── theme/        # design tokens ✓ พร้อมแล้ว
```
