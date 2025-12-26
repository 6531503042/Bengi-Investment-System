import { SafeAreaView } from 'react-native'
import { useRouter } from 'expo-router'
import { RegisterForm } from '@/components/auth/RegisterForm'

export default function RegisterScreen() {
    const router = useRouter()

    return (
        <SafeAreaView style={{ flex: 1 }}>
            <RegisterForm onLoginPress={() => router.back()} />
        </SafeAreaView>
    )
}
