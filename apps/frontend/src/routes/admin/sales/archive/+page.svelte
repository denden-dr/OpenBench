<script lang="ts">
  import { saleService, type Sale } from '$lib/services/sales';
  import { onMount } from 'svelte';
  import { toastService } from '$lib/services/toast.svelte';

  import ArchiveHeader from './components/ArchiveHeader.svelte';
  import ArchiveStats from './components/ArchiveStats.svelte';
  import ArchiveSearch from './components/ArchiveSearch.svelte';
  import ArchiveTable from './components/ArchiveTable.svelte';
  import ArchiveReceiptModal from './components/ArchiveReceiptModal.svelte';

  let sales = $state<Sale[]>([]);
  let isLoading = $state(true);
  let searchQuery = $state('');

  // Selected sale for receipt reprint modal
  let selectedSale = $state<Sale | null>(null);

  onMount(async () => {
    await loadSales();
  });

  async function loadSales() {
    try {
      sales = await saleService.getSales();
    } catch (err: any) {
      console.error('Error loading sales archive:', err);
      toastService.error('Failed to load sales: ' + err.message);
      sales = [];
    } finally {
      isLoading = false;
    }
  }

  // Filter sales reactively
  let filteredSales = $derived(
    sales.filter(s => {
      if (searchQuery.trim() !== '') {
        const query = searchQuery.toLowerCase();
        const invMatch = s.invoice_number.toLowerCase().includes(query);
        const itemMatch = s.items.some(item => item.name.toLowerCase().includes(query));
        return invMatch || itemMatch;
      }
      return true;
    })
  );

  // Total gross revenue calculation
  let totalRevenue = $derived(sales.reduce((sum, s) => sum + s.total, 0));

  const formatDate = (dateStr: string) => {
    const d = new Date(dateStr);
    return d.toLocaleDateString('en-US', { 
      day: 'numeric', 
      month: 'short', 
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };
</script>

<svelte:head>
  <title>Sales Archive - OpenBench Admin</title>
</svelte:head>

<div class="flex flex-col gap-6">

  <ArchiveHeader />

  <ArchiveStats
    {totalRevenue}
    salesCount={sales.length}
  />

  <ArchiveSearch
    bind:searchQuery
  />

  <ArchiveTable
    {filteredSales}
    {isLoading}
    onSelectSale={(s) => { selectedSale = s; }}
    {formatDate}
  />

  <!-- Receipt Reprint Modal Overlay -->
  {#if selectedSale !== null}
    <ArchiveReceiptModal
      {selectedSale}
      onClose={() => { selectedSale = null; }}
    />
  {/if}

</div>
