import { useState, type FC } from 'react'
import { StyleSheet } from 'react-native'
import { YStack, XStack, Input, Button, Text, Spinner, View } from 'tamagui'
import { LinearGradient } from 'expo-linear-gradient'
import { useAuthStore } from '@/stores/auth'

interface LoginFormProps {
    onRegisterPress?: () => void
}

export const LoginForm: FC<LoginFormProps> = ({ onRegisterPress }) => {
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const { login, isLoading, error, clearError } = useAuthStore()

    const handleLogin = async () => {
        if (!email || !password) return
        clearError()
        await login(email, password)
    }

    return (
        <LinearGradient
            colors={['#0a0a0a', '#1a1a2e', '#16213e']}
            start={{ x: 0, y: 0 }}
            end={{ x: 1, y: 1 }}
            style={styles.container}
        >
            <YStack flex={1} padding="$6" justifyContent="center">
                {/* Logo Area */}
                <YStack alignItems="center" marginBottom="$8">
                    <View
                        width={80}
                        height={80}
                        borderRadius={40}
                        backgroundColor="$green10"
                        alignItems="center"
                        justifyContent="center"
                        marginBottom="$4"
                    >
                        <Text fontSize={36} fontWeight="bold" color="white">B</Text>
                    </View>
                    <Text fontSize="$8" fontWeight="bold" color="white">
                        Welcome Back
                    </Text>
                    <Text fontSize="$3" color="$gray10" marginTop="$2">
                        Sign in to continue trading
                    </Text>
                </YStack>

                {/* Error Message */}
                {error && (
                    <View
                        backgroundColor="rgba(239, 68, 68, 0.15)"
                        borderRadius="$3"
                        padding="$3"
                        marginBottom="$4"
                        borderWidth={1}
                        borderColor="$red10"
                    >
                        <Text color="$red10" textAlign="center" fontSize="$3">
                            {error}
                        </Text>
                    </View>
                )}

                {/* Form */}
                <YStack gap="$4">
                    <YStack gap="$2">
                        <Text color="$gray11" fontSize="$2" marginLeft="$1">
                            Email Address
                        </Text>
                        <Input
                            placeholder="you@example.com"
                            value={email}
                            onChangeText={setEmail}
                            autoCapitalize="none"
                            keyboardType="email-address"
                            size="$5"
                            backgroundColor="rgba(255,255,255,0.08)"
                            borderColor="rgba(255,255,255,0.15)"
                            color="white"
                            placeholderTextColor="$gray10"
                            borderRadius="$4"
                        />
                    </YStack>

                    <YStack gap="$2">
                        <Text color="$gray11" fontSize="$2" marginLeft="$1">
                            Password
                        </Text>
                        <Input
                            placeholder="••••••••"
                            value={password}
                            onChangeText={setPassword}
                            secureTextEntry
                            size="$5"
                            backgroundColor="rgba(255,255,255,0.08)"
                            borderColor="rgba(255,255,255,0.15)"
                            color="white"
                            placeholderTextColor="$gray10"
                            borderRadius="$4"
                        />
                    </YStack>

                    <Button
                        onPress={handleLogin}
                        size="$5"
                        disabled={isLoading || !email || !password}
                        marginTop="$2"
                        borderRadius="$4"
                        backgroundColor={isLoading || !email || !password ? '$gray8' : '$green9'}
                        pressStyle={{ backgroundColor: '$green10', scale: 0.98 }}
                    >
                        {isLoading ? (
                            <Spinner color="white" />
                        ) : (
                            <Text color="white" fontWeight="bold" fontSize="$4">
                                Sign In
                            </Text>
                        )}
                    </Button>
                </YStack>

                {/* Footer */}
                <XStack justifyContent="center" gap="$2" marginTop="$6">
                    <Text color="$gray10" fontSize="$3">
                        Don't have an account?
                    </Text>
                    <Text
                        color="$green10"
                        onPress={onRegisterPress}
                        fontWeight="bold"
                        fontSize="$3"
                    >
                        Sign Up
                    </Text>
                </XStack>
            </YStack>
        </LinearGradient>
    )
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
    },
})
