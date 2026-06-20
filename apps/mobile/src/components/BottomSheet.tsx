import React from 'react';
import { Modal, Pressable, StyleSheet, View } from 'react-native';
import { colors, radius } from '../theme';

interface Props {
  visible: boolean;
  onClose: () => void;
  children: React.ReactNode;
}

// Bottom sheet เบาๆ ด้วย RN Modal (slide-up ในตัว) + drag handle + dim backdrop
export const BottomSheet: React.FC<Props> = ({ visible, onClose, children }) => (
  <Modal visible={visible} transparent animationType="slide" onRequestClose={onClose}>
    <View style={styles.fill}>
      <Pressable style={styles.backdrop} onPress={onClose} />
      <View style={styles.sheet}>
        <View style={styles.handle} />
        {children}
      </View>
    </View>
  </Modal>
);

const styles = StyleSheet.create({
  fill: { flex: 1, justifyContent: 'flex-end' },
  backdrop: { ...StyleSheet.absoluteFillObject, backgroundColor: 'rgba(46,29,18,0.45)' },
  sheet: {
    backgroundColor: colors.surface,
    borderTopLeftRadius: radius.xxl,
    borderTopRightRadius: radius.xxl,
    paddingHorizontal: 22,
    paddingTop: 14,
    paddingBottom: 28,
  },
  handle: {
    width: 38,
    height: 4,
    borderRadius: 3,
    backgroundColor: '#e4d6c5',
    alignSelf: 'center',
    marginBottom: 16,
  },
});
