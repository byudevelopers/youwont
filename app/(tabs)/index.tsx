import { StyleSheet, View, Text, TouchableOpacity, ScrollView, SafeAreaView, Image } from 'react-native';
import { supabase } from '../../lib/supabase';
import { IconSymbol } from '@/components/ui/icon-symbol';

export default function HomeScreen() {

  const handleSignOut = async () => {
    const { error } = await supabase.auth.signOut();
    if (error) console.error(error.message);
  };

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        
        {/* Header Section */}
        <View style={styles.header}>
            <View>
                <Text style={styles.greeting}>Hello, User</Text>
                <Text style={styles.subGreeting}>Welcome back to YouWont</Text>
            </View>
            <TouchableOpacity style={styles.profileButton}>
                 <IconSymbol size={24} name="person.crop.circle" color="#64748b" />
            </TouchableOpacity>
        </View>

        {/* Balance Card */}
        <View style={styles.card}>
            <Text style={styles.cardLabel}>Current Balance</Text>
            <Text style={styles.cardAmount}>$1,240.50</Text>
            <View style={styles.cardRow}>
                <View style={styles.badge}>
                    <Text style={styles.badgeText}>+ 12% this week</Text>
                </View>
            </View>
        </View>

        {/* Action Buttons */}
        <View style={styles.actionRow}>
            <TouchableOpacity style={styles.actionButton}>
                <View style={styles.iconCircle}>
                    <IconSymbol size={24} name="plus" color="#7c3aed" />
                </View>
                <Text style={styles.actionText}>New Bet</Text>
            </TouchableOpacity>

            <TouchableOpacity style={styles.actionButton}>
                <View style={[styles.iconCircle, { backgroundColor: '#eff6ff' }]}>
                    <IconSymbol size={24} name="arrow.up.right" color="#7c3aed" />
                </View>
                <Text style={styles.actionText}>Deposit</Text>
            </TouchableOpacity>
        </View>

        {/* Recent Activity Section */}
        <Text style={styles.sectionTitle}>Recent Activity</Text>
        
        {[1, 2, 3].map((item, index) => (
            <View key={index} style={styles.activityItem}>
                <View style={styles.activityIcon}>
                     <IconSymbol size={20} name="flag.fill" color="#ffffff" />
                </View>
                <View style={styles.activityInfo}>
                    <Text style={styles.activityTitle}>Won bet against @alex</Text>
                    <Text style={styles.activityDate}>2 hours ago</Text>
                </View>
                <Text style={styles.activityAmount}>+$50.00</Text>
            </View>
        ))}

        {/* Sign Out */}
        <TouchableOpacity style={styles.signOutButton} onPress={handleSignOut}>
            <Text style={styles.signOutText}>Sign Out</Text>
        </TouchableOpacity>

      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f8fafc', // Slate 50
  },
  scrollContent: {
    padding: 24,
    paddingBottom: 100, // Space for tab bar
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 24,
    marginTop: 10,
  },
  greeting: {
    fontSize: 24,
    fontWeight: '700',
    color: '#0f172a',
  },
  subGreeting: {
    fontSize: 14,
    color: '#64748b',
  },
  profileButton: {
    padding: 8,
    backgroundColor: '#ffffff',
    borderRadius: 20,
    shadowColor: '#000',
    shadowOpacity: 0.05,
    shadowRadius: 5,
  },
  card: {
    backgroundColor: '#7c3aed', // Startup Blue
    borderRadius: 20,
    padding: 24,
    marginBottom: 24,
    shadowColor: '#7c3aed',
    shadowOffset: { width: 0, height: 8 },
    shadowOpacity: 0.25,
    shadowRadius: 12,
    elevation: 8,
  },
  cardLabel: {
    color: '#bfdbfe', // Blue 200
    fontSize: 14,
    fontWeight: '500',
    marginBottom: 8,
  },
  cardAmount: {
    color: '#ffffff',
    fontSize: 36,
    fontWeight: '700',
    marginBottom: 16,
  },
  cardRow: {
    flexDirection: 'row',
  },
  badge: {
    backgroundColor: 'rgba(255,255,255,0.2)',
    paddingHorizontal: 10,
    paddingVertical: 4,
    borderRadius: 12,
  },
  badgeText: {
    color: '#ffffff',
    fontSize: 12,
    fontWeight: '600',
  },
  actionRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 32,
  },
  actionButton: {
    flex: 1,
    backgroundColor: '#ffffff',
    borderRadius: 16,
    padding: 16,
    marginHorizontal: 6,
    alignItems: 'center',
    shadowColor: '#000',
    shadowOpacity: 0.05,
    shadowRadius: 10,
    elevation: 2,
  },
  iconCircle: {
    width: 48,
    height: 48,
    borderRadius: 24,
    backgroundColor: '#eff6ff',
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 12,
  },
  actionText: {
    fontSize: 14,
    fontWeight: '600',
    color: '#0f172a',
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '700',
    color: '#0f172a',
    marginBottom: 16,
  },
  activityItem: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#ffffff',
    padding: 16,
    borderRadius: 16,
    marginBottom: 12,
    borderWidth: 1,
    borderColor: '#f1f5f9',
  },
  activityIcon: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: '#7c3aed', // Emerald 500
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 16,
  },
  activityInfo: {
    flex: 1,
  },
  activityTitle: {
    fontSize: 15,
    fontWeight: '600',
    color: '#0f172a',
  },
  activityDate: {
    fontSize: 13,
    color: '#64748b',
    marginTop: 2,
  },
  activityAmount: {
    fontSize: 16,
    fontWeight: '700',
    color: '#10b981',
  },
  signOutButton: {
    marginTop: 24,
    paddingVertical: 16,
    borderWidth: 1,
    borderColor: '#e2e8f0',
    borderRadius: 12,
    alignItems: 'center',
  },
  signOutText: {
    color: '#64748b',
    fontWeight: '600',
    fontSize: 16,
  },
});