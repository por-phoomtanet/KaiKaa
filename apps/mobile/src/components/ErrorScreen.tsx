import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { colors, fontSize } from '../theme';
import { PrimaryButton } from './PrimaryButton';

interface Props {
  message?: string;
  onRetry?: () => void;
}

export const ErrorScreen: React.FC<Props> = ({ message, onRetry }) => (
  <View style={styles.container}>
    <Text style={styles.emoji}>😕</Text>
    <Text style={styles.message}>{message ?? 'เกิดข้อผิดพลาด'}</Text>
    {onRetry && <PrimaryButton label="ลองอีกครั้ง" onPress={onRetry} />}
  </View>
);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: colors.background,
    padding: 24,
    gap: 14,
  },
  emoji: { fontSize: 40 },
  message: { fontSize: fontSize.md, color: colors.textMuted, textAlign: 'center' },
});
