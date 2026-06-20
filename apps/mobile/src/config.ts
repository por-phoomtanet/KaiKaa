import { Platform } from 'react-native';

// Android emulator เข้าถึง host ผ่าน 10.0.2.2, iOS simulator ใช้ localhost
// แก้เป็น IP เครื่อง backend จริงเมื่อรันบนอุปกรณ์จริง
export const API_BASE_URL =
  Platform.OS === 'android'
    ? 'http://10.0.2.2:8080/api'
    : 'http://localhost:8080/api';
