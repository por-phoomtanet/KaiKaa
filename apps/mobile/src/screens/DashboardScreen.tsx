import React, { useCallback, useRef } from 'react';
import { ScrollView, StyleSheet, Text, View } from 'react-native';
import { useFocusEffect } from '@react-navigation/native';
import { reportApi } from '../api';
import {
  Card,
  EmojiAvatar,
  ErrorScreen,
  LoadingScreen,
  ProgressBar,
  ScreenWrapper,
  SectionTitle,
} from '../components';
import { useApi } from '../hooks/useApi';
import { useAuth } from '../context/AuthContext';
import { colors, fontSize, radius } from '../theme';
import { baht, greeting } from '../utils/format';

export const DashboardScreen: React.FC = () => {
  const { shop } = useAuth();
  const { data, loading, error, refetch } = useApi(() => reportApi.dailyReport());

  // refetch เมื่อกลับมาโฟกัสหน้านี้ (เช่น หลังขายสินค้า) — ข้ามครั้งแรกที่ mount
  const firstFocus = useRef(true);
  useFocusEffect(
    useCallback(() => {
      if (firstFocus.current) {
        firstFocus.current = false;
        return;
      }
      refetch();
    }, [refetch]),
  );

  if (loading && !data) {
    return <LoadingScreen />;
  }
  if (error || !data) {
    return <ErrorScreen message={error ?? undefined} onRetry={refetch} />;
  }

  const up = data.change_pct >= 0;

  return (
    <ScreenWrapper>
      <ScrollView showsVerticalScrollIndicator={false}>
        {/* header */}
        <View style={styles.header}>
          <View style={styles.headerRow}>
            <View>
              <Text style={styles.greeting}>{greeting()}</Text>
              <Text style={styles.shop}>{shop?.name ?? 'ร้านของคุณ'}</Text>
            </View>
            <View style={styles.avatar}>
              <Text style={{ fontSize: 18 }}>🧑‍🍳</Text>
            </View>
          </View>

          {/* total card */}
          <View style={styles.totalCard}>
            <Text style={styles.totalLabel}>ยอดขายวันนี้</Text>
            <Text style={styles.totalValue}>{baht(data.total)}</Text>
            <View style={styles.badge}>
              <Text style={[styles.badgeChange, { color: up ? '#bfe3cf' : '#f3b0a4' }]}>
                {up ? '▲' : '▼'} {Math.abs(data.change_pct)}%
              </Text>
              <Text style={styles.badgeCount}>{data.count} รายการ</Text>
            </View>
          </View>
        </View>

        {/* cash / transfer */}
        <View style={styles.body}>
          <View style={styles.statRow}>
            <Card style={styles.statCard}>
              <Text style={styles.statLabel}>เงินสด</Text>
              <Text style={styles.statValue}>{baht(data.cash)}</Text>
              <ProgressBar pct={data.cash_pct} color={colors.cash} />
              <Text style={styles.statHint}>{data.cash_pct}% ของยอดรวม</Text>
            </Card>
            <Card style={styles.statCard}>
              <Text style={styles.statLabel}>โอนเงิน</Text>
              <Text style={styles.statValue}>{baht(data.transfer)}</Text>
              <ProgressBar pct={data.transfer_pct} color={colors.transfer} />
              <Text style={styles.statHint}>{data.transfer_pct}% ของยอดรวม</Text>
            </Card>
          </View>

          {/* best sellers */}
          <View style={styles.sectionHead}>
            <SectionTitle>สินค้าขายดี</SectionTitle>
            <Text style={styles.sectionTag}>วันนี้</Text>
          </View>

          {data.best_sellers.length === 0 ? (
            <Card>
              <Text style={styles.empty}>ยังไม่มีการขายวันนี้</Text>
            </Card>
          ) : (
            <Card style={styles.listCard}>
              {data.best_sellers.map((b, i) => (
                <View key={b.rank} style={[styles.row, i > 0 && styles.rowBorder]}>
                  <EmojiAvatar emoji={b.emoji} size={34} />
                  <View style={styles.rowMid}>
                    <Text style={styles.rowName}>{b.name}</Text>
                    <ProgressBar pct={b.pct_width} height={4} />
                  </View>
                  <View style={styles.rowRight}>
                    <Text style={styles.rowRev}>{baht(b.revenue)}</Text>
                    <Text style={styles.rowQty}>{b.qty} แก้ว</Text>
                  </View>
                </View>
              ))}
            </Card>
          )}
        </View>
      </ScrollView>
    </ScreenWrapper>
  );
};

const styles = StyleSheet.create({
  header: {
    backgroundColor: colors.primary,
    paddingHorizontal: 20,
    paddingTop: 6,
    paddingBottom: 26,
    borderBottomLeftRadius: radius.xxl,
    borderBottomRightRadius: radius.xxl,
  },
  headerRow: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center' },
  greeting: { fontSize: fontSize.base, color: 'rgba(255,255,255,0.7)' },
  shop: { fontSize: fontSize.xxl, color: '#fff', fontWeight: '600', marginTop: 2 },
  avatar: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: 'rgba(255,255,255,0.14)',
    alignItems: 'center',
    justifyContent: 'center',
  },
  totalCard: {
    marginTop: 20,
    backgroundColor: '#c8895c',
    borderWidth: 1,
    borderColor: 'rgba(255,255,255,0.12)',
    borderRadius: radius.xxl,
    padding: 18,
  },
  totalLabel: { fontSize: fontSize.base, color: 'rgba(255,255,255,0.72)' },
  totalValue: { fontSize: fontSize.hero, color: '#fff', fontWeight: '700', marginTop: 2 },
  badge: {
    marginTop: 6,
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    alignSelf: 'flex-start',
    backgroundColor: 'rgba(255,255,255,0.12)',
    paddingHorizontal: 10,
    paddingVertical: 3,
    borderRadius: 20,
  },
  badgeChange: { fontSize: fontSize.sm },
  badgeCount: { fontSize: fontSize.sm, color: 'rgba(255,255,255,0.7)' },
  body: { padding: 20 },
  statRow: { flexDirection: 'row', gap: 12 },
  statCard: { flex: 1, borderRadius: radius.xl },
  statLabel: { fontSize: fontSize.base, color: '#ab9683' },
  statValue: { fontSize: fontSize.title, fontWeight: '700', color: '#4a3526', marginVertical: 8 },
  statHint: { fontSize: fontSize.xs, color: '#ab9683', marginTop: 5 },
  sectionHead: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginTop: 24,
    marginBottom: 12,
  },
  sectionTag: { fontSize: fontSize.sm, color: colors.primary, fontWeight: '600' },
  listCard: { paddingVertical: 6 },
  empty: { textAlign: 'center', color: colors.textMuted, fontSize: fontSize.md, paddingVertical: 12 },
  row: { flexDirection: 'row', alignItems: 'center', gap: 12, paddingVertical: 11 },
  rowBorder: { borderTopWidth: 1, borderTopColor: colors.divider },
  rowMid: { flex: 1, gap: 6 },
  rowName: { fontSize: fontSize.md, fontWeight: '500', color: colors.textDark },
  rowRight: { alignItems: 'flex-end' },
  rowRev: { fontSize: fontSize.base, fontWeight: '600', color: colors.textDark },
  rowQty: { fontSize: fontSize.xs, color: '#ab9683' },
});
