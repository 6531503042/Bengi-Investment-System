import { useState, type FC } from 'react'
import { YStack, XStack, Input, Button, Text, H3, Spinner } from 'tamagui'
import { orderService } from '@/services/api'
import { usePortfolioStore } from '@/stores/portfolio'
import { useMarketStore } from '@/stores/market'
import type { OrderSide, OrderType } from '@/types/trade'

interface OrderFormProps {
    symbol: string
    onOrderPlaced?: () => void
}

export const OrderForm: FC<OrderFormProps> = ({ symbol, onOrderPlaced }) => {
    const [side, setSide] = useState<OrderSide>('BUY')
    const [orderType, setOrderType] = useState<OrderType>('MARKET')
    const [quantity, setQuantity] = useState('')
    const [price, setPrice] = useState('')
    const [isLoading, setIsLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)

    const { selectedPortfolioId } = usePortfolioStore()
    const quote = useMarketStore((s) => s.quotes[symbol])

    const handleSubmit = async () => {
        if (!selectedPortfolioId || !quantity) return

        setIsLoading(true)
        setError(null)

        try {
            await orderService.create({
                portfolioId: selectedPortfolioId,
                symbol,
                side,
                type: orderType,
                quantity: parseFloat(quantity),
                price: orderType === 'LIMIT' ? parseFloat(price) : undefined,
            })
            onOrderPlaced?.()
        } catch (err: unknown) {
            setError((err as { response?: { data?: { error?: string } } }).response?.data?.error ?? 'Order failed')
        } finally {
            setIsLoading(false)
        }
    }

    const estimatedValue = quote?.price && quantity
        ? (quote.price * parseFloat(quantity)).toFixed(2)
        : '0.00'

    return (
        <YStack gap="$4" padding="$4">
            <H3>{symbol}</H3>
            <Text color="$gray10">Current Price: ${quote?.price?.toFixed(2) ?? '---'}</Text>

            {error && <Text color="$red10">{error}</Text>}

            {/* Side Selection */}
            <XStack gap="$2">
                <Button flex={1} theme={side === 'BUY' ? 'green' : 'gray'} onPress={() => setSide('BUY')}>
                    Buy
                </Button>
                <Button flex={1} theme={side === 'SELL' ? 'red' : 'gray'} onPress={() => setSide('SELL')}>
                    Sell
                </Button>
            </XStack>

            {/* Order Type */}
            <XStack gap="$2">
                <Button flex={1} size="$3" theme={orderType === 'MARKET' ? 'blue' : 'gray'} onPress={() => setOrderType('MARKET')}>
                    Market
                </Button>
                <Button flex={1} size="$3" theme={orderType === 'LIMIT' ? 'blue' : 'gray'} onPress={() => setOrderType('LIMIT')}>
                    Limit
                </Button>
            </XStack>

            <Input placeholder="Quantity" value={quantity} onChangeText={setQuantity} keyboardType="decimal-pad" size="$4" />

            {orderType === 'LIMIT' && (
                <Input placeholder="Limit Price" value={price} onChangeText={setPrice} keyboardType="decimal-pad" size="$4" />
            )}

            <XStack justifyContent="space-between">
                <Text color="$gray10">Estimated Value:</Text>
                <Text fontWeight="bold">${estimatedValue}</Text>
            </XStack>

            <Button
                theme={side === 'BUY' ? 'green' : 'red'}
                size="$5"
                onPress={handleSubmit}
                disabled={isLoading || !quantity}
            >
                {isLoading ? <Spinner color="white" /> : `${side} ${quantity || '0'} ${symbol}`}
            </Button>
        </YStack>
    )
}
