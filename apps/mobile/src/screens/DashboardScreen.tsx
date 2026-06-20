import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { ScreenWrapper } from '../components';
import { colors, fontSize } from '../theme';

// Placeholder — T10 จะใส่ยอดขายวันนี้ + cash/transfer + สินค้าขายดี
export const DashboardScreen: React.FC = () => (
  <ScreenWrapper>
    <View style={styles.c}>
      <Text style={styles.t}>หน้าหลัก</Text>
      <Text style={styles.s}>Dashboard (T10)</Text>
    </View>
  </ScreenWrapper>
);

const styles = StyleSheet.create({
  c: { flex: 1, alignItems: 'center', justifyContent: 'center', gap: 4 },
  t: { fontSize: fontSize.title, fontWeight: '700', color: colors.primary },
  s: { fontSize: fontSize.md, color: colors.textMuted },
});
