import React, { useState } from 'react';
import { ActivityIndicator, ScrollView, StyleSheet, Text, View } from 'react-native';
import { aiApi, errorMessage } from '../api';
import { AiSummary } from '../api/types';
import { Card, PrimaryButton, ScreenWrapper } from '../components';
import { colors, fontSize, radius } from '../theme';

export const AiScreen: React.FC = () => {
  const [data, setData] = useState<AiSummary | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const run = async () => {
    if (loading) {
      return;
    }
    setLoading(true);
    setError(null);
    try {
      setData(await aiApi.aiSummary());
    } catch (e) {
      setError(errorMessage(e));
    } finally {
      setLoading(false);
    }
  };

  return (
    <ScreenWrapper>
      <View style={styles.header}>
        <View style={styles.headRow}>
          <Text style={styles.sparkle}>✨</Text>
          <Text style={styles.title}>AI สรุปประจำวัน</Text>
        </View>
        <Text style={styles.subtitle}>วิเคราะห์อัตโนมัติจากยอดขายของคุณ</Text>
      </View>

      <ScrollView contentContainerStyle={styles.body} showsVerticalScrollIndicator={false}>
        {loading ? (
          <View style={styles.center}>
            <ActivityIndicator size="large" color={colors.primary} />
            <Text style={styles.loadingText}>AI กำลังวิเคราะห์ยอดขาย...</Text>
          </View>
        ) : error ? (
          <Card>
            <Text style={styles.errorEmoji}>😕</Text>
            <Text style={styles.errorText}>{error}</Text>
            <PrimaryButton label="ลองอีกครั้ง" onPress={run} style={styles.retry} />
          </Card>
        ) : data ? (
          <View style={styles.stack}>
            {/* สรุปยอดขาย */}
            <Card>
              <View style={styles.cardHead}>
                <Text style={styles.cardIcon}>📈</Text>
                <Text style={styles.cardTitle}>สรุปยอดขายวันนี้</Text>
              </View>
              <Text style={styles.summaryText}>{data.summary}</Text>
            </Card>

            {/* สินค้าขายดี */}
            <Card>
              <View style={styles.cardHead}>
                <Text style={styles.cardIcon}>⭐</Text>
                <Text style={styles.cardTitle}>สินค้าขายดี</Text>
              </View>
              {data.best_sellers.length === 0 ? (
                <Text style={styles.muted}>ยังไม่มีข้อมูลการขาย</Text>
              ) : (
                data.best_sellers.map((b, i) => (
                  <View key={i} style={styles.bestRow}>
                    <Text style={styles.bestEmoji}>{b.emoji}</Text>
                    <Text style={styles.bestName}>{b.name}</Text>
                    <Text style={styles.bestQty}>{b.qty} แก้ว</Text>
                  </View>
                ))
              )}
            </Card>

            {/* คำแนะนำ */}
            <View style={styles.tipCard}>
              <View style={styles.cardHead}>
                <Text style={styles.cardIcon}>💡</Text>
                <Text style={[styles.cardTitle, { color: '#fff' }]}>คำแนะนำสำหรับร้าน</Text>
              </View>
              {data.recommendations.map((r, i) => (
                <View key={i} style={styles.tipRow}>
                  <Text style={styles.tipNum}>{i + 1}.</Text>
                  <Text style={styles.tipText}>{r}</Text>
                </View>
              ))}
            </View>

            <PrimaryButton label="สรุปใหม่อีกครั้ง" onPress={run} style={styles.refresh} />
          </View>
        ) : (
          <Card>
            <Text style={styles.introEmoji}>🤖</Text>
            <Text style={styles.introTitle}>ให้ AI ช่วยสรุปยอดขายวันนี้</Text>
            <Text style={styles.introText}>
              วิเคราะห์ยอดขาย สินค้าขายดี และคำแนะนำเพื่อเพิ่มยอดให้ร้านคุณ
            </Text>
            <PrimaryButton label="✨ สรุปวันนี้" onPress={run} style={styles.retry} />
          </Card>
        )}
      </ScrollView>
    </ScreenWrapper>
  );
};

const styles = StyleSheet.create({
  header: {
    backgroundColor: colors.primary,
    paddingHorizontal: 20,
    paddingTop: 8,
    paddingBottom: 22,
    borderBottomLeftRadius: radius.xxl,
    borderBottomRightRadius: radius.xxl,
  },
  headRow: { flexDirection: 'row', alignItems: 'center', gap: 8 },
  sparkle: { fontSize: 18 },
  title: { fontSize: fontSize.xxl, color: '#fff', fontWeight: '600' },
  subtitle: { fontSize: fontSize.base, color: 'rgba(255,255,255,0.7)', marginTop: 2 },
  body: { padding: 20 },
  center: { alignItems: 'center', justifyContent: 'center', paddingVertical: 60, gap: 14 },
  loadingText: { fontSize: fontSize.md, color: colors.textMuted },
  stack: { gap: 14 },
  cardHead: { flexDirection: 'row', alignItems: 'center', gap: 8, marginBottom: 8 },
  cardIcon: { fontSize: 15 },
  cardTitle: { fontSize: fontSize.md, fontWeight: '600', color: colors.textDark },
  summaryText: { fontSize: fontSize.md, lineHeight: 24, color: '#6b5544' },
  muted: { color: colors.textMuted, fontSize: fontSize.md },
  bestRow: { flexDirection: 'row', alignItems: 'center', gap: 10, paddingVertical: 5 },
  bestEmoji: { fontSize: 18 },
  bestName: { flex: 1, fontSize: fontSize.md, color: colors.textDark },
  bestQty: { fontSize: fontSize.base, color: colors.textMuted },
  tipCard: { backgroundColor: '#c8895c', borderRadius: radius.xl, padding: 16 },
  tipRow: { flexDirection: 'row', gap: 9, marginBottom: 10 },
  tipNum: { color: 'rgba(255,255,255,0.9)', fontSize: fontSize.base },
  tipText: { flex: 1, color: 'rgba(255,255,255,0.9)', fontSize: fontSize.base, lineHeight: 21 },
  refresh: { marginTop: 4 },
  retry: { marginTop: 16 },
  introEmoji: { fontSize: 40, textAlign: 'center' },
  introTitle: {
    fontSize: fontSize.lg,
    fontWeight: '600',
    color: colors.textDark,
    textAlign: 'center',
    marginTop: 10,
  },
  introText: {
    fontSize: fontSize.md,
    color: colors.textMuted,
    textAlign: 'center',
    marginTop: 6,
    lineHeight: 21,
  },
  errorEmoji: { fontSize: 36, textAlign: 'center' },
  errorText: { fontSize: fontSize.md, color: colors.textMuted, textAlign: 'center', marginTop: 8 },
});
