// Design tokens — อิงจาก prototype "Sales Tracker AI Warm"
// ใช้เป็นแหล่งความจริงเดียว (single source of truth) สำหรับสี/ฟอนต์/ระยะ

export const colors = {
  primary: '#b06a43', // กาแฟน้ำตาล — header, button, active
  background: '#f7efe4', // ครีมอุ่น — พื้นหลังหลัก
  surface: '#ffffff', // card, sheet
  border: '#efe3d4', // เส้นขอบ card
  divider: '#f3ebe0', // เส้นคั่นในลิสต์
  textDark: '#4a3526', // ข้อความหลัก
  textMuted: '#ab9683', // ข้อความรอง
  cash: '#2f9e6b', // เงินสด — เขียว
  transfer: '#3b6fe0', // โอนเงิน — น้ำเงิน
} as const;

export const fontSize = {
  xs: 11,
  sm: 12,
  base: 13,
  md: 14,
  lg: 15,
  xl: 17,
  xxl: 19,
  title: 22,
  display: 26,
  hero: 34,
} as const;

export const radius = {
  sm: 9,
  md: 12,
  lg: 18,
  xl: 20,
  xxl: 24,
  pill: 32,
} as const;

export const spacing = {
  xs: 4,
  sm: 8,
  md: 12,
  lg: 16,
  xl: 20,
  xxl: 24,
} as const;

export const fontFamily = 'Anuphan';

export const theme = { colors, fontSize, radius, spacing, fontFamily } as const;

export type Theme = typeof theme;
