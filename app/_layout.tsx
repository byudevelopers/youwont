import { DarkTheme, DefaultTheme, ThemeProvider } from '@react-navigation/native';
import { QueryClientProvider } from '@tanstack/react-query';
import { Stack, useRouter, useSegments } from 'expo-router';
import { useEffect } from 'react';
import { ApiError } from '../api/client';
import { AuthProvider, useSession } from '../ctx';
import { useColorScheme } from '../hooks/use-color-scheme';
import { useMe } from '../hooks/use-user';
import { queryClient } from '../lib/query-client';

function RootLayoutNav() {
    const { session, isLoading } = useSession();
    const segments = useSegments();
    const router = useRouter();

    const { data: me, error: meError, isLoading: meLoading } = useMe({ enabled: !!session });

    const needsProfile = !!session && meError instanceof ApiError && meError.status === 401;

    useEffect(() => {
        if (isLoading) return;
        if (session && meLoading) return;

        const inAuthGroup = segments[0] === '(auth)';

        if (session) {
            if (needsProfile) {
                // Prevent infinite redirect if already on create-profile
                if (segments[1] !== 'create-profile') {
                    router.replace('/(auth)/create-profile');
                }
            } else if (inAuthGroup || segments.length === 0) {
                // If they are authenticated and don't need a profile, but are on an auth screen or root, go to app
                router.replace('/(tabs)');
            }
        } else if (!inAuthGroup) {
            // Not authenticated and not in auth group, go to login
            router.replace('/(auth)/login');
        }
    }, [session, isLoading, me, meError, meLoading, segments]);

    return (
        <Stack screenOptions={{ headerShown: false }}>
            <Stack.Screen name="(auth)" />
            <Stack.Screen name="(tabs)" />
            <Stack.Screen
                name="group/[id]"
                options={{
                    animation: 'slide_from_right',
                }}
            />
            <Stack.Screen
                name="bet/[id]"
                options={{
                    animation: 'slide_from_right',
                }}
            />
            <Stack.Screen
                name="create-group"
                options={{
                    animation: 'slide_from_right',
                }}
            />
            <Stack.Screen
                name="create-bet/[groupId]"
                options={{
                    animation: 'slide_from_right',
                }}
            />
            <Stack.Screen
                name="modal"
                options={{
                    presentation: 'modal',
                }}
            />
        </Stack>
    );
}

export default function RootLayout() {
    const colorScheme = useColorScheme();

    return (
        <QueryClientProvider client={queryClient}>
            <ThemeProvider value={colorScheme === 'dark' ? DarkTheme : DefaultTheme}>
                <AuthProvider>
                    <RootLayoutNav />
                </AuthProvider>
            </ThemeProvider>
        </QueryClientProvider>
    );
}