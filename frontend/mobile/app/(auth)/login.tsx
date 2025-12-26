import { SafeAreaView } from 'react-native'
import { useRouter } from 'expo-router'
import { LoginForm } from '@/components/auth/LoginForm'

export default function LoginScreen() {
    const router = useRouter()

    return (
        <SafeAreaView style={{ flex: 1 }}>
            <LoginForm onRegisterPress={() => router.push('/(auth)/register')} />
        </SafeAreaView>
    )
}
