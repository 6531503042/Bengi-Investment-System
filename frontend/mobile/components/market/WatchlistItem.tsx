import { type FC } from 'react'
import { TouchableOpacity, StyleSheet } from 'react-native'
import { XStack, YStack, Text, View } from 'tamagui'
import { useRouter } from 'expo-router'
import { PriceChip } from '@/components/common/PriceChip'
import { dimeTheme } from '@/constants/theme'

interface WatchlistItemProps {
    symbol: string
    name: string
    price: number
    change: number
    changePercent: number
    onPress?: () => void
}

export const WatchlistItem: FC<WatchlistItemProps> = ({
    symbol,
    name,
    price,
    change,
    changePercent,
    onPress,
}) => {
    const router = useRouter()
    const isPositive = change >= 0

    const handlePress = () => {
        if (onPress) {
            onPress()
        } else {
            router.push(`/(tabs)/market/${symbol}`)
        }
    }

    return (
        <TouchableOpacity onPress={handlePress} activeOpacity={0.7}>
            <XStack
                backgroundColor={dimeTheme.colors.surface}
                padding="$4"
                borderRadius="$4"
                marginBottom="$3"
                alignItems="center"
                justifyContent="space-between"
                borderWidth={1}
                borderColor={dimeTheme.colors.border}
            >
                {/* Left: Symbol & Name */}
                <XStack alignItems="center" gap="$3" flex={1}>
                    {/* Symbol Icon */}
                    <View
                        width={44}
                        height={44}
                        borderRadius={22}
                        backgroundColor={dimeTheme.colors.backgroundSecondary}
                        alignItems="center"
                        justifyContent="center"
                    >
                        <Text
                            color={dimeTheme.colors.primary}
                            fontWeight="bold"
                            fontSize="$4"
                        >
                            {symbol.charAt(0)}
                        </Text>
                    </View>

                    <YStack>
                        <Text
                            color={dimeTheme.colors.textPrimary}
                            fontWeight="bold"
                            fontSize="$4"
                        >
                            {symbol}
                        </Text>
                        <Text
                            color={dimeTheme.colors.textSecondary}
                            fontSize="$2"
                            numberOfLines={1}
                        >
                            {name}
                        </Text>
                    </YStack>
                </XStack>

                {/* Right: Price & Change */}
                <YStack alignItems="flex-end" gap="$1">
                    <Text
                        color={dimeTheme.colors.textPrimary}
                        fontWeight="bold"
                        fontSize="$5"
                    >
                        ${price.toFixed(2)}
                    </Text>
                    <PriceChip value={changePercent} size="sm" />
                </YStack>
            </XStack>
        </TouchableOpacity>
    )
}

const styles = StyleSheet.create({})
