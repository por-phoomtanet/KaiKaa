import React from 'react';
import { StyleSheet, View } from 'react-native';
import { colors } from '../theme';

interface Props {
  pct: number; // 0-100
  color?: string;
  height?: number;
}

export const ProgressBar: React.FC<Props> = ({ pct, color = colors.primary, height = 5 }) => {
  const clamped = Math.max(0, Math.min(100, pct));
  return (
    <View style={[styles.track, { height, borderRadius: height }]}>
      <View
        style={{ width: `${clamped}%`, height: '100%', backgroundColor: color, borderRadius: height }}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  track: { backgroundColor: colors.background, overflow: 'hidden', width: '100%' },
});
