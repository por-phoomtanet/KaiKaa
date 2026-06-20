import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { radius } from '../theme';

interface Props {
  emoji: string;
  size?: number;
}

export const EmojiAvatar: React.FC<Props> = ({ emoji, size = 40 }) => (
  <View style={[styles.box, { width: size, height: size, borderRadius: radius.md }]}>
    <Text style={{ fontSize: size * 0.5 }}>{emoji}</Text>
  </View>
);

const styles = StyleSheet.create({
  box: {
    backgroundColor: '#f6ede1', // เฉดครีมเข้มสำหรับพื้น emoji
    alignItems: 'center',
    justifyContent: 'center',
  },
});
