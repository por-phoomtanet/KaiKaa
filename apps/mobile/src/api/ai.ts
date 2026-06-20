import { apiClient, unwrap } from './client';
import { AiSummary } from './types';

export const aiSummary = (date?: string) =>
  unwrap<AiSummary>(apiClient.post('/v1/ai/summary', date ? { date } : {}));
