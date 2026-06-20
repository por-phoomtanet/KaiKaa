import React, { useState } from 'react';
import { LoginScreen } from './LoginScreen';
import { RegisterScreen } from './RegisterScreen';

// สลับ Login/Register ด้วย local state (ไม่ต้องใช้ navigator แยกสำหรับ auth)
export const AuthScreen: React.FC = () => {
  const [mode, setMode] = useState<'login' | 'register'>('login');

  return mode === 'login' ? (
    <LoginScreen onSwitch={() => setMode('register')} />
  ) : (
    <RegisterScreen onSwitch={() => setMode('login')} />
  );
};
