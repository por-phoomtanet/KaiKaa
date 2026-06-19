import React from 'react';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import { View, Text, StyleSheet } from 'react-native';
import { theme } from './src/theme';

// Placeholder shell — navigation จริงเพิ่มใน T09
export default function App() {
  return (
    <SafeAreaProvider>
      <View style={styles.container}>
        <Text style={styles.title}>KaiKaa</Text>
        <Text style={styles.subtitle}>POS ร้านกาแฟ · scaffold พร้อมแล้ว</Text>
      </View>
    </SafeAreaProvider>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: theme.colors.background,
  },
  title: { fontSize: 34, fontWeight: '700', color: theme.colors.primary },
  subtitle: { fontSize: 14, color: theme.colors.textMuted, marginTop: 6 },
});
