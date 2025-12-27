import { useEffect } from 'react'
import { StyleSheet, ScrollView, StatusBar } from 'react-native'
import { YStack, XStack, Text, View, Button } from 'tamagui'
import { LinearGradient } from 'expo-linear-gradient'
import { useRouter } from 'expo-router'
import { SafeAreaView } from 'react-native-safe-area-context'
import { useMarketStore } from '@/stores/market'
import { useAuthStore } from '@/stores/auth'
import { useDemoStore } from '@/stores/demo'
import { dimeTheme } from '@/constants/theme'
import { PriceChip } from '@/components/common/PriceChip'
import { Ionicons } from '@expo/vector-icons'

export default function HomeScreen() {
  const router = useRouter()
  const { user } = useAuthStore()
  const { wsConnected } = useMarketStore()
  const { account: demoAccount, isLoading, fetchDemo } = useDemoStore()

  useEffect(() => {
    fetchDemo()
  }, [])

  // Use demo account balance instead of portfolio
  const balance = demoAccount?.balance ?? 0
  const initialBalance = demoAccount?.initialBalance ?? 10000
  const pnl = balance - initialBalance
  const pnlPercent = demoAccount?.pnlPercentage ?? 0
  const leverage = demoAccount?.leverage ?? 10

  return (
    <View style={styles.container}>
      <StatusBar barStyle="light-content" />
      <SafeAreaView style={styles.safeArea} edges={['top']}>
        <ScrollView showsVerticalScrollIndicator={false}>
          {/* Header */}
          <XStack justifyContent="space-between" alignItems="center" padding="$4">
            <YStack>
              <Text color={dimeTheme.colors.textSecondary} fontSize="$3">
                Welcome back,
              </Text>
              <Text color={dimeTheme.colors.textPrimary} fontSize="$6" fontWeight="bold">
                {user?.fullName ?? 'Trader'}
              </Text>
            </YStack>
            <XStack alignItems="center" gap="$2">
              <View style={styles.demoBadge}>
                <Text color={dimeTheme.colors.primary} fontSize="$1" fontWeight="bold">
                  DEMO
                </Text>
              </View>
              <View
                width={10}
                height={10}
                borderRadius={5}
                backgroundColor={wsConnected ? dimeTheme.colors.profit : dimeTheme.colors.textTertiary}
              />
            </XStack>
          </XStack>

          {/* Portfolio Value Card */}
          <YStack paddingHorizontal="$4" marginBottom="$4">
            <LinearGradient
              colors={[dimeTheme.colors.surface, dimeTheme.colors.backgroundSecondary]}
              start={{ x: 0, y: 0 }}
              end={{ x: 1, y: 1 }}
              style={styles.portfolioCard}
            >
              <Text color={dimeTheme.colors.textSecondary} fontSize="$3" marginBottom="$2">
                Demo Account Balance
              </Text>
              <Text color={dimeTheme.colors.textPrimary} fontSize="$10" fontWeight="bold">
                ${balance.toLocaleString('en-US', { minimumFractionDigits: 2 })}
              </Text>
              <XStack marginTop="$3" alignItems="center" gap="$2">
                <PriceChip value={pnlPercent} size="md" />
                <Text color={dimeTheme.colors.textSecondary} fontSize="$2">
                  {pnl >= 0 ? '+' : ''}${pnl.toFixed(2)} total P&L
                </Text>
              </XStack>
            </LinearGradient>
          </YStack>

          {/* Quick Stats */}
          <XStack paddingHorizontal="$4" gap="$3" marginBottom="$4">
            <View style={styles.statCard}>
              <XStack alignItems="center" gap="$2" marginBottom="$1">
                <Ionicons name="trending-up" size={16} color={dimeTheme.colors.primary} />
                <Text color={dimeTheme.colors.textSecondary} fontSize="$2">
                  Leverage
                </Text>
              </XStack>
              <Text color={dimeTheme.colors.textPrimary} fontSize="$5" fontWeight="bold">
                {leverage}x
              </Text>
            </View>
            <View style={styles.statCard}>
              <XStack alignItems="center" gap="$2" marginBottom="$1">
                <Ionicons name="wallet-outline" size={16} color={dimeTheme.colors.primary} />
                <Text color={dimeTheme.colors.textSecondary} fontSize="$2">
                  Initial
                </Text>
              </XStack>
              <Text color={dimeTheme.colors.textPrimary} fontSize="$5" fontWeight="bold">
                ${initialBalance.toLocaleString()}
              </Text>
            </View>
            <View style={styles.statCard}>
              <XStack alignItems="center" gap="$2" marginBottom="$1">
                <Ionicons name="analytics-outline" size={16} color={dimeTheme.colors.primary} />
                <Text color={dimeTheme.colors.textSecondary} fontSize="$2">
                  Positions
                </Text>
              </XStack>
              <Text color={dimeTheme.colors.textPrimary} fontSize="$5" fontWeight="bold">
                0
              </Text>
            </View>
          </XStack>

          {/* Quick Actions */}
          <YStack paddingHorizontal="$4" marginBottom="$4">
            <Text color={dimeTheme.colors.textPrimary} fontSize="$5" fontWeight="bold" marginBottom="$3">
              Quick Actions
            </Text>
            <XStack gap="$3">
              <Button
                flex={1}
                size="$5"
                backgroundColor={dimeTheme.colors.primary}
                pressStyle={{ backgroundColor: dimeTheme.colors.primaryDark }}
                onPress={() => router.push('/(tabs)/market/index')}
              >
                <XStack alignItems="center" gap="$2">
                  <Ionicons name="search" size={18} color={dimeTheme.colors.background} />
                  <Text color={dimeTheme.colors.background} fontWeight="bold">
                    Market
                  </Text>
                </XStack>
              </Button>
              <Button
                flex={1}
                size="$5"
                backgroundColor={dimeTheme.colors.surface}
                borderWidth={1}
                borderColor={dimeTheme.colors.profit}
                pressStyle={{ backgroundColor: dimeTheme.colors.backgroundSecondary }}
                onPress={() => router.push('/trade')}
              >
                <XStack alignItems="center" gap="$2">
                  <Ionicons name="swap-vertical" size={18} color={dimeTheme.colors.profit} />
                  <Text color={dimeTheme.colors.profit} fontWeight="bold">
                    Trade
                  </Text>
                </XStack>
              </Button>
            </XStack>
          </YStack>

          {/* Demo Info */}
          <YStack paddingHorizontal="$4" paddingBottom="$4">
            <View style={styles.infoCard}>
              <XStack alignItems="center" gap="$2" marginBottom="$2">
                <Ionicons name="information-circle" size={20} color={dimeTheme.colors.primary} />
                <Text color={dimeTheme.colors.textPrimary} fontWeight="600">
                  Demo Trading Mode
                </Text>
              </XStack>
              <Text color={dimeTheme.colors.textSecondary} fontSize="$2">
                Trade with virtual funds. Your demo balance can be reset or topped up anytime from your Profile.
              </Text>
            </View>
          </YStack>
        </ScrollView>
      </SafeAreaView>
    </View>
  )
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: dimeTheme.colors.background,
  },
  safeArea: {
    flex: 1,
  },
  portfolioCard: {
    padding: 20,
    borderRadius: dimeTheme.radius.xl,
    borderWidth: 1,
    borderColor: dimeTheme.colors.border,
  },
  demoBadge: {
    backgroundColor: 'rgba(0, 230, 118, 0.15)',
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 6,
    borderWidth: 1,
    borderColor: dimeTheme.colors.primary,
  },
  statCard: {
    flex: 1,
    backgroundColor: dimeTheme.colors.surface,
    padding: 12,
    borderRadius: dimeTheme.radius.lg,
    borderWidth: 1,
    borderColor: dimeTheme.colors.border,
  },
  infoCard: {
    backgroundColor: dimeTheme.colors.surface,
    padding: 16,
    borderRadius: dimeTheme.radius.lg,
    borderWidth: 1,
    borderColor: dimeTheme.colors.border,
  },
})
