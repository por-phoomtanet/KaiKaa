// จัด format เงินบาท — ใส่ comma หลักพันเอง(ไม่พึ่ง Intl ที่ Hermes รองรับไม่ครบ)
export const baht = (n: number): string =>
  '฿' + Math.round(n).toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');

// ทักทายตามช่วงเวลา
export const greeting = (): string => {
  const h = new Date().getHours();
  if (h < 12) {
    return 'สวัสดีตอนเช้า ☕';
  }
  if (h < 17) {
    return 'สวัสดีตอนบ่าย ☀️';
  }
  return 'สวัสดีตอนเย็น 🌙';
};
