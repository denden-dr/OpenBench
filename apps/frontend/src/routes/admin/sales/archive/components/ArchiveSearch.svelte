<script lang="ts">
  import { Search } from 'lucide-svelte';

  let {
    searchQuery = $bindable('')
  }: {
    searchQuery: string;
  } = $props();
</script>

<div class="flex flex-col md:flex-row gap-4 items-center justify-between">
  <div class="relative w-full md:max-w-md">
    <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-zinc-500">
      <Search class="w-4 h-4" />
    </div>
    <input 
      type="text" 
      placeholder="Search invoice number or item name..."
      bind:value={searchQuery}
      class="w-full pl-9 pr-4 py-2 border-4 border-neubrutalism-charcoal bg-white focus:outline-none focus:bg-zinc-50 font-mono text-sm shadow-neubrutalism-sm"
    />
  </div>
</div>
