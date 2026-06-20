import { apiClient, unwrap } from './client';
import { PaymentMethod, Sale } from './types';

export const createSale = (body: { product_id: string; method: PaymentMethod }) =>
  unwrap<Sale>(apiClient.post('/v1/sales', body));

export const listSales = (date?: string) =>
  unwrap<Sale[]>(apiClient.get('/v1/sales', { params: date ? { date } : undefined }));
