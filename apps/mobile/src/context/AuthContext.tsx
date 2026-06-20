import React, {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from 'react';
import { setUnauthorizedHandler } from '../api/client';
import { clearShop, clearToken, getShop, getToken } from '../api/storage';
import { Shop } from '../api/types';

interface AuthContextValue {
  token: string | null;
  shop: Shop | null;
  initializing: boolean;
  signIn: (token: string, shop: Shop) => void;
  signOut: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [token, setTokenState] = useState<string | null>(null);
  const [shop, setShopState] = useState<Shop | null>(null);
  const [initializing, setInitializing] = useState(true);

  const signIn = useCallback((t: string, s: Shop) => {
    setTokenState(t);
    setShopState(s);
  }, []);

  const signOut = useCallback(async () => {
    await clearToken();
    await clearShop();
    setTokenState(null);
    setShopState(null);
  }, []);

  useEffect(() => {
    // โหลด session ที่ค้างไว้ (เปิดแอปแล้วยัง login อยู่)
    (async () => {
      const [t, s] = await Promise.all([getToken(), getShop()]);
      setTokenState(t);
      setShopState(s);
      setInitializing(false);
    })();

    // token หมดอายุ/401 → client ล้าง token แล้ว เราแค่อัปเดต state ให้เด้งกลับ Login
    setUnauthorizedHandler(() => {
      setTokenState(null);
      setShopState(null);
    });
  }, []);

  return (
    <AuthContext.Provider value={{ token, shop, initializing, signIn, signOut }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextValue => {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return ctx;
};
