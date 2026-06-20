import React from 'react';
import { Pressable, StyleSheet, Text, ViewStyle } from 'react-native';
import { colors, fontSize, radius } from '../theme';

interface Props {
  label: string;
  onPress: () => void;
  style?: ViewStyle;
}

export const GhostButton: React.FC<Props> = ({ label, onPress, style }) => (
  <Pressable
    onPress={onPress}
    style={({ pressed }) => [styles.btn, style, pressed && styles.pressed]}>
    <Text style={styles.label}>{label}</Text>
  </Pressable>
);

const styles = StyleSheet.create({
  btn: {
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radius.lg,
    paddingVertical: 14,
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: colors.surface,
  },
  pressed: { opacity: 0.6 },
  label: { color: colors.textMuted, fontSize: fontSize.md, fontWeight: '600' },
});
