import { StyleSheet, ScrollView, StatusBar } from 'react-native'
import { YStack, XStack, Text, View, Button, Separator } from 'tamagui'
import { SafeAreaView } from 'react-native-safe-area-context'
import { useAuthStore } from '@/stores/auth'
import { usePortfolioStore } from '@/stores/portfolio'
import { dimeTheme } from '@/constants/theme'
import { Ionicons } from '@expo/vector-icons'

export default function ProfileScreen() {
    const { user, logout } = useAuthStore()
    const { accounts = [] } = usePortfolioStore()

    const safeAccounts = accounts ?? []
    const totalBalance = safeAccounts.reduce((sum, a) => sum + (a.balance ?? 0), 0)

    const menuItems = [
        { icon: 'person-outline', label: 'Account Settings', value: '' },
        { icon: 'notifications-outline', label: 'Notifications', value: 'On' },
        { icon: 'moon-outline', label: 'Theme', value: 'Dark' },
        { icon: 'shield-checkmark-outline', label: 'Security', value: '' },
        { icon: 'help-circle-outline', label: 'Help & Support', value: '' },
        { icon: 'document-text-outline', label: 'Terms & Privacy', value: '' },
    ]

    return (
        <View style={styles.container}>
            <StatusBar barStyle="light-content" />
            <SafeAreaView style={styles.safeArea} edges={['top']}>
                <ScrollView showsVerticalScrollIndicator={false}>
                    {/* Header */}
                    <YStack padding="$4" paddingBottom="$2">
                        <Text color={dimeTheme.colors.textPrimary} fontSize="$8" fontWeight="bold">
                            Profile
                        </Text>
                    </YStack>

                    {/* User Card */}
                    <YStack paddingHorizontal="$4" marginBottom="$4">
                        <View style={styles.userCard}>
                            <XStack alignItems="center" gap="$4">
                                <View style={styles.avatar}>
                                    <Text color={dimeTheme.colors.background} fontSize="$6" fontWeight="bold">
                                        {user?.fullName?.charAt(0).toUpperCase() ?? 'U'}
                                    </Text>
                                </View>
                                <YStack flex={1}>
                                    <Text color={dimeTheme.colors.textPrimary} fontSize="$5" fontWeight="bold">
                                        {user?.fullName ?? 'User'}
                                    </Text>
                                    <Text color={dimeTheme.colors.textSecondary} fontSize="$3">
                                        {user?.email ?? 'email@example.com'}
                                    </Text>
                                </YStack>
                                <Button size="$3" circular backgroundColor={dimeTheme.colors.surface}>
                                    <Ionicons name="pencil" size={16} color={dimeTheme.colors.primary} />
                                </Button>
                            </XStack>
                        </View>
                    </YStack>

                    {/* Balance Card */}
                    <YStack paddingHorizontal="$4" marginBottom="$4">
                        <View style={styles.balanceCard}>
                            <XStack justifyContent="space-between" alignItems="center">
                                <YStack>
                                    <Text color={dimeTheme.colors.textSecondary} fontSize="$2">
                                        Available Balance
                                    </Text>
                                    <Text color={dimeTheme.colors.textPrimary} fontSize="$7" fontWeight="bold">
                                        ${totalBalance.toFixed(2)}
                                    </Text>
                                </YStack>
                                <XStack gap="$2">
                                    <Button
                                        size="$4"
                                        backgroundColor={dimeTheme.colors.primary}
                                        pressStyle={{ backgroundColor: dimeTheme.colors.primaryDark }}
                                    >
                                        <Text color={dimeTheme.colors.background} fontWeight="600">
                                            Deposit
                                        </Text>
                                    </Button>
                                    <Button
                                        size="$4"
                                        backgroundColor={dimeTheme.colors.surface}
                                        borderWidth={1}
                                        borderColor={dimeTheme.colors.primary}
                                    >
                                        <Text color={dimeTheme.colors.primary} fontWeight="600">
                                            Withdraw
                                        </Text>
                                    </Button>
                                </XStack>
                            </XStack>
                        </View>
                    </YStack>

                    {/* Menu Items */}
                    <YStack paddingHorizontal="$4" marginBottom="$4">
                        <Text color={dimeTheme.colors.textPrimary} fontSize="$5" fontWeight="bold" marginBottom="$3">
                            Settings
                        </Text>
                        <View style={styles.menuCard}>
                            {menuItems.map((item, index) => (
                                <View key={item.label}>
                                    <XStack
                                        paddingVertical="$3"
                                        alignItems="center"
                                        justifyContent="space-between"
                                    >
                                        <XStack alignItems="center" gap="$3">
                                            <Ionicons
                                                name={item.icon as any}
                                                size={22}
                                                color={dimeTheme.colors.textSecondary}
                                            />
                                            <Text color={dimeTheme.colors.textPrimary} fontSize="$4">
                                                {item.label}
                                            </Text>
                                        </XStack>
                                        <XStack alignItems="center" gap="$2">
                                            {item.value && (
                                                <Text color={dimeTheme.colors.textTertiary} fontSize="$3">
                                                    {item.value}
                                                </Text>
                                            )}
                                            <Ionicons
                                                name="chevron-forward"
                                                size={18}
                                                color={dimeTheme.colors.textTertiary}
                                            />
                                        </XStack>
                                    </XStack>
                                    {index < menuItems.length - 1 && (
                                        <View style={styles.separator} />
                                    )}
                                </View>
                            ))}
                        </View>
                    </YStack>

                    {/* App Info */}
                    <YStack paddingHorizontal="$4" marginBottom="$4">
                        <View style={styles.infoCard}>
                            <XStack justifyContent="space-between">
                                <Text color={dimeTheme.colors.textTertiary}>Version</Text>
                                <Text color={dimeTheme.colors.textSecondary}>1.0.0</Text>
                            </XStack>
                        </View>
                    </YStack>

                    {/* Logout Button */}
                    <YStack paddingHorizontal="$4" paddingBottom="$8">
                        <Button
                            size="$5"
                            backgroundColor={dimeTheme.colors.surface}
                            borderWidth={1}
                            borderColor={dimeTheme.colors.loss}
                            onPress={logout}
                        >
                            <XStack alignItems="center" gap="$2">
                                <Ionicons name="log-out-outline" size={20} color={dimeTheme.colors.loss} />
                                <Text color={dimeTheme.colors.loss} fontWeight="bold" fontSize="$4">
                                    Sign Out
                                </Text>
                            </XStack>
                        </Button>
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
    userCard: {
        backgroundColor: dimeTheme.colors.surface,
        padding: 20,
        borderRadius: dimeTheme.radius.xl,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    avatar: {
        width: 60,
        height: 60,
        borderRadius: 30,
        backgroundColor: dimeTheme.colors.primary,
        alignItems: 'center',
        justifyContent: 'center',
    },
    balanceCard: {
        backgroundColor: dimeTheme.colors.surface,
        padding: 20,
        borderRadius: dimeTheme.radius.lg,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    menuCard: {
        backgroundColor: dimeTheme.colors.surface,
        padding: 16,
        borderRadius: dimeTheme.radius.lg,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    separator: {
        height: 1,
        backgroundColor: dimeTheme.colors.border,
    },
    infoCard: {
        backgroundColor: dimeTheme.colors.surface,
        padding: 16,
        borderRadius: dimeTheme.radius.lg,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
})
