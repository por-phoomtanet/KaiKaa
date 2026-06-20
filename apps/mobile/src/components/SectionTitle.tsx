import React from 'react';
import { StyleSheet, Text, ViewStyle } from 'react-native';
import { colors, fontSize } from '../theme';

interface Props {
  children: React.ReactNode;
  style?: ViewStyle;
}

export const SectionTitle: React.FC<Props> = ({ children }) => (
  <Text style={styles.title}>{children}</Text>
);

const styles = StyleSheet.create({
  title: { fontSize: fontSize.lg, fontWeight: '600', color: colors.textDark },
});
