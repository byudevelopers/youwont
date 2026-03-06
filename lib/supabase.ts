import 'react-native-url-polyfill/auto';
import * as SecureStore from 'expo-secure-store';
import { Platform } from 'react-native';
import { createClient } from '@supabase/supabase-js';

// SecureStore adapter for Supabase persistence (native only)
const ExpoSecureStoreAdapter = {
    getItem: (key: string) => {
        return SecureStore.getItemAsync(key);
    },
    setItem: (key: string, value: string) => {
        SecureStore.setItemAsync(key, value);
    },
    removeItem: (key: string) => {
        SecureStore.deleteItemAsync(key);
    },
};

const supabaseUrl = 'https://pubiltfghmqmiirephsd.supabase.co'; // From Dashboard
const supabaseAnonKey = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InB1YmlsdGZnaG1xbWlpcmVwaHNkIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzA5Mzk5MDcsImV4cCI6MjA4NjUxNTkwN30.-WP0ACKHKxKhyKMoMtBPbQrUZRrPEQhu-opQoZAgonk'; // From Dashboard

export const supabase = createClient(supabaseUrl, supabaseAnonKey, {
    auth: {
        storage: Platform.OS === 'web' ? undefined : ExpoSecureStoreAdapter,
        autoRefreshToken: true,
        persistSession: true,
        detectSessionInUrl: false,
    },
});