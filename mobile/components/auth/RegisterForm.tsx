import { useState, type FC } from 'react'
import { StyleSheet } from 'react-native'
import { YStack, XStack, Input, Button, Text, Spinner, View } from 'tamagui'
import { LinearGradient } from 'expo-linear-gradient'
import { useAuthStore } from '@/stores/auth'

interface RegisterFormProps {
    onLoginPress?: () => void
}

export const RegisterForm: FC<RegisterFormProps> = ({ onLoginPress }) => {
    const [fullName, setFullName] = useState('')
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [confirmPassword, setConfirmPassword] = useState('')
    const { register, isLoading, error, clearError } = useAuthStore()

    const handleRegister = async () => {
        if (password !== confirmPassword) return
        clearError()
        await register(email, password, fullName)
    }

    const isValid = fullName && email && password && password === confirmPassword && password.length >= 8
    const passwordMismatch = confirmPassword.length > 0 && password !== confirmPassword

    return (
        <LinearGradient
            colors={['#0a0a0a', '#1a1a2e', '#16213e']}
            start={{ x: 0, y: 0 }}
            end={{ x: 1, y: 1 }}
            style={styles.container}
        >
            <YStack flex={1} padding="$6" justifyContent="center">
                {/* Logo Area */}
                <YStack alignItems="center" marginBottom="$6">
                    <View
                        width={70}
                        height={70}
                        borderRadius={35}
                        backgroundColor="$green10"
                        alignItems="center"
                        justifyContent="center"
                        marginBottom="$3"
                    >
                        <Text fontSize={32} fontWeight="bold" color="white">B</Text>
                    </View>
                    <Text fontSize="$7" fontWeight="bold" color="white">
                        Create Account
                    </Text>
                    <Text fontSize="$3" color="$gray10" marginTop="$2">
                        Start your trading journey
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
                <YStack gap="$3">
                    <YStack gap="$1">
                        <Text color="$gray11" fontSize="$2" marginLeft="$1">
                            Full Name
                        </Text>
                        <Input
                            placeholder="John Doe"
                            value={fullName}
                            onChangeText={setFullName}
                            size="$4"
                            backgroundColor="rgba(255,255,255,0.08)"
                            borderColor="rgba(255,255,255,0.15)"
                            color="white"
                            placeholderTextColor="$gray10"
                            borderRadius="$4"
                        />
                    </YStack>

                    <YStack gap="$1">
                        <Text color="$gray11" fontSize="$2" marginLeft="$1">
                            Email Address
                        </Text>
                        <Input
                            placeholder="you@example.com"
                            value={email}
                            onChangeText={setEmail}
                            autoCapitalize="none"
                            keyboardType="email-address"
                            size="$4"
                            backgroundColor="rgba(255,255,255,0.08)"
                            borderColor="rgba(255,255,255,0.15)"
                            color="white"
                            placeholderTextColor="$gray10"
                            borderRadius="$4"
                        />
                    </YStack>

                    <YStack gap="$1">
                        <Text color="$gray11" fontSize="$2" marginLeft="$1">
                            Password (min 8 characters)
                        </Text>
                        <Input
                            placeholder="••••••••"
                            value={password}
                            onChangeText={setPassword}
                            secureTextEntry
                            size="$4"
                            backgroundColor="rgba(255,255,255,0.08)"
                            borderColor="rgba(255,255,255,0.15)"
                            color="white"
                            placeholderTextColor="$gray10"
                            borderRadius="$4"
                        />
                    </YStack>

                    <YStack gap="$1">
                        <Text color="$gray11" fontSize="$2" marginLeft="$1">
                            Confirm Password
                        </Text>
                        <Input
                            placeholder="••••••••"
                            value={confirmPassword}
                            onChangeText={setConfirmPassword}
                            secureTextEntry
                            size="$4"
                            backgroundColor="rgba(255,255,255,0.08)"
                            borderColor={passwordMismatch ? '$red10' : 'rgba(255,255,255,0.15)'}
                            color="white"
                            placeholderTextColor="$gray10"
                            borderRadius="$4"
                        />
                        {passwordMismatch && (
                            <Text color="$red10" fontSize="$2" marginLeft="$1" marginTop="$1">
                                Passwords don't match
                            </Text>
                        )}
                    </YStack>

                    <Button
                        onPress={handleRegister}
                        size="$5"
                        disabled={isLoading || !isValid}
                        marginTop="$3"
                        borderRadius="$4"
                        backgroundColor={isLoading || !isValid ? '$gray8' : '$green9'}
                        pressStyle={{ backgroundColor: '$green10', scale: 0.98 }}
                    >
                        {isLoading ? (
                            <Spinner color="white" />
                        ) : (
                            <Text color="white" fontWeight="bold" fontSize="$4">
                                Create Account
                            </Text>
                        )}
                    </Button>
                </YStack>

                {/* Footer */}
                <XStack justifyContent="center" gap="$2" marginTop="$5">
                    <Text color="$gray10" fontSize="$3">
                        Already have an account?
                    </Text>
                    <Text
                        color="$green10"
                        onPress={onLoginPress}
                        fontWeight="bold"
                        fontSize="$3"
                    >
                        Sign In
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
