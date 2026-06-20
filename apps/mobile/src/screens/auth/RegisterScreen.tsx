import React, { useState } from 'react';
import { Pressable, ScrollView, StyleSheet, Text, View } from 'react-native';
import { authApi, errorMessage } from '../../api';
import { PrimaryButton, ScreenWrapper, TextField } from '../../components';
import { useAuth } from '../../context/AuthContext';
import { colors, fontSize } from '../../theme';

interface Props {
  onSwitch: () => void;
}

export const RegisterScreen: React.FC<Props> = ({ onSwitch }) => {
  const { signIn } = useAuth();
  const [shopName, setShopName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const submit = async () => {
    if (submitting) {
      return;
    }
    setSubmitting(true);
    setError(null);
    try {
      const data = await authApi.register({
        shop_name: shopName.trim(),
        email: email.trim(),
        password,
      });
      signIn(data.token, data.shop);
    } catch (e) {
      setError(errorMessage(e));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <ScreenWrapper>
      <ScrollView contentContainerStyle={styles.container} keyboardShouldPersistTaps="handled">
        <Text style={styles.brand}>สมัครร้านใหม่</Text>
        <Text style={styles.subtitle}>เริ่มใช้งาน KaiKaa ฟรี</Text>

        <View style={styles.form}>
          <TextField
            label="ชื่อร้าน"
            value={shopName}
            onChangeText={setShopName}
            placeholder="เช่น ร้านกาแฟ Brewtiful"
            autoCapitalize="sentences"
          />
          <TextField
            label="อีเมล"
            value={email}
            onChangeText={setEmail}
            placeholder="owner@shop.com"
            keyboardType="email-address"
          />
          <TextField
            label="รหัสผ่าน (อย่างน้อย 6 ตัว)"
            value={password}
            onChangeText={setPassword}
            placeholder="••••••••"
            secureTextEntry
          />

          {error ? <Text style={styles.error}>{error}</Text> : null}

          <PrimaryButton
            label={submitting ? 'กำลังสมัคร...' : 'สมัครและเริ่มใช้งาน'}
            onPress={submit}
            disabled={submitting}
          />
        </View>

        <Pressable onPress={onSwitch} style={styles.switch}>
          <Text style={styles.switchText}>
            มีบัญชีอยู่แล้ว? <Text style={styles.switchLink}>เข้าสู่ระบบ</Text>
          </Text>
        </Pressable>
      </ScrollView>
    </ScreenWrapper>
  );
};

const styles = StyleSheet.create({
  container: { flexGrow: 1, justifyContent: 'center', padding: 28 },
  brand: { fontSize: fontSize.title, fontWeight: '700', color: colors.primary, textAlign: 'center' },
  subtitle: {
    fontSize: fontSize.md,
    color: colors.textMuted,
    textAlign: 'center',
    marginTop: 4,
    marginBottom: 28,
  },
  form: { gap: 14 },
  error: { color: '#d2604f', fontSize: fontSize.base },
  switch: { marginTop: 22, alignItems: 'center' },
  switchText: { fontSize: fontSize.md, color: colors.textMuted },
  switchLink: { color: colors.primary, fontWeight: '600' },
});
