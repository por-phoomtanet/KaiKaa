import React from 'react';
import { StyleSheet, View } from 'react-native';
import { colors } from '../theme';

export const Divider: React.FC = () => <View style={styles.line} />;

const styles = StyleSheet.create({
  line: { height: 1, backgroundColor: colors.divider, width: '100%' },
});
