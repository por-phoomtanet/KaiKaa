import { apiClient, unwrap } from './client';
import { DailyReport, MonthlyReport } from './types';

export const dailyReport = (date?: string) =>
  unwrap<DailyReport>(apiClient.get('/v1/reports/daily', { params: date ? { date } : undefined }));

export const monthlyReport = (month?: string) =>
  unwrap<MonthlyReport>(
    apiClient.get('/v1/reports/monthly', { params: month ? { month } : undefined }),
  );
