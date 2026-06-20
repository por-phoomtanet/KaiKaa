import { apiClient, unwrap } from './client';
import { Product, ProductInput } from './types';

export const listProducts = () => unwrap<Product[]>(apiClient.get('/v1/products'));

export const createProduct = (body: ProductInput) =>
  unwrap<Product>(apiClient.post('/v1/products', body));

export const updateProduct = (id: string, body: ProductInput) =>
  unwrap<Product>(apiClient.put(`/v1/products/${id}`, body));

export const deleteProduct = (id: string) =>
  unwrap<{ id: string; deleted: boolean }>(apiClient.delete(`/v1/products/${id}`));
