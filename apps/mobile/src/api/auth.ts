import { apiClient, unwrap } from './client';
import { setToken, setShop, clearToken, clearShop } from './storage';
import { AuthData } from './types';

export const register = async (body: {
  shop_name: string;
  email: string;
  password: string;
}): Promise<AuthData> => {
  const data = await unwrap<AuthData>(apiClient.post('/auth/register', body));
  await setToken(data.token);
  await setShop(data.shop);
  return data;
};

export const login = async (body: {
  email: string;
  password: string;
}): Promise<AuthData> => {
  const data = await unwrap<AuthData>(apiClient.post('/auth/login', body));
  await setToken(data.token);
  await setShop(data.shop);
  return data;
};

export const logout = async () => {
  await clearToken();
  await clearShop();
};
