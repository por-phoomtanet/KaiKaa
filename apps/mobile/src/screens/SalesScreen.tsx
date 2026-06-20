import React, { useCallback, useState } from 'react';
import { Pressable, ScrollView, StyleSheet, Text, View } from 'react-native';
import { errorMessage, productApi, saleApi } from '../api';
import { Product } from '../api/types';
import {
  EmojiAvatar,
  ErrorScreen,
  HeaderBar,
  LoadingScreen,
  PaymentSheet,
  ScreenWrapper,
  Toast,
} from '../components';
import { useApi } from '../hooks/useApi';
import { useToast } from '../hooks/useToast';
import { colors, fontSize, radius } from '../theme';
import { baht } from '../utils/format';

export const SalesScreen: React.FC = () => {
  const { data, loading, error, refetch } = useApi(() => productApi.listProducts());
  const [selected, setSelected] = useState<Product | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const toast = useToast();

  const pay = useCallback(
    async (method: 'cash' | 'transfer') => {
      if (!selected || submitting) {
        return; // กันกดรัวซ้ำ
      }
      setSubmitting(true);
      try {
        await saleApi.createSale({ product_id: selected.id, method });
        toast.show(
          (method === 'cash' ? 'รับเงินสด ' : 'รับโอนเงิน ') +
            selected.name +
            ' · ' +
            baht(selected.price),
        );
        setSelected(null);
      } catch (e) {
        toast.show(errorMessage(e));
      } finally {
        setSubmitting(false);
      }
    },
    [selected, submitting, toast],
  );

  if (loading && !data) {
    return <LoadingScreen />;
  }
  if (error || !data) {
    return <ErrorScreen message={error ?? undefined} onRetry={refetch} />;
  }

  return (
    <ScreenWrapper>
      <HeaderBar title="ขายสินค้า" subtitle="แตะเมนูเพื่อบันทึกการขายทันที" />

      {data.length === 0 ? (
        <View style={styles.empty}>
          <Text style={styles.emptyText}>ยังไม่มีสินค้า</Text>
          <Text style={styles.emptyHint}>เพิ่มสินค้าได้ที่แท็บ “สินค้า”</Text>
        </View>
      ) : (
        <ScrollView contentContainerStyle={styles.body} showsVerticalScrollIndicator={false}>
          <View style={styles.grid}>
            {data.map((p) => (
              <Pressable
                key={p.id}
                style={styles.card}
                onPress={() => setSelected(p)}>
                <EmojiAvatar emoji={p.emoji} size={46} />
                <Text style={styles.name}>{p.name}</Text>
                <Text style={styles.price}>{baht(p.price)}</Text>
              </Pressable>
            ))}
          </View>
        </ScrollView>
      )}

      <PaymentSheet
        product={selected}
        submitting={submitting}
        onPay={pay}
        onClose={() => setSelected(null)}
      />
      <Toast message={toast.message} visible={toast.visible} />
    </ScreenWrapper>
  );
};

const styles = StyleSheet.create({
  body: { padding: 20 },
  grid: { flexDirection: 'row', flexWrap: 'wrap', justifyContent: 'space-between', rowGap: 12 },
  card: {
    width: '48%',
    backgroundColor: colors.surface,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radius.xl,
    padding: 14,
    gap: 10,
  },
  name: { fontSize: fontSize.md, fontWeight: '500', color: colors.textDark, lineHeight: 18 },
  price: { fontSize: fontSize.lg, fontWeight: '700', color: colors.primary },
  empty: { flex: 1, alignItems: 'center', justifyContent: 'center', gap: 6 },
  emptyText: { fontSize: fontSize.lg, color: colors.textDark, fontWeight: '600' },
  emptyHint: { fontSize: fontSize.md, color: colors.textMuted },
});
