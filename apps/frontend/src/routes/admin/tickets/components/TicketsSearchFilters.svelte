<script lang="ts">
  import { Search, Filter } from 'lucide-svelte';

  interface Props {
    searchQuery: string;
    statusFilter: string;
    activeTab: 'active' | 'archive';
  }

  let {
    searchQuery = $bindable(),
    statusFilter = $bindable(),
    activeTab
  }: Props = $props();
</script>

<div class="flex flex-col md:flex-row gap-4 items-center justify-between">
  <!-- Search Box -->
  <div class="relative w-full md:max-w-md">
    <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-zinc-500">
      <Search class="w-4 h-4" />
    </div>
    <input 
      type="text" 
      aria-label="Search tickets"
      placeholder="Search ticket number, customer name, device..."
      bind:value={searchQuery}
      class="w-full pl-9 pr-4 py-2 border-4 border-neubrutalism-charcoal bg-white focus:outline-none focus:bg-zinc-50 font-mono text-sm shadow-neubrutalism-sm"
    />
  </div>

  <!-- Status Dropdown Filter -->
  <div class="flex items-center gap-2 w-full md:w-auto">
    <Filter class="w-4 h-4 text-zinc-600 shrink-0" />
    <span class="font-mono text-xs font-bold shrink-0">STATUS:</span>
    <select 
      bind:value={statusFilter}
      class="w-full md:w-48 bg-white border-4 border-neubrutalism-charcoal px-3 py-1.5 font-mono text-xs shadow-neubrutalism-sm focus:outline-none"
    >
      <option value="all">ALL STATUS</option>
      {#if activeTab === 'active'}
        <option value="received">RECEIVED</option>
        <option value="in_repair">IN REPAIR</option>
        <option value="ready_for_pickup">READY FOR PICKUP</option>
      {:else}
        <option value="completed">COMPLETED</option>
        <option value="cancelled">CANCELLED</option>
      {/if}
    </select>
  </div>
</div>
