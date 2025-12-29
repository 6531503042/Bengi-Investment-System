import { useEffect, useState } from 'react'
import { StyleSheet, ScrollView, StatusBar, Alert, TouchableOpacity } from 'react-native'
import { YStack, XStack, Text, View, Button } from 'tamagui'
import { SafeAreaView } from 'react-native-safe-area-context'
import { useAuthStore } from '@/stores/auth'
import { useDemoStore } from '@/stores/demo'
import { dimeTheme } from '@/constants/theme'
import { Ionicons } from '@expo/vector-icons'

export default function ProfileScreen() {
    const { user, logout } = useAuthStore()
    const { account: demoAccount, isLoading, fetchDemo, deposit, reset } = useDemoStore()
    const [depositModalVisible, setDepositModalVisible] = useState(false)

    // Fetch demo account on mount
    useEffect(() => {
        fetchDemo()
    }, [])

    const handleDeposit = async () => {
        Alert.prompt(
            'Deposit Demo Funds',
            'Enter amount to deposit (USD)',
            [
                { text: 'Cancel', style: 'cancel' },
                {
                    text: 'Deposit',
                    onPress: async (value) => {
                        const amount = parseFloat(value || '0')
                        if (amount > 0) {
                            const result = await deposit(amount)
                            if (result) {
                                Alert.alert('Success', result.message)
                            }
                        }
                    },
                },
            ],
            'plain-text',
            '5000'
        )
    }

    const handleReset = () => {
        Alert.alert(
            'Reset Demo Account',
            'This will reset your demo account to $10,000. All positions and history will be cleared.',
            [
                { text: 'Cancel', style: 'cancel' },
                {
                    text: 'Reset',
                    style: 'destructive',
                    onPress: async () => {
                        const result = await reset(10000)
                        if (result) {
                            Alert.alert('Success', result.message)
                        }
                    },
                },
            ]
        )
    }

    // Settings sections with organized categories
    const settingsSections = [
        {
            title: 'Account',
            items: [
                { icon: 'person-outline', label: 'Edit Profile', value: '', hasArrow: true },
                { icon: 'mail-outline', label: 'Email', value: user?.email ?? 'Not set', hasArrow: true },
                { icon: 'key-outline', label: 'Change Password', value: '', hasArrow: true },
                { icon: 'card-outline', label: 'Payment Methods', value: '', hasArrow: true },
            ]
        },
        {
            title: 'Preferences',
            items: [
                { icon: 'moon-outline', label: 'Dark Mode', value: 'On', isToggle: true, toggleValue: true },
                { icon: 'globe-outline', label: 'Language', value: 'English', hasArrow: true },
                { icon: 'cash-outline', label: 'Currency', value: 'USD', hasArrow: true },
                { icon: 'notifications-outline', label: 'Push Notifications', value: '', isToggle: true, toggleValue: true },
                { icon: 'mail-unread-outline', label: 'Email Alerts', value: '', isToggle: true, toggleValue: false },
            ]
        },
        {
            title: 'Trading',
            items: [
                { icon: 'trending-up-outline', label: 'Default Leverage', value: `${demoAccount?.leverage ?? 1}x`, hasArrow: true },
                { icon: 'analytics-outline', label: 'Chart Style', value: 'Candles', hasArrow: true },
                { icon: 'timer-outline', label: 'Default Timeframe', value: '1D', hasArrow: true },
                { icon: 'warning-outline', label: 'Risk Management', value: '', hasArrow: true },
            ]
        },
        {
            title: 'Security',
            items: [
                { icon: 'finger-print-outline', label: 'Face ID / Touch ID', value: '', isToggle: true, toggleValue: true },
                { icon: 'lock-closed-outline', label: 'App Lock PIN', value: '', hasArrow: true },
                { icon: 'shield-checkmark-outline', label: 'Two-Factor Auth', value: 'Enabled', hasArrow: true },
                { icon: 'eye-off-outline', label: 'Hide Balances', value: '', isToggle: true, toggleValue: false },
            ]
        },
        {
            title: 'Support',
            items: [
                { icon: 'help-circle-outline', label: 'Help Center', value: '', hasArrow: true },
                { icon: 'chatbubble-outline', label: 'Contact Support', value: '', hasArrow: true },
                { icon: 'document-text-outline', label: 'Terms of Service', value: '', hasArrow: true },
                { icon: 'shield-outline', label: 'Privacy Policy', value: '', hasArrow: true },
                { icon: 'information-circle-outline', label: 'About', value: 'v1.0.0', hasArrow: true },
            ]
        }
    ]

    // For backward compatibility with existing render code
    const menuItems = settingsSections.flatMap(s => s.items)

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
                                    <XStack alignItems="center" gap="$2" marginTop="$1">
                                        <View style={styles.demoBadge}>
                                            <Text color={dimeTheme.colors.primary} fontSize="$2" fontWeight="bold">
                                                DEMO ACCOUNT
                                            </Text>
                                        </View>
                                    </XStack>
                                </YStack>
                            </XStack>
                        </View>
                    </YStack>

                    {/* Demo Balance Card */}
                    <YStack paddingHorizontal="$4" marginBottom="$4">
                        <View style={styles.balanceCard}>
                            <YStack marginBottom="$3">
                                <Text color={dimeTheme.colors.textSecondary} fontSize="$2">
                                    Demo Account Balance
                                </Text>
                                <Text color={dimeTheme.colors.textPrimary} fontSize="$9" fontWeight="bold">
                                    ${(demoAccount?.balance ?? 0).toLocaleString('en-US', { minimumFractionDigits: 2 })}
                                </Text>
                                {demoAccount && (
                                    <XStack gap="$3" marginTop="$1">
                                        <Text color={dimeTheme.colors.textTertiary} fontSize="$2">
                                            Initial: ${demoAccount.initialBalance.toLocaleString()}
                                        </Text>
                                        <Text
                                            color={demoAccount.pnlPercentage >= 0 ? dimeTheme.colors.profit : dimeTheme.colors.loss}
                                            fontSize="$2"
                                            fontWeight="600"
                                        >
                                            {demoAccount.pnlPercentage >= 0 ? '+' : ''}{demoAccount.pnlPercentage.toFixed(2)}%
                                        </Text>
                                    </XStack>
                                )}
                            </YStack>
                            <XStack gap="$3">
                                <Button
                                    flex={1}
                                    size="$4"
                                    backgroundColor={dimeTheme.colors.primary}
                                    pressStyle={{ backgroundColor: dimeTheme.colors.primaryDark }}
                                    onPress={handleDeposit}
                                    disabled={isLoading}
                                >
                                    <XStack alignItems="center" gap="$2">
                                        <Ionicons name="add-circle-outline" size={18} color={dimeTheme.colors.background} />
                                        <Text color={dimeTheme.colors.background} fontWeight="600">
                                            Deposit
                                        </Text>
                                    </XStack>
                                </Button>
                                <Button
                                    flex={1}
                                    size="$4"
                                    backgroundColor={dimeTheme.colors.surface}
                                    borderWidth={1}
                                    borderColor={dimeTheme.colors.border}
                                    onPress={handleReset}
                                    disabled={isLoading}
                                >
                                    <XStack alignItems="center" gap="$2">
                                        <Ionicons name="refresh-outline" size={18} color={dimeTheme.colors.textSecondary} />
                                        <Text color={dimeTheme.colors.textSecondary} fontWeight="600">
                                            Reset
                                        </Text>
                                    </XStack>
                                </Button>
                            </XStack>
                        </View>
                    </YStack>

                    {/* Leverage Info */}
                    {demoAccount && (
                        <YStack paddingHorizontal="$4" marginBottom="$4">
                            <View style={styles.infoCard}>
                                <XStack justifyContent="space-between" alignItems="center">
                                    <XStack alignItems="center" gap="$2">
                                        <Ionicons name="trending-up" size={20} color={dimeTheme.colors.primary} />
                                        <Text color={dimeTheme.colors.textPrimary} fontWeight="600">
                                            Default Leverage
                                        </Text>
                                    </XStack>
                                    <View style={styles.leverageBadge}>
                                        <Text color={dimeTheme.colors.primary} fontWeight="bold">
                                            {demoAccount.leverage}x
                                        </Text>
                                    </View>
                                </XStack>
                            </View>
                        </YStack>
                    )}

                    {/* Settings Sections */}
                    {settingsSections.map((section, sectionIndex) => (
                        <YStack key={section.title} paddingHorizontal="$4" marginBottom="$4">
                            <Text color={dimeTheme.colors.textSecondary} fontSize={12} fontWeight="600" marginBottom="$2" marginLeft="$1">
                                {section.title.toUpperCase()}
                            </Text>
                            <View style={styles.menuCard}>
                                {section.items.map((item, index) => (
                                    <TouchableOpacity key={item.label} activeOpacity={0.7}>
                                        <XStack
                                            paddingVertical="$3"
                                            alignItems="center"
                                            justifyContent="space-between"
                                        >
                                            <XStack alignItems="center" gap="$3">
                                                <View style={styles.settingIcon}>
                                                    <Ionicons
                                                        name={item.icon as any}
                                                        size={20}
                                                        color={dimeTheme.colors.textSecondary}
                                                    />
                                                </View>
                                                <Text color={dimeTheme.colors.textPrimary} fontSize={15}>
                                                    {item.label}
                                                </Text>
                                            </XStack>
                                            <XStack alignItems="center" gap="$2">
                                                {item.value && !item.isToggle && (
                                                    <Text color={dimeTheme.colors.textTertiary} fontSize={13}>
                                                        {item.value}
                                                    </Text>
                                                )}
                                                {item.isToggle ? (
                                                    <View style={[
                                                        styles.toggleSwitch,
                                                        item.toggleValue && styles.toggleSwitchOn
                                                    ]}>
                                                        <View style={[
                                                            styles.toggleKnob,
                                                            item.toggleValue && styles.toggleKnobOn
                                                        ]} />
                                                    </View>
                                                ) : item.hasArrow && (
                                                    <Ionicons
                                                        name="chevron-forward"
                                                        size={18}
                                                        color={dimeTheme.colors.textTertiary}
                                                    />
                                                )}
                                            </XStack>
                                        </XStack>
                                        {index < section.items.length - 1 && (
                                            <View style={styles.separator} />
                                        )}
                                    </TouchableOpacity>
                                ))}
                            </View>
                        </YStack>
                    ))}

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
            </SafeAreaView >
        </View >
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
    demoBadge: {
        backgroundColor: 'rgba(0, 230, 118, 0.15)',
        paddingHorizontal: 8,
        paddingVertical: 4,
        borderRadius: 6,
        borderWidth: 1,
        borderColor: dimeTheme.colors.primary,
    },
    balanceCard: {
        backgroundColor: dimeTheme.colors.surface,
        padding: 20,
        borderRadius: dimeTheme.radius.lg,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
    leverageBadge: {
        backgroundColor: 'rgba(0, 230, 118, 0.15)',
        paddingHorizontal: 12,
        paddingVertical: 6,
        borderRadius: 8,
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
    settingIcon: {
        width: 32,
        height: 32,
        borderRadius: 8,
        backgroundColor: 'rgba(255, 255, 255, 0.05)',
        alignItems: 'center',
        justifyContent: 'center',
    },
    toggleSwitch: {
        width: 44,
        height: 26,
        borderRadius: 13,
        backgroundColor: 'rgba(255, 255, 255, 0.1)',
        padding: 2,
        justifyContent: 'center',
    },
    toggleSwitchOn: {
        backgroundColor: dimeTheme.colors.primary,
    },
    toggleKnob: {
        width: 22,
        height: 22,
        borderRadius: 11,
        backgroundColor: '#fff',
    },
    toggleKnobOn: {
        alignSelf: 'flex-end',
    },
})
