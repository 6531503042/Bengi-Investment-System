import { type FC } from 'react'
import { XStack, Text } from 'tamagui'
import { dimeTheme } from '@/constants/theme'

interface PriceChipProps {
    value: number
    showPercent?: boolean
    size?: 'sm' | 'md' | 'lg'
}

export const PriceChip: FC<PriceChipProps> = ({
    value,
    showPercent = true,
    size = 'md'
}) => {
    const isPositive = value >= 0
    const color = isPositive ? dimeTheme.colors.profit : dimeTheme.colors.loss
    const prefix = isPositive ? '+' : ''

    const fontSizes = {
        sm: '$2',
        md: '$3',
        lg: '$4',
    } as const

    return (
        <XStack
            alignItems="center"
            backgroundColor={`${color}20`}
            paddingHorizontal="$2"
            paddingVertical="$1"
            borderRadius="$2"
        >
            <Text
                color={color}
                fontSize={fontSizes[size]}
                fontWeight="600"
            >
                {prefix}{value.toFixed(2)}{showPercent ? '%' : ''}
            </Text>
        </XStack>
    )
}
