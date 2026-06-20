// Types ตรงตาม API Contract (snake_case) — แหล่งความจริงเดียวฝั่ง mobile

export type PaymentMethod = 'cash' | 'transfer';

export interface Shop {
  id: string;
  name: string;
}

export interface AuthData {
  token: string;
  shop: Shop;
}

export interface Product {
  id: string;
  name: string;
  price: number;
  emoji: string;
}

export interface ProductInput {
  name: string;
  price: number;
  emoji: string;
}

export interface Sale {
  id: string;
  product_id: string;
  name: string;
  emoji: string;
  total: number;
  method: PaymentMethod;
  time: string;
  sold_at: string;
}

export interface BestSeller {
  rank: number;
  name: string;
  emoji: string;
  qty: number;
  revenue: number;
  pct_width: number;
}

export interface DayBar {
  day: string;
  amount: number;
}

export interface RecentSale {
  name: string;
  emoji: string;
  time: string;
  total: number;
  method: PaymentMethod;
}

export interface DailyReport {
  total: number;
  count: number;
  cash: number;
  transfer: number;
  cash_pct: number;
  transfer_pct: number;
  change_pct: number;
  best_sellers: BestSeller[];
  week_bars: DayBar[];
  recent: RecentSale[];
}

export interface WeekBar {
  label: string;
  amount: number;
}

export interface MonthlyReport {
  month_total: number;
  week_bars: WeekBar[];
  avg_per_day: number;
  bill_count: number;
  avg_per_bill: number;
}

export interface AiBestSeller {
  name: string;
  emoji: string;
  qty: number;
}

export interface AiSummary {
  summary: string;
  best_sellers: AiBestSeller[];
  recommendations: string[];
}
