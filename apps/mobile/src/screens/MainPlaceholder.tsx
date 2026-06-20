import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { PrimaryButton, ScreenWrapper } from '../components';
import { useAuth } from '../context/AuthContext';
import { colors, fontSize } from '../theme';

// Placeholder ชั่วคราว — T09 จะแทนด้วย Bottom Tab Navigator (5 หน้า)
export const MainPlaceholder: React.FC = () => {
  const { shop, signOut } = useAuth();
  return (
    <ScreenWrapper>
      <View style={styles.container}>
        <Text style={styles.hi}>เข้าสู่ระบบแล้ว ✓</Text>
        <Text style={styles.shop}>{shop?.name ?? 'ร้านของคุณ'}</Text>
        <Text style={styles.note}>(หน้าหลัก + เมนูจะมาใน T09–T14)</Text>
        <PrimaryButton label="ออกจากระบบ" onPress={signOut} style={styles.btn} />
      </View>
    </ScreenWrapper>
  );
};

const styles = StyleSheet.create({
  container: { flex: 1, alignItems: 'center', justifyContent: 'center', padding: 28, gap: 8 },
  hi: { fontSize: fontSize.md, color: colors.textMuted },
  shop: { fontSize: fontSize.title, fontWeight: '700', color: colors.primary },
  note: { fontSize: fontSize.sm, color: colors.textMuted, marginBottom: 24 },
  btn: { alignSelf: 'stretch' },
});
