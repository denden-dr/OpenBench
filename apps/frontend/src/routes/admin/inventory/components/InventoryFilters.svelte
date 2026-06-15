<script lang="ts">
  import { Search, Filter } from 'lucide-svelte';

  let {
    searchQuery = $bindable(''),
    categoryFilter = $bindable('all'),
    stockFilter = $bindable('all')
  }: {
    searchQuery: string;
    categoryFilter: string;
    stockFilter: string;
  } = $props();
</script>

<div class="flex flex-col md:flex-row gap-4 items-center justify-between">
  <!-- Search -->
  <div class="relative w-full md:max-w-md">
    <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-zinc-500">
      <Search class="w-4 h-4" />
    </div>
    <input 
      type="text" 
      placeholder="Search item name..."
      bind:value={searchQuery}
      class="w-full pl-9 pr-4 py-2 border-4 border-neubrutalism-charcoal bg-white focus:outline-none focus:bg-zinc-50 font-mono text-sm shadow-neubrutalism-sm"
    />
  </div>

  <!-- Dropdown Filters -->
  <div class="flex flex-col sm:flex-row items-center gap-3 w-full md:w-auto">
    <div class="flex items-center gap-2 w-full sm:w-auto">
      <Filter class="w-4 h-4 text-zinc-600 shrink-0" />
      <span class="font-mono text-xs font-bold shrink-0">CATEGORY:</span>
      <select 
        bind:value={categoryFilter}
        class="w-full sm:w-36 bg-white border-4 border-neubrutalism-charcoal px-3 py-1.5 font-mono text-xs shadow-neubrutalism-sm focus:outline-none"
      >
        <option value="all">ALL</option>
        <option value="retail">RETAIL</option>
        <option value="spare_part">SPARE PARTS</option>
      </select>
    </div>

    <div class="flex items-center gap-2 w-full sm:w-auto">
      <span class="font-mono text-xs font-bold shrink-0">STOCK:</span>
      <select 
        bind:value={stockFilter}
        class="w-full sm:w-36 bg-white border-4 border-neubrutalism-charcoal px-3 py-1.5 font-mono text-xs shadow-neubrutalism-sm focus:outline-none"
      >
        <option value="all">ALL QUANTITY</option>
        <option value="low_stock">LOW STOCK WARNINGS</option>
      </select>
    </div>
  </div>
</div>
