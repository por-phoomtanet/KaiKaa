import React from 'react';
import { DimensionValue, StyleSheet, Text, View } from 'react-native';
import { colors, fontSize, radius } from '../theme';

export interface Bar {
  label: string;
  value: number;
  highlight?: boolean;
  topLabel?: string; // ข้อความเหนือแท่ง (เช่น ยอดเงิน)
}

interface Props {
  bars: Bar[];
  height?: number;
}

// Bar chart แบบ View ล้วน (ไม่พึ่ง library) — สูงสัมพันธ์กับค่าสูงสุด
export const BarChart: React.FC<Props> = ({ bars, height = 120 }) => {
  const max = Math.max(1, ...bars.map((b) => b.value));
  return (
    <View style={[styles.row, { height }]}>
      {bars.map((b, i) => {
        const h: DimensionValue = `${Math.round((b.value / max) * 100)}%`;
        return (
          <View key={i} style={styles.col}>
            {b.topLabel ? <Text style={styles.topLabel}>{b.topLabel}</Text> : null}
            <View
              style={[
                styles.bar,
                { height: h, backgroundColor: b.highlight ? colors.primary : '#ecdcc8' },
              ]}
            />
            <Text style={styles.label}>{b.label}</Text>
          </View>
        );
      })}
    </View>
  );
};

const styles = StyleSheet.create({
  row: { flexDirection: 'row', alignItems: 'flex-end', gap: 9 },
  col: { flex: 1, alignItems: 'center', justifyContent: 'flex-end', height: '100%', gap: 7 },
  bar: { width: '100%', borderTopLeftRadius: radius.sm, borderTopRightRadius: radius.sm },
  label: { fontSize: fontSize.xs, color: colors.textMuted },
  topLabel: { fontSize: fontSize.xs, fontWeight: '600', color: colors.textDark },
});
