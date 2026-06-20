import React from 'react';
import { Pressable, StyleSheet, Text, View } from 'react-native';
import { Product } from '../api/types';
import { colors, fontSize, radius } from '../theme';
import { baht } from '../utils/format';
import { BottomSheet } from './BottomSheet';
import { EmojiAvatar } from './EmojiAvatar';

interface Props {
  product: Product | null;
  submitting?: boolean;
  onPay: (method: 'cash' | 'transfer') => void;
  onClose: () => void;
}

export const PaymentSheet: React.FC<Props> = ({ product, submitting, onPay, onClose }) => (
  <BottomSheet visible={product !== null} onClose={onClose}>
    {product && (
      <View>
        <View style={styles.head}>
          <EmojiAvatar emoji={product.emoji} size={54} />
          <View style={styles.headMid}>
            <Text style={styles.name}>{product.name}</Text>
            <Text style={styles.sub}>เลือกวิธีชำระเงิน</Text>
          </View>
          <Text style={styles.price}>{baht(product.price)}</Text>
        </View>

        <View style={styles.btnRow}>
          <Pressable
            style={[styles.payBtn, { backgroundColor: '#e8f5ee' }]}
            disabled={submitting}
            onPress={() => onPay('cash')}>
            <Text style={styles.payEmoji}>💵</Text>
            <Text style={[styles.payLabel, { color: colors.cash }]}>เงินสด</Text>
          </Pressable>
          <Pressable
            style={[styles.payBtn, { backgroundColor: '#eaf0fd' }]}
            disabled={submitting}
            onPress={() => onPay('transfer')}>
            <Text style={styles.payEmoji}>📲</Text>
            <Text style={[styles.payLabel, { color: colors.transfer }]}>โอนเงิน</Text>
          </Pressable>
        </View>

        <Pressable onPress={onClose} style={styles.cancel}>
          <Text style={styles.cancelText}>ยกเลิก</Text>
        </Pressable>
      </View>
    )}
  </BottomSheet>
);

const styles = StyleSheet.create({
  head: { flexDirection: 'row', alignItems: 'center', gap: 14 },
  headMid: { flex: 1 },
  name: { fontSize: fontSize.xl, fontWeight: '600', color: colors.textDark },
  sub: { fontSize: fontSize.base, color: colors.textMuted },
  price: { fontSize: fontSize.display, fontWeight: '700', color: colors.primary },
  btnRow: { flexDirection: 'row', gap: 12, marginTop: 20 },
  payBtn: {
    flex: 1,
    borderRadius: radius.xl,
    paddingVertical: 16,
    alignItems: 'center',
    gap: 5,
  },
  payEmoji: { fontSize: 24 },
  payLabel: { fontSize: fontSize.lg, fontWeight: '700' },
  cancel: { marginTop: 12, padding: 8, alignItems: 'center' },
  cancelText: { color: colors.textMuted, fontSize: fontSize.md },
});
