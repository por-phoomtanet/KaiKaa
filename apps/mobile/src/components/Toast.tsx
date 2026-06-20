import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { fontSize, radius } from '../theme';

interface Props {
  message: string;
  visible: boolean;
}

// Toast แบบ presentational — parent คุม visible (auto-dismiss ทำที่ระดับ screen/hook)
export const Toast: React.FC<Props> = ({ message, visible }) => {
  if (!visible) {
    return null;
  }
  return (
    <View style={styles.toast} pointerEvents="none">
      <Text style={styles.check}>✓</Text>
      <Text style={styles.msg}>{message}</Text>
    </View>
  );
};

const styles = StyleSheet.create({
  toast: {
    position: 'absolute',
    bottom: 96,
    alignSelf: 'center',
    flexDirection: 'row',
    alignItems: 'center',
    gap: 8,
    backgroundColor: '#4a3526',
    paddingVertical: 11,
    paddingHorizontal: 18,
    borderRadius: radius.xxl,
  },
  check: { color: '#7fe3a8', fontSize: fontSize.md },
  msg: { color: '#fff', fontSize: fontSize.base, fontWeight: '500' },
});
