import React, { useCallback, useRef, useState } from 'react';
import { ScrollView, StyleSheet, Text, View } from 'react-native';
import { useFocusEffect } from '@react-navigation/native';
import { reportApi } from '../api';
import type { Bar } from '../components';
import {
  BarChart,
  Card,
  EmojiAvatar,
  ErrorScreen,
  LoadingScreen,
  ScreenWrapper,
  TabSwitch,
} from '../components';
import { useApi } from '../hooks/useApi';
import { colors, fontSize, radius } from '../theme';
import { baht } from '../utils/format';

type TabKey = 'daily' | 'monthly';

export const ReportsScreen: React.FC = () => {
  const [tab, setTab] = useState<TabKey>('daily');
  const daily = useApi(() => reportApi.dailyReport());
  const monthly = useApi(() => reportApi.monthlyReport());

  // refresh เมื่อกลับมาโฟกัส (ข้าม mount แรก)
  const firstFocus = useRef(true);
  useFocusEffect(
    useCallback(() => {
      if (firstFocus.current) {
        firstFocus.current = false;
        return;
      }
      daily.refetch();
      monthly.refetch();
    }, [daily, monthly]),
  );

  return (
    <ScreenWrapper>
      <View style={styles.header}>
        <Text style={styles.title}>รายงานยอดขาย</Text>
        <TabSwitch
          tabs={[
            { key: 'daily', label: 'รายวัน' },
            { key: 'monthly', label: 'รายเดือน' },
          ]}
          active={tab}
          onChange={(k) => setTab(k as TabKey)}
        />
      </View>

      {tab === 'daily' ? <DailyView state={daily} /> : <MonthlyView state={monthly} />}
    </ScreenWrapper>
  );
};

// ===== Daily =====
const DailyView: React.FC<{ state: ReturnType<typeof useApi<Awaited<ReturnType<typeof reportApi.dailyReport>>>> }> = ({
  state,
}) => {
  const { data, loading, error, refetch } = state;
  if (loading && !data) {
    return <LoadingScreen />;
  }
  if (error || !data) {
    return <ErrorScreen message={error ?? undefined} onRetry={refetch} />;
  }

  const bars: Bar[] = data.week_bars.map((w, i) => ({
    label: w.day,
    value: w.amount,
    highlight: i === data.week_bars.length - 1,
  }));

  return (
    <ScrollView contentContainerStyle={styles.body} showsVerticalScrollIndicator={false}>
      <Card>
        <Text style={styles.cardLabel}>ยอดขายวันนี้</Text>
        <Text style={styles.cardValue}>{baht(data.total)}</Text>
        <View style={styles.chart}>
          <BarChart bars={bars} />
        </View>
      </Card>

      <Text style={styles.sectionTitle}>รายการล่าสุด</Text>
      {data.recent.length === 0 ? (
        <Card>
          <Text style={styles.empty}>ยังไม่มีรายการวันนี้</Text>
        </Card>
      ) : (
        <Card style={styles.listCard}>
          {data.recent.map((s, i) => (
            <View key={i} style={[styles.row, i > 0 && styles.rowBorder]}>
              <EmojiAvatar emoji={s.emoji} size={32} />
              <View style={styles.rowMid}>
                <Text style={styles.rowName}>{s.name}</Text>
                <Text style={styles.rowTime}>{s.time} น.</Text>
              </View>
              <View style={styles.rowRight}>
                <Text style={styles.rowTotal}>{baht(s.total)}</Text>
                <Text
                  style={[
                    styles.rowMethod,
                    { color: s.method === 'cash' ? colors.cash : colors.transfer },
                  ]}>
                  {s.method === 'cash' ? 'เงินสด' : 'โอนเงิน'}
                </Text>
              </View>
            </View>
          ))}
        </Card>
      )}
    </ScrollView>
  );
};

// ===== Monthly =====
const MonthlyView: React.FC<{ state: ReturnType<typeof useApi<Awaited<ReturnType<typeof reportApi.monthlyReport>>>> }> = ({
  state,
}) => {
  const { data, loading, error, refetch } = state;
  if (loading && !data) {
    return <LoadingScreen />;
  }
  if (error || !data) {
    return <ErrorScreen message={error ?? undefined} onRetry={refetch} />;
  }

  const bars: Bar[] = data.week_bars.map((w) => ({
    label: w.label,
    value: w.amount,
    highlight: true,
    topLabel: baht(w.amount),
  }));

  return (
    <ScrollView contentContainerStyle={styles.body} showsVerticalScrollIndicator={false}>
      <Card>
        <Text style={styles.cardLabel}>ยอดขายเดือนนี้</Text>
        <Text style={styles.cardValue}>{baht(data.month_total)}</Text>
        <View style={styles.chart}>
          <BarChart bars={bars} height={130} />
        </View>
      </Card>

      <Card style={styles.statsCard}>
        <StatRow label="เฉลี่ยต่อวัน" value={baht(data.avg_per_day)} border />
        <StatRow label="จำนวนบิล" value={`${data.bill_count} บิล`} border />
        <StatRow label="ยอดเฉลี่ยต่อบิล" value={baht(data.avg_per_bill)} />
      </Card>
    </ScrollView>
  );
};

const StatRow: React.FC<{ label: string; value: string; border?: boolean }> = ({
  label,
  value,
  border,
}) => (
  <View style={[styles.statRow, border && styles.statBorder]}>
    <Text style={styles.statLabel}>{label}</Text>
    <Text style={styles.statValue}>{value}</Text>
  </View>
);

const styles = StyleSheet.create({
  header: {
    backgroundColor: colors.primary,
    paddingHorizontal: 20,
    paddingTop: 8,
    paddingBottom: 18,
    borderBottomLeftRadius: radius.xxl,
    borderBottomRightRadius: radius.xxl,
  },
  title: { fontSize: fontSize.xxl, color: '#fff', fontWeight: '600', marginBottom: 14 },
  body: { padding: 20 },
  cardLabel: { fontSize: fontSize.base, color: '#ab9683' },
  cardValue: { fontSize: fontSize.display, fontWeight: '700', color: colors.textDark, marginTop: 2 },
  chart: { marginTop: 18 },
  sectionTitle: {
    fontSize: fontSize.lg,
    fontWeight: '600',
    color: colors.textDark,
    marginTop: 22,
    marginBottom: 12,
  },
  listCard: { paddingVertical: 4 },
  empty: { textAlign: 'center', color: colors.textMuted, fontSize: fontSize.md, paddingVertical: 12 },
  row: { flexDirection: 'row', alignItems: 'center', gap: 12, paddingVertical: 11 },
  rowBorder: { borderTopWidth: 1, borderTopColor: colors.divider },
  rowMid: { flex: 1 },
  rowName: { fontSize: fontSize.md, color: colors.textDark },
  rowTime: { fontSize: fontSize.xs, color: '#ab9683' },
  rowRight: { alignItems: 'flex-end' },
  rowTotal: { fontSize: fontSize.md, fontWeight: '600', color: colors.textDark },
  rowMethod: { fontSize: fontSize.xs, fontWeight: '600' },
  statsCard: { marginTop: 14 },
  statRow: { flexDirection: 'row', justifyContent: 'space-between', paddingVertical: 8 },
  statBorder: { borderBottomWidth: 1, borderBottomColor: colors.divider },
  statLabel: { fontSize: fontSize.md, color: '#ab9683' },
  statValue: { fontSize: fontSize.md, fontWeight: '600', color: colors.textDark },
});
