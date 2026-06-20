import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { ScreenWrapper } from '../components';
import { colors, fontSize } from '../theme';

// Placeholder — T14 จะใส่ปุ่มสรุป + 3 การ์ดผล AI
export const AiScreen: React.FC = () => (
  <ScreenWrapper>
    <View style={styles.c}>
      <Text style={styles.t}>AI สรุป</Text>
      <Text style={styles.s}>AI Summary (T14)</Text>
    </View>
  </ScreenWrapper>
);

const styles = StyleSheet.create({
  c: { flex: 1, alignItems: 'center', justifyContent: 'center', gap: 4 },
  t: { fontSize: fontSize.title, fontWeight: '700', color: colors.primary },
  s: { fontSize: fontSize.md, color: colors.textMuted },
});
