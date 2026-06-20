import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { ScreenWrapper } from '../components';
import { colors, fontSize } from '../theme';

// Placeholder — T13 จะใส่ tab รายวัน/รายเดือน + bar chart
export const ReportsScreen: React.FC = () => (
  <ScreenWrapper>
    <View style={styles.c}>
      <Text style={styles.t}>รายงาน</Text>
      <Text style={styles.s}>Reports (T13)</Text>
    </View>
  </ScreenWrapper>
);

const styles = StyleSheet.create({
  c: { flex: 1, alignItems: 'center', justifyContent: 'center', gap: 4 },
  t: { fontSize: fontSize.title, fontWeight: '700', color: colors.primary },
  s: { fontSize: fontSize.md, color: colors.textMuted },
});
