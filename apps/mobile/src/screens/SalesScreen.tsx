import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { ScreenWrapper } from '../components';
import { colors, fontSize } from '../theme';

// Placeholder — T11 จะใส่ grid สินค้า + payment sheet
export const SalesScreen: React.FC = () => (
  <ScreenWrapper>
    <View style={styles.c}>
      <Text style={styles.t}>ขายสินค้า</Text>
      <Text style={styles.s}>POS (T11)</Text>
    </View>
  </ScreenWrapper>
);

const styles = StyleSheet.create({
  c: { flex: 1, alignItems: 'center', justifyContent: 'center', gap: 4 },
  t: { fontSize: fontSize.title, fontWeight: '700', color: colors.primary },
  s: { fontSize: fontSize.md, color: colors.textMuted },
});
