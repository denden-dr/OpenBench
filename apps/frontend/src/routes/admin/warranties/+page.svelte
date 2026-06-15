<script lang="ts">
  import { Card, Button } from '$lib';
  import { warrantyService, type Warranty } from '$lib/services/warranty';
  import { onMount } from 'svelte';
  import { 
    ShieldCheck, Search, Calendar, User, 
    ArrowUpRight, Clock, Award 
  } from 'lucide-svelte';

  let warranties = $state<Warranty[]>([]);
  let isLoading = $state(true);
  let searchQuery = $state('');

  onMount(async () => {
    warranties = await warrantyService.getWarranties();
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
      // Determine real expiration status
      const isExpired = new Date(w.end_date) <= new Date();
      const realStatus = isExpired ? 'expired' : 'active';
      
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

  const getWarrantyStatus = (endDateStr: string): 'active' | 'expired' => {
    const end = new Date(endDateStr);
    return end > new Date() ? 'active' : 'expired';
  };

  const getRemainingDays = (endDateStr: string): number => {
    const end = new Date(endDateStr).getTime();
    const now = new Date().getTime();
    const diffTime = end - now;
    if (diffTime <= 0) return 0;
    return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  };

  const formatDate = (dateStr: string) => {
    const d = new Date(dateStr);
    return d.toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric' });
  };
</script>

<svelte:head>
  <title>Warranties - OpenBench Admin</title>
</svelte:head>

<div class="flex flex-col gap-6">

  <!-- Header Area -->
  <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
    <div>
      <h2 class="font-display font-extrabold text-2xl md:text-3xl uppercase tracking-tight flex items-center gap-2">
        <ShieldCheck class="w-8 h-8 text-neubrutalism-green" />
        <span>Warranty Guarantee Monitor</span>
      </h2>
      <p class="font-sans text-xs sm:text-sm text-zinc-500 font-semibold mt-1">
        Monitor active device warranties registered upon successful repair pickups.
      </p>
    </div>
  </div>

  <!-- Stats Grid -->
  <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
    <!-- Active Warranties -->
    <Card bgColor="bg-white" class="p-4 flex items-center justify-between border-4 border-neubrutalism-charcoal shadow-neubrutalism-sm">
      <div class="flex flex-col">
        <span class="font-mono text-[10px] font-bold text-zinc-500 uppercase">Active Guarantees</span>
        <span class="font-mono text-2xl font-extrabold text-neubrutalism-green">{activeWarrantiesCount}</span>
      </div>
      <Award class="w-8 h-8 text-neubrutalism-green" />
    </Card>

    <!-- Expired Warranties -->
    <Card bgColor="bg-white" class="p-4 flex items-center justify-between border-4 border-neubrutalism-charcoal shadow-neubrutalism-sm">
      <div class="flex flex-col">
        <span class="font-mono text-[10px] font-bold text-zinc-500 uppercase font-semibold">Expired Guarantees</span>
        <span class="font-mono text-2xl font-extrabold text-zinc-400">{expiredWarrantiesCount}</span>
      </div>
      <Clock class="w-8 h-8 text-zinc-300" />
    </Card>
  </div>

  <!-- Search area -->
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

  <!-- Warranties list table -->
  {#if isLoading}
    <div class="flex flex-col gap-3">
      {#each Array(4) as _}
        <div class="h-14 bg-zinc-200 border-4 border-neubrutalism-charcoal animate-pulse"></div>
      {/each}
    </div>
  {:else}
    <div class="overflow-x-auto border-4 border-neubrutalism-charcoal shadow-neubrutalism-md">
      <table class="w-full bg-white text-left font-sans text-xs border-collapse">
        <thead class="bg-zinc-100 font-display font-extrabold uppercase border-b-4 border-neubrutalism-charcoal text-[10px] sm:text-xs">
          <tr>
            <th class="p-3 border-r-2 border-neubrutalism-charcoal">Ticket Number</th>
            <th class="p-3 border-r-2 border-neubrutalism-charcoal">Customer Name</th>
            <th class="p-3 border-r-2 border-neubrutalism-charcoal">Device Info</th>
            <th class="p-3 border-r-2 border-neubrutalism-charcoal">Guarantee Period</th>
            <th class="p-3 border-r-2 border-neubrutalism-charcoal text-center">Warranty Status</th>
            <th class="p-3 text-center">Actions</th>
          </tr>
        </thead>
        <tbody class="font-mono">
          {#if filteredWarranties.length === 0}
            <tr>
              <td colspan="6" class="p-8 text-center text-zinc-500 font-sans text-sm">
                No active or expired warranties match current search filters.
              </td>
            </tr>
          {:else}
            {#each filteredWarranties as w}
              {@const status = getWarrantyStatus(w.end_date)}
              {@const daysLeft = getRemainingDays(w.end_date)}
              <tr class="border-b-2 border-neubrutalism-charcoal hover:bg-zinc-50 transition-colors">
                <!-- Ticket Number -->
                <td class="p-3 border-r-2 border-neubrutalism-charcoal font-bold">
                  {w.ticket_number}
                </td>

                <!-- Customer -->
                <td class="p-3 border-r-2 border-neubrutalism-charcoal font-sans text-zinc-700">
                  <div class="flex items-center gap-1">
                    <User class="w-3.5 h-3.5 text-zinc-400" />
                    <span class="font-bold">{w.customer_name}</span>
                  </div>
                </td>

                <!-- Device Info -->
                <td class="p-3 border-r-2 border-neubrutalism-charcoal font-sans text-zinc-800">
                  {w.device_info}
                </td>

                <!-- Period -->
                <td class="p-3 border-r-2 border-neubrutalism-charcoal font-sans text-zinc-600">
                  <div class="flex flex-col gap-0.5">
                    <span>Starts: {formatDate(w.start_date)}</span>
                    <span class="font-bold">Expires: {formatDate(w.end_date)}</span>
                  </div>
                </td>

                <!-- Status & Time Remaining -->
                <td class="p-3 border-r-2 border-neubrutalism-charcoal text-center">
                  <div class="flex flex-col items-center gap-1">
                    {#if status === 'active'}
                      <span class="bg-neubrutalism-green border border-neubrutalism-charcoal py-0.5 px-2 font-mono text-[9px] font-bold shadow-neubrutalism-sm">
                        ACTIVE
                      </span>
                      <span class="text-[8px] text-emerald-600 font-bold uppercase">{daysLeft} days remaining</span>
                    {:else}
                      <span class="bg-zinc-200 border border-zinc-400 text-zinc-500 py-0.5 px-2 font-mono text-[9px] font-bold">
                        EXPIRED
                      </span>
                      <span class="text-[8px] text-zinc-400 font-bold uppercase">Guarantee Ended</span>
                    {/if}
                  </div>
                </td>

                <!-- View Ticket Link -->
                <td class="p-3 text-center">
                  <a href="/admin/tickets/{w.ticket_id}">
                    <button 
                      class="p-1.5 border-2 border-neubrutalism-charcoal bg-white hover:bg-zinc-100 shadow-neubrutalism-sm active:translate-y-0.5 active:shadow-none cursor-pointer"
                      title="View Associated Repair Ticket"
                    >
                      <ArrowUpRight class="w-4 h-4 text-neubrutalism-charcoal" />
                    </button>
                  </a>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  {/if}

</div>
