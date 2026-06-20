import React, { useCallback, useState } from 'react';
import { Alert, Pressable, ScrollView, StyleSheet, Text, View } from 'react-native';
import { useNavigation } from '@react-navigation/native';
import type { BottomTabNavigationProp } from '@react-navigation/bottom-tabs';
import { errorMessage, productApi } from '../api';
import { Product, ProductInput } from '../api/types';
import {
  EmojiAvatar,
  ErrorScreen,
  HeaderBar,
  LoadingScreen,
  ProductEditorSheet,
  ScreenWrapper,
  Toast,
} from '../components';
import { useApi } from '../hooks/useApi';
import { useToast } from '../hooks/useToast';
import { MainTabParamList } from '../navigation/types';
import { colors, fontSize, radius } from '../theme';
import { baht } from '../utils/format';

type Nav = BottomTabNavigationProp<MainTabParamList>;

export const ProductsScreen: React.FC = () => {
  const navigation = useNavigation<Nav>();
  const { data, loading, error, refetch } = useApi(() => productApi.listProducts());
  const [editorOpen, setEditorOpen] = useState(false);
  const [editing, setEditing] = useState<Product | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const toast = useToast();

  const openAdd = () => {
    setEditing(null);
    setEditorOpen(true);
  };
  const openEdit = (p: Product) => {
    setEditing(p);
    setEditorOpen(true);
  };

  const save = useCallback(
    async (input: ProductInput) => {
      if (submitting) {
        return;
      }
      setSubmitting(true);
      try {
        if (editing) {
          await productApi.updateProduct(editing.id, input);
          toast.show('บันทึกการแก้ไขแล้ว');
        } else {
          await productApi.createProduct(input);
          toast.show('เพิ่มสินค้าใหม่แล้ว');
        }
        setEditorOpen(false);
        setEditing(null);
        refetch();
      } catch (e) {
        toast.show(errorMessage(e));
      } finally {
        setSubmitting(false);
      }
    },
    [editing, submitting, toast, refetch],
  );

  const confirmDelete = (p: Product) => {
    Alert.alert('ลบสินค้า', `ต้องการลบ “${p.name}” ใช่ไหม?`, [
      { text: 'ยกเลิก', style: 'cancel' },
      {
        text: 'ลบ',
        style: 'destructive',
        onPress: async () => {
          try {
            await productApi.deleteProduct(p.id);
            toast.show('ลบสินค้าแล้ว');
            refetch();
          } catch (e) {
            toast.show(errorMessage(e));
          }
        },
      },
    ]);
  };

  if (loading && !data) {
    return <LoadingScreen />;
  }
  if (error || !data) {
    return <ErrorScreen message={error ?? undefined} onRetry={refetch} />;
  }

  return (
    <ScreenWrapper>
      <HeaderBar
        title="สินค้าทั้งหมด"
        subtitle="แตะไอคอนเพื่อแก้ไขหรือลบ"
        onBack={() => navigation.navigate('Hub')}
      />

      <ScrollView contentContainerStyle={styles.body} showsVerticalScrollIndicator={false}>
        {data.length === 0 ? (
          <Text style={styles.empty}>ยังไม่มีสินค้า — กดปุ่ม “เพิ่มสินค้า” ด้านล่าง</Text>
        ) : (
          data.map((p) => (
            <View key={p.id} style={styles.row}>
              <EmojiAvatar emoji={p.emoji} size={40} />
              <View style={styles.rowMid}>
                <Text style={styles.name}>{p.name}</Text>
                <Text style={styles.price}>{baht(p.price)}</Text>
              </View>
              <Pressable style={[styles.iconBtn, styles.editBtn]} onPress={() => openEdit(p)}>
                <Text style={styles.editIcon}>✎</Text>
              </Pressable>
              <Pressable style={[styles.iconBtn, styles.delBtn]} onPress={() => confirmDelete(p)}>
                <Text style={styles.delIcon}>🗑</Text>
              </Pressable>
            </View>
          ))
        )}
      </ScrollView>

      <Pressable style={styles.fab} onPress={openAdd}>
        <Text style={styles.fabText}>＋ เพิ่มสินค้า</Text>
      </Pressable>

      <ProductEditorSheet
        visible={editorOpen}
        initial={editing}
        submitting={submitting}
        onSave={save}
        onClose={() => {
          setEditorOpen(false);
          setEditing(null);
        }}
      />
      <Toast message={toast.message} visible={toast.visible} />
    </ScreenWrapper>
  );
};

const styles = StyleSheet.create({
  body: { padding: 20, paddingBottom: 110 },
  empty: { textAlign: 'center', color: colors.textMuted, fontSize: fontSize.md, marginTop: 40 },
  row: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
    backgroundColor: colors.surface,
    borderWidth: 1,
    borderColor: colors.border,
    borderRadius: radius.lg,
    padding: 12,
    marginBottom: 10,
  },
  rowMid: { flex: 1 },
  name: { fontSize: fontSize.lg, fontWeight: '500', color: colors.textDark },
  price: { fontSize: fontSize.base, color: colors.primary, fontWeight: '600', marginTop: 1 },
  iconBtn: {
    width: 34,
    height: 34,
    borderRadius: radius.sm,
    borderWidth: 1,
    alignItems: 'center',
    justifyContent: 'center',
  },
  editBtn: { borderColor: '#ecdfce' },
  editIcon: { color: colors.primary, fontSize: 15 },
  delBtn: { borderColor: '#f6e2e2' },
  delIcon: { fontSize: 15 },
  fab: {
    position: 'absolute',
    right: 20,
    bottom: 24,
    height: 50,
    paddingHorizontal: 22,
    borderRadius: 26,
    backgroundColor: colors.primary,
    alignItems: 'center',
    justifyContent: 'center',
    shadowColor: colors.primary,
    shadowOpacity: 0.35,
    shadowRadius: 12,
    shadowOffset: { width: 0, height: 10 },
    elevation: 6,
  },
  fabText: { color: '#fff', fontSize: fontSize.lg, fontWeight: '600' },
});
