import { apiClient, unwrap } from './client';
import { setToken, clearToken } from './storage';
import { AuthData } from './types';

export const register = async (body: {
  shop_name: string;
  email: string;
  password: string;
}): Promise<AuthData> => {
  const data = await unwrap<AuthData>(apiClient.post('/auth/register', body));
  await setToken(data.token);
  return data;
};

export const login = async (body: {
  email: string;
  password: string;
}): Promise<AuthData> => {
  const data = await unwrap<AuthData>(apiClient.post('/auth/login', body));
  await setToken(data.token);
  return data;
};

export const logout = () => clearToken();
