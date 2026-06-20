import AsyncStorage from '@react-native-async-storage/async-storage';
import { Shop } from './types';

const TOKEN_KEY = 'kaikaa.token';
const SHOP_KEY = 'kaikaa.shop';

export const getToken = () => AsyncStorage.getItem(TOKEN_KEY);
export const setToken = (token: string) => AsyncStorage.setItem(TOKEN_KEY, token);
export const clearToken = () => AsyncStorage.removeItem(TOKEN_KEY);

export const setShop = (shop: Shop) => AsyncStorage.setItem(SHOP_KEY, JSON.stringify(shop));
export const getShop = async (): Promise<Shop | null> => {
  const raw = await AsyncStorage.getItem(SHOP_KEY);
  return raw ? (JSON.parse(raw) as Shop) : null;
};
export const clearShop = () => AsyncStorage.removeItem(SHOP_KEY);
