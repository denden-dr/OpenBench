<script lang="ts">
  import { warrantyService, type Warranty } from '$lib/services/warranty';
  import { onMount } from 'svelte';
  import { Search } from 'lucide-svelte';

  import WarrantiesHeader from './components/WarrantiesHeader.svelte';
  import WarrantiesTable from './components/WarrantiesTable.svelte';

  let warranties = $state<Warranty[]>([]);
  let isLoading = $state(true);
  let searchQuery = $state('');

  onMount(async () => {
    warranties = await warrantyService.getWarranties() ?? [];
    isLoading = false;
  });

  // Calculate stats reactively
  let activeWarrantiesCount = $derived(
    warranties.filter(w => new Date(w.end_date) > new Date()).length
  );
  let expiredWarrantiesCount = $derived(
    warranties.filter(w => new Date(w.end_date) <= new Date()).length
  );

  // Filter warranties reactively
  let filteredWarranties = $derived(
    warranties.filter(w => {
      if (searchQuery.trim() !== '') {
        const query = searchQuery.toLowerCase();
        const numMatch = w.ticket_number.toLowerCase().includes(query);
        const nameMatch = w.customer_name.toLowerCase().includes(query);
        const deviceMatch = w.device_info.toLowerCase().includes(query);
        return numMatch || nameMatch || deviceMatch;
      }
      return true;
    })
  );
</script>

<svelte:head>
  <title>Warranties - OpenBench Admin</title>
</svelte:head>

<div class="flex flex-col gap-6">

  <!-- Header & Stats Area -->
  <WarrantiesHeader {activeWarrantiesCount} {expiredWarrantiesCount} />

  <!-- Search Area -->
  <div class="flex flex-col md:flex-row gap-4 items-center justify-between">
    <div class="relative w-full md:max-w-md">
      <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-zinc-500">
        <Search class="w-4 h-4" />
      </div>
      <input 
        type="text" 
        placeholder="Search ticket number, customer name, device..."
        bind:value={searchQuery}
        class="w-full pl-9 pr-4 py-2 border-4 border-neubrutalism-charcoal bg-white focus:outline-none focus:bg-zinc-50 font-mono text-sm shadow-neubrutalism-sm"
      />
    </div>
  </div>

  <!-- Warranties List Table -->
  <WarrantiesTable {filteredWarranties} {isLoading} />

</div>
