import axios, { AxiosError } from 'axios';
import { API_BASE_URL } from '../config';
import { getToken, clearToken } from './storage';

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 15000,
});

// แนบ JWT ทุก request อัตโนมัติ
apiClient.interceptors.request.use(async (config) => {
  const token = await getToken();
  if (token) {
    (config.headers as Record<string, string>).Authorization = `Bearer ${token}`;
  }
  return config;
});

// callback ให้ navigation layer ลงทะเบียน — ใช้เด้งกลับ Login เมื่อ token หมดอายุ
type UnauthorizedHandler = () => void;
let unauthorizedHandler: UnauthorizedHandler = () => {};
export const setUnauthorizedHandler = (fn: UnauthorizedHandler) => {
  unauthorizedHandler = fn;
};

// 401 → ล้าง token + แจ้ง navigation
apiClient.interceptors.response.use(
  (res) => res,
  async (error: AxiosError) => {
    if (error.response?.status === 401) {
      await clearToken();
      unauthorizedHandler();
    }
    return Promise.reject(error);
  },
);

interface Envelope<T> {
  success: boolean;
  data: T;
  message?: string;
}

// แกะ envelope { success, data } → คืนเฉพาะ data
export async function unwrap<T>(req: Promise<{ data: Envelope<T> }>): Promise<T> {
  const res = await req;
  return res.data.data;
}

// ดึงข้อความ error ที่อ่านได้จาก response (ตาม contract { message })
export function errorMessage(e: unknown): string {
  const err = e as AxiosError<{ message?: string }>;
  return err.response?.data?.message ?? err.message ?? 'เกิดข้อผิดพลาด';
}
