import React from 'react';
import { LoadingScreen } from '../components';
import { useAuth } from '../context/AuthContext';
import { AuthScreen } from '../screens/auth/AuthScreen';
import { MainPlaceholder } from '../screens/MainPlaceholder';

// Auth guard: ยังโหลด session → loading · ไม่มี token → Login · มี token → แอปหลัก
export const RootNavigator: React.FC = () => {
  const { token, initializing } = useAuth();

  if (initializing) {
    return <LoadingScreen />;
  }
  return token ? <MainPlaceholder /> : <AuthScreen />;
};
