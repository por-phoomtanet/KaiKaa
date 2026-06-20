import React from 'react';
import { Pressable, ScrollView, StyleSheet, Text, View } from 'react-native';
import { useNavigation } from '@react-navigation/native';
import type { BottomTabNavigationProp } from '@react-navigation/bottom-tabs';
import { HeaderBar, ScreenWrapper } from '../components';
import { useAuth } from '../context/AuthContext';
import { MainTabParamList } from '../navigation/types';
import { colors, fontSize, radius } from '../theme';

type Nav = BottomTabNavigationProp<MainTabParamList>;

interface Shortcut {
  label: string;
  emoji: string;
  bg: string;
  target: keyof MainTabParamList;
}

const shortcuts: Shortcut[] = [
  { label: 'หน้าหลัก', emoji: '🏠', bg: '#f3e2d3', target: 'Dashboard' },
  { label: 'ขายสินค้า', emoji: '🛒', bg: '#e8f5ee', target: 'Sales' },
  { label: 'สินค้า', emoji: '🧋', bg: '#eaf0fd', target: 'Products' },
  { label: 'รายงาน', emoji: '📊', bg: '#fbe8df', target: 'Reports' },
  { label: 'AI สรุป', emoji: '✨', bg: '#f3e7d0', target: 'AI' },
];

// เมนูรวม — shortcut ไปแต่ละหน้า + ปุ่มออกจากระบบ
export const HubScreen: React.FC = () => {
  const navigation = useNavigation<Nav>();
  const { shop, signOut } = useAuth();

  return (
    <ScreenWrapper>
      <HeaderBar title="เมนูทั้งหมด" subtitle={shop?.name ?? 'ร้านของคุณ'} />
      <ScrollView contentContainerStyle={styles.body}>
        <View style={styles.grid}>
          {shortcuts.map((s) => (
            <Pressable key={s.target} style={styles.item} onPress={() => navigation.navigate(s.target)}>
              <View style={[styles.icon, { backgroundColor: s.bg }]}>
                <Text style={styles.emoji}>{s.emoji}</Text>
              </View>
              <Text style={styles.label}>{s.label}</Text>
            </Pressable>
          ))}
        </View>

        <Pressable style={styles.logout} onPress={signOut}>
          <Text style={styles.logoutText}>ออกจากระบบ</Text>
        </Pressable>
      </ScrollView>
    </ScreenWrapper>
  );
};

const styles = StyleSheet.create({
  body: { padding: 22 },
  grid: { flexDirection: 'row', flexWrap: 'wrap', gap: 12 },
  item: { width: '30%', alignItems: 'center', gap: 9, marginBottom: 14 },
  icon: { width: 62, height: 62, borderRadius: radius.xl, alignItems: 'center', justifyContent: 'center' },
  emoji: { fontSize: 28 },
  label: { fontSize: fontSize.sm, color: colors.textDark, fontWeight: '500' },
  logout: { marginTop: 24, alignItems: 'center', padding: 12 },
  logoutText: { color: '#d2604f', fontSize: fontSize.md, fontWeight: '600' },
});
