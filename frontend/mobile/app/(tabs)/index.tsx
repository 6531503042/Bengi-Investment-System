import { useEffect } from 'react'
import { StyleSheet, ScrollView, StatusBar } from 'react-native'
import { YStack, XStack, Text, View, Button } from 'tamagui'
import { LinearGradient } from 'expo-linear-gradient'
import { useRouter } from 'expo-router'
import { SafeAreaView } from 'react-native-safe-area-context'
import { usePortfolioStore } from '@/stores/portfolio'
import { useMarketStore } from '@/stores/market'
import { useAuthStore } from '@/stores/auth'
import { dimeTheme } from '@/constants/theme'
import { PriceChip } from '@/components/common/PriceChip'
import type { Portfolio, Account } from '@/types/portfolio'

export default function HomeScreen() {
  const router = useRouter()
  const { user } = useAuthStore()
  const { portfolios = [], accounts = [], fetchPortfolios, fetchAccounts } = usePortfolioStore()
  const { wsConnected } = useMarketStore()

  useEffect(() => {
    fetchPortfolios()
    fetchAccounts()
  }, [])

  const safePortfolios = portfolios ?? []
  const safeAccounts = accounts ?? []

  const totalValue = safePortfolios.reduce((sum: number, p: Portfolio) => sum + (p.totalValue ?? 0), 0)
  const totalPL = safePortfolios.reduce((sum: number, p: Portfolio) => sum + (p.totalPL ?? 0), 0)
  const totalPLPercent = totalValue > 0 ? (totalPL / totalValue) * 100 : 0
  const cashBalance = safeAccounts.reduce((sum: number, a: Account) => sum + (a.balance ?? 0), 0)

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
              <View
                width={10}
                height={10}
                borderRadius={5}
                backgroundColor={wsConnected ? dimeTheme.colors.profit : dimeTheme.colors.loss}
              />
              <Text color={wsConnected ? dimeTheme.colors.profit : dimeTheme.colors.loss} fontSize="$2">
                {wsConnected ? 'Live' : 'Offline'}
              </Text>
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
                Total Portfolio Value
              </Text>
              <Text color={dimeTheme.colors.textPrimary} fontSize="$10" fontWeight="bold">
                ${totalValue.toLocaleString('en-US', { minimumFractionDigits: 2 })}
              </Text>
              <XStack marginTop="$3" alignItems="center" gap="$2">
                <PriceChip value={totalPLPercent} size="md" />
                <Text color={dimeTheme.colors.textSecondary} fontSize="$2">
                  {totalPL >= 0 ? '+' : ''}${totalPL.toFixed(2)} today
                </Text>
              </XStack>
            </LinearGradient>
          </YStack>

          {/* Quick Stats */}
          <XStack paddingHorizontal="$4" gap="$3" marginBottom="$4">
            <View style={styles.statCard}>
              <Text color={dimeTheme.colors.textSecondary} fontSize="$2">
                Cash Balance
              </Text>
              <Text color={dimeTheme.colors.textPrimary} fontSize="$5" fontWeight="bold">
                ${cashBalance.toFixed(2)}
              </Text>
            </View>
            <View style={styles.statCard}>
              <Text color={dimeTheme.colors.textSecondary} fontSize="$2">
                Positions
              </Text>
              <Text color={dimeTheme.colors.textPrimary} fontSize="$5" fontWeight="bold">
                {safePortfolios.reduce((sum, p) => sum + (p.positions?.length ?? 0), 0)}
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
                <Text color={dimeTheme.colors.background} fontWeight="bold">
                  Browse Market
                </Text>
              </Button>
              <Button
                flex={1}
                size="$5"
                backgroundColor={dimeTheme.colors.surface}
                borderWidth={1}
                borderColor={dimeTheme.colors.primary}
                pressStyle={{ backgroundColor: dimeTheme.colors.backgroundSecondary }}
                onPress={() => router.push('/trade')}
              >
                <Text color={dimeTheme.colors.primary} fontWeight="bold">
                  Place Order
                </Text>
              </Button>
            </XStack>
          </YStack>

          {/* Recent Activity */}
          <YStack paddingHorizontal="$4" paddingBottom="$4">
            <Text color={dimeTheme.colors.textPrimary} fontSize="$5" fontWeight="bold" marginBottom="$3">
              Recent Activity
            </Text>
            <View style={styles.emptyCard}>
              <Text color={dimeTheme.colors.textTertiary} textAlign="center">
                No recent activity
              </Text>
              <Text color={dimeTheme.colors.textSecondary} fontSize="$2" textAlign="center" marginTop="$1">
                Start trading to see your activity here
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
  statCard: {
    flex: 1,
    backgroundColor: dimeTheme.colors.surface,
    padding: 16,
    borderRadius: dimeTheme.radius.lg,
    borderWidth: 1,
    borderColor: dimeTheme.colors.border,
  },
  emptyCard: {
    backgroundColor: dimeTheme.colors.surface,
    padding: 24,
    borderRadius: dimeTheme.radius.lg,
    borderWidth: 1,
    borderColor: dimeTheme.colors.border,
  },
})
