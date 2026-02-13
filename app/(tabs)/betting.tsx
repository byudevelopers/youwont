import { StyleSheet, View, Text, ScrollView, SafeAreaView, TouchableOpacity } from 'react-native';
import { IconSymbol } from '@/components/ui/icon-symbol';

export default function BettingScreen() {
  const activeBets = [
    { id: 1, title: 'Lakers vs Warriors', type: 'Spread', odds: '-110', status: 'Live', amount: '$20' },
    { id: 2, title: 'Chiefs to Win', type: 'Moneyline', odds: '+150', status: 'Pending', amount: '$50' },
    { id: 3, title: 'Over 2.5 Goals', type: 'Total', odds: '-105', status: 'Live', amount: '$15' },
  ];

  return (
    <SafeAreaView style={styles.container}>
        <View style={styles.headerContainer}>
            <Text style={styles.headerTitle}>My Bets</Text>
            <TouchableOpacity>
                <IconSymbol size={24} name="slider.horizontal.3" color="#0f172a" />
            </TouchableOpacity>
        </View>

        <ScrollView contentContainerStyle={styles.scrollContent}>
            
            {/* Filter Pills */}
            <View style={styles.filterRow}>
                <TouchableOpacity style={[styles.filterPill, styles.filterPillActive]}>
                    <Text style={styles.filterTextActive}>Active</Text>
                </TouchableOpacity>
                <TouchableOpacity style={styles.filterPill}>
                    <Text style={styles.filterText}>Settled</Text>
                </TouchableOpacity>
                <TouchableOpacity style={styles.filterPill}>
                    <Text style={styles.filterText}>Won</Text>
                </TouchableOpacity>
            </View>

            {/* Bet List */}
            {activeBets.map((bet) => (
                <View key={bet.id} style={styles.betCard}>
                    <View style={styles.cardHeader}>
                        <View style={styles.badgeContainer}>
                            <View style={[styles.statusDot, { backgroundColor: bet.status === 'Live' ? '#ef4444' : '#f59e0b' }]} />
                            <Text style={styles.statusText}>{bet.status}</Text>
                        </View>
                        <Text style={styles.oddsText}>{bet.odds}</Text>
                    </View>
                    
                    <Text style={styles.betTitle}>{bet.title}</Text>
                    <Text style={styles.betType}>{bet.type}</Text>

                    <View style={styles.divider} />

                    <View style={styles.cardFooter}>
                        <Text style={styles.wagerLabel}>Wager</Text>
                        <Text style={styles.wagerAmount}>{bet.amount}</Text>
                    </View>
                </View>
            ))}

        </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f8fafc',
  },
  headerContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 24,
    paddingVertical: 16,
    marginTop: 10,
  },
  headerTitle: {
    fontSize: 28,
    fontWeight: '800',
    color: '#0f172a',
  },
  scrollContent: {
    paddingHorizontal: 24,
    paddingBottom: 100,
  },
  filterRow: {
    flexDirection: 'row',
    marginBottom: 24,
    gap: 12,
  },
  filterPill: {
    paddingVertical: 8,
    paddingHorizontal: 20,
    borderRadius: 20,
    borderWidth: 1,
    borderColor: '#e2e8f0',
    backgroundColor: '#ffffff',
  },
  filterPillActive: {
    backgroundColor: '#0f172a',
    borderColor: '#0f172a',
  },
  filterText: {
    fontSize: 14,
    fontWeight: '600',
    color: '#64748b',
  },
  filterTextActive: {
    fontSize: 14,
    fontWeight: '600',
    color: '#ffffff',
  },
  betCard: {
    backgroundColor: '#ffffff',
    borderRadius: 16,
    padding: 20,
    marginBottom: 16,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.05,
    shadowRadius: 8,
    elevation: 2,
    borderWidth: 1,
    borderColor: '#f1f5f9',
  },
  cardHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  badgeContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#f8fafc',
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 8,
  },
  statusDot: {
    width: 6,
    height: 6,
    borderRadius: 3,
    marginRight: 6,
  },
  statusText: {
    fontSize: 12,
    fontWeight: '600',
    color: '#475569',
  },
  oddsText: {
    fontSize: 16,
    fontWeight: '700',
    color: '#7c3aed',
  },
  betTitle: {
    fontSize: 18,
    fontWeight: '700',
    color: '#0f172a',
    marginBottom: 4,
  },
  betType: {
    fontSize: 14,
    color: '#64748b',
    fontWeight: '500',
  },
  divider: {
    height: 1,
    backgroundColor: '#f1f5f9',
    marginVertical: 16,
  },
  cardFooter: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  wagerLabel: {
    fontSize: 14,
    color: '#94a3b8',
    fontWeight: '500',
  },
  wagerAmount: {
    fontSize: 16,
    fontWeight: '700',
    color: '#0f172a',
  },
});