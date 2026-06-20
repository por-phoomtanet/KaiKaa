import React from 'react';
import { Pressable, StyleSheet, Text, View } from 'react-native';
import type { BottomTabBarProps } from '@react-navigation/bottom-tabs';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import Svg, { Path } from 'react-native-svg';
import { TabIcon } from './TabIcon';
import { colors } from '../theme';

const ACTIVE = colors.primary;
const INACTIVE = '#c3b3a2';

const labels: Record<string, string> = {
  Dashboard: 'หน้าหลัก',
  Reports: 'รายงาน',
  AI: 'AI',
  Hub: 'รวม',
};

// ปุ่มกลาง (ขาย) — วงกลมน้ำตาลยกสูง
const CenterButton: React.FC<{ onPress: () => void }> = ({ onPress }) => (
  <Pressable style={styles.centerWrap} onPress={onPress}>
    <View style={styles.centerCircle}>
      <Svg width={24} height={24} viewBox="0 0 24 24" fill="none" stroke="#fff" strokeWidth={2} strokeLinecap="round">
        <Path d="M12 6v12M6 12h12" />
      </Svg>
    </View>
    <Text style={styles.centerLabel}>ขาย</Text>
  </Pressable>
);

export const TabBar: React.FC<BottomTabBarProps> = ({ state, navigation }) => {
  const insets = useSafeAreaInsets();

  return (
    <View style={[styles.bar, { paddingBottom: Math.max(insets.bottom, 12) }]}>
      {state.routes.map((route, index) => {
        // Products เป็นหน้าในระบบนำทาง แต่ไม่แสดงปุ่มใน tab bar
        if (route.name === 'Products') {
          return null;
        }
        const focused = state.index === index;
        const onPress = () => {
          const event = navigation.emit({
            type: 'tabPress',
            target: route.key,
            canPreventDefault: true,
          });
          if (!focused && !event.defaultPrevented) {
            navigation.navigate(route.name);
          }
        };

        if (route.name === 'Sales') {
          return <CenterButton key={route.key} onPress={onPress} />;
        }

        const color = focused ? ACTIVE : INACTIVE;
        return (
          <Pressable key={route.key} style={styles.item} onPress={onPress}>
            <TabIcon name={route.name} color={color} />
            <Text style={[styles.label, { color }]}>{labels[route.name] ?? route.name}</Text>
          </Pressable>
        );
      })}
    </View>
  );
};

const styles = StyleSheet.create({
  bar: {
    flexDirection: 'row',
    alignItems: 'flex-end',
    backgroundColor: colors.surface,
    borderTopWidth: 1,
    borderTopColor: '#efe5d8',
    paddingTop: 8,
    paddingHorizontal: 6,
  },
  item: { flex: 1, alignItems: 'center', gap: 3 },
  label: { fontSize: 10.5, fontWeight: '600' },
  centerWrap: { flex: 1, alignItems: 'center' },
  centerCircle: {
    width: 50,
    height: 50,
    borderRadius: 25,
    backgroundColor: colors.primary,
    alignItems: 'center',
    justifyContent: 'center',
    transform: [{ translateY: -14 }],
    shadowColor: colors.primary,
    shadowOpacity: 0.4,
    shadowRadius: 9,
    shadowOffset: { width: 0, height: 8 },
    elevation: 6,
  },
  centerLabel: { fontSize: 10.5, fontWeight: '700', color: colors.primary, marginTop: -10 },
});
