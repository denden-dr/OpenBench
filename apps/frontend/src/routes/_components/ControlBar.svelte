<script lang="ts">
  import { Search, Plus, ShieldCheck } from "lucide-svelte";

  let {
    searchQuery = $bindable(),
    onAddTicket,
  }: {
    searchQuery: string;
    onAddTicket: () => void;
  } = $props();
</script>

<div class="flex flex-col lg:flex-row justify-between items-stretch lg:items-center gap-4 mb-6">
  <!-- Left Side: Search -->
  <div class="flex items-stretch md:items-center gap-4 flex-1 min-w-0">
    <!-- Search Input -->
    <div class="relative w-full md:w-72 lg:w-96 shrink-0">
      <span class="absolute inset-y-0 left-0 flex items-center pl-3.5 pointer-events-none text-slate-400">
        <Search size={18} />
      </span>
      <input
        type="text"
        bind:value={searchQuery}
        placeholder="Cari pelanggan, model, kendala..."
        class="w-full bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl py-2.5 pl-10 pr-4 text-sm text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all placeholder-slate-400"
      />
    </div>
  </div>

  <!-- Right Side: Actions -->
  <div class="flex items-center gap-2 w-full md:w-auto shrink-0 justify-end">
    <a
      href="/warranty"
      class="w-full md:w-auto px-5 py-2.5 border border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 font-bold text-sm rounded-xl hover:bg-slate-50 dark:hover:bg-slate-900 transition-all inline-flex items-center justify-center gap-2"
    >
      <ShieldCheck size={16} class="text-indigo-600" />
      Klaim Garansi
    </a>
    <button
      onclick={onAddTicket}
      class="w-full md:w-auto px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-white font-bold text-sm rounded-xl transition-all shadow-sm hover:shadow active:scale-95 inline-flex items-center justify-center gap-2"
    >
      <Plus size={16} />
      New Repair
    </button>
  </div>
</div>
