import React from 'react';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { MainTabParamList } from './types';
import { TabBar } from './TabBar';
import { DashboardScreen } from '../screens/DashboardScreen';
import { ReportsScreen } from '../screens/ReportsScreen';
import { SalesScreen } from '../screens/SalesScreen';
import { AiScreen } from '../screens/AiScreen';
import { HubScreen } from '../screens/HubScreen';
import { ProductsScreen } from '../screens/ProductsScreen';

const Tab = createBottomTabNavigator<MainTabParamList>();

// แท็บล่าง 5 หน้า — ใช้ custom TabBar (ปุ่มขายกลางยกสูง)
export const MainTabs: React.FC = () => (
  <Tab.Navigator
    screenOptions={{ headerShown: false }}
    tabBar={(props) => <TabBar {...props} />}>
    <Tab.Screen name="Dashboard" component={DashboardScreen} />
    <Tab.Screen name="Reports" component={ReportsScreen} />
    <Tab.Screen name="Sales" component={SalesScreen} />
    <Tab.Screen name="AI" component={AiScreen} />
    <Tab.Screen name="Hub" component={HubScreen} />
    {/* ซ่อนจาก tab bar (กรองใน TabBar) — เข้าจาก Hub */}
    <Tab.Screen name="Products" component={ProductsScreen} />
  </Tab.Navigator>
);
