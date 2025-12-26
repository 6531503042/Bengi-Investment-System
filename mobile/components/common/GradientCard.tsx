import { type FC } from 'react'
import { StyleSheet, type ViewStyle } from 'react-native'
import { LinearGradient } from 'expo-linear-gradient'
import { YStack } from 'tamagui'
import { dimeTheme } from '@/constants/theme'

interface GradientCardProps {
    children: React.ReactNode
    style?: ViewStyle
    variant?: 'default' | 'primary' | 'dark'
}

export const GradientCard: FC<GradientCardProps> = ({
    children,
    style,
    variant = 'default'
}) => {
    const gradients = {
        default: [dimeTheme.colors.surface, dimeTheme.colors.backgroundSecondary],
        primary: [dimeTheme.colors.primaryDark, dimeTheme.colors.primary],
        dark: [dimeTheme.colors.backgroundSecondary, dimeTheme.colors.background],
    }

    return (
        <LinearGradient
            colors={gradients[variant] as [string, string]}
            start={{ x: 0, y: 0 }}
            end={{ x: 1, y: 1 }}
            style={[styles.card, style]}
        >
            <YStack flex={1}>
                {children}
            </YStack>
        </LinearGradient>
    )
}

const styles = StyleSheet.create({
    card: {
        borderRadius: dimeTheme.radius.lg,
        padding: dimeTheme.spacing.md,
        borderWidth: 1,
        borderColor: dimeTheme.colors.border,
    },
})
