import React, { useEffect, useState } from 'react';
import { StyleSheet, Text } from 'react-native';
import { Product, ProductInput } from '../api/types';
import { colors, fontSize } from '../theme';
import { BottomSheet } from './BottomSheet';
import { PrimaryButton } from './PrimaryButton';
import { TextField } from './TextField';

interface Props {
  visible: boolean;
  initial?: Product | null; // null/undefined = โหมดเพิ่มใหม่
  submitting?: boolean;
  onSave: (input: ProductInput) => void;
  onClose: () => void;
}

export const ProductEditorSheet: React.FC<Props> = ({
  visible,
  initial,
  submitting,
  onSave,
  onClose,
}) => {
  const [name, setName] = useState('');
  const [price, setPrice] = useState('');
  const [emoji, setEmoji] = useState('☕');

  // reset ฟอร์มทุกครั้งที่เปิด sheet (ตาม initial)
  useEffect(() => {
    if (visible) {
      setName(initial?.name ?? '');
      setPrice(initial ? String(initial.price) : '');
      setEmoji(initial?.emoji ?? '☕');
    }
  }, [visible, initial]);

  const priceNum = Number(price);
  const valid = name.trim().length > 0 && price.trim() !== '' && !isNaN(priceNum) && priceNum >= 0;

  const save = () => {
    if (!valid || submitting) {
      return;
    }
    onSave({ name: name.trim(), price: Math.round(priceNum), emoji: emoji.trim() || '☕' });
  };

  return (
    <BottomSheet visible={visible} onClose={onClose}>
      <Text style={styles.title}>{initial ? 'แก้ไขสินค้า' : 'เพิ่มสินค้าใหม่'}</Text>

      <TextField
        label="ชื่อสินค้า"
        value={name}
        onChangeText={setName}
        placeholder="เช่น ลาเต้เย็น"
        autoCapitalize="sentences"
        style={styles.field}
      />
      <TextField
        label="ราคา (บาท)"
        value={price}
        onChangeText={setPrice}
        placeholder="0"
        keyboardType="number-pad"
        style={styles.field}
      />
      <TextField
        label="อิโมจิ"
        value={emoji}
        onChangeText={setEmoji}
        placeholder="☕"
        style={styles.field}
      />

      <PrimaryButton
        label={submitting ? 'กำลังบันทึก...' : 'บันทึก'}
        onPress={save}
        disabled={!valid || submitting}
        style={styles.save}
      />
    </BottomSheet>
  );
};

const styles = StyleSheet.create({
  title: { fontSize: fontSize.xl, fontWeight: '600', color: colors.textDark, marginBottom: 16 },
  field: { marginBottom: 14 },
  save: { marginTop: 6 },
});
