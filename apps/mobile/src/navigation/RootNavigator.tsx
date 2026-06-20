import React from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { LoadingScreen } from '../components';
import { useAuth } from '../context/AuthContext';
import { AuthScreen } from '../screens/auth/AuthScreen';
import { MainTabs } from './MainTabs';

// Auth guard: ยังโหลด session → loading · ไม่มี token → Login · มี token → แท็บหลัก
export const RootNavigator: React.FC = () => {
  const { token, initializing } = useAuth();

  if (initializing) {
    return <LoadingScreen />;
  }

  if (!token) {
    return <AuthScreen />;
  }

  return (
    <NavigationContainer>
      <MainTabs />
    </NavigationContainer>
  );
};
