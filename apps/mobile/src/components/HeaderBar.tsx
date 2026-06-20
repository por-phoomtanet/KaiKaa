import React from 'react';
import { Pressable, StyleSheet, Text, View } from 'react-native';
import { fontSize, radius } from '../theme';

interface Props {
  title: string;
  subtitle?: string;
  onBack?: () => void;
}

export const HeaderBar: React.FC<Props> = ({ title, subtitle, onBack }) => (
  <View style={styles.header}>
    {onBack && (
      <Pressable onPress={onBack} hitSlop={10}>
        <Text style={styles.back}>‹ เมนูทั้งหมด</Text>
      </Pressable>
    )}
    <Text style={styles.title}>{title}</Text>
    {subtitle ? <Text style={styles.subtitle}>{subtitle}</Text> : null}
  </View>
);

const styles = StyleSheet.create({
  header: {
    backgroundColor: '#b06a43',
    paddingHorizontal: 20,
    paddingTop: 8,
    paddingBottom: 22,
    borderBottomLeftRadius: radius.xxl,
    borderBottomRightRadius: radius.xxl,
  },
  back: { color: 'rgba(255,255,255,0.85)', fontSize: fontSize.base, marginBottom: 6 },
  title: { color: '#fff', fontSize: fontSize.xxl, fontWeight: '600' },
  subtitle: { color: 'rgba(255,255,255,0.72)', fontSize: fontSize.base, marginTop: 2 },
});
