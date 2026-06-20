import React from 'react';
import { Pressable, StyleSheet, Text, View } from 'react-native';
import { fontSize, radius } from '../theme';

interface Tab {
  key: string;
  label: string;
}

interface Props {
  tabs: Tab[];
  active: string;
  onChange: (key: string) => void;
}

// สลับแท็บบนพื้น header (เช่น รายวัน/รายเดือน)
export const TabSwitch: React.FC<Props> = ({ tabs, active, onChange }) => (
  <View style={styles.container}>
    {tabs.map((t) => {
      const isActive = t.key === active;
      return (
        <Pressable
          key={t.key}
          onPress={() => onChange(t.key)}
          style={[styles.tab, isActive && styles.tabActive]}>
          <Text style={[styles.label, isActive && styles.labelActive]}>{t.label}</Text>
        </Pressable>
      );
    })}
  </View>
);

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    gap: 6,
    backgroundColor: 'rgba(255,255,255,0.12)',
    borderRadius: radius.md,
    padding: 4,
  },
  tab: { flex: 1, borderRadius: radius.sm, paddingVertical: 9, alignItems: 'center' },
  tabActive: { backgroundColor: '#b06a43' },
  label: { fontSize: fontSize.base, fontWeight: '600', color: 'rgba(255,255,255,0.65)' },
  labelActive: { color: '#fff' },
});
