<script lang="ts">
  import { Search, Plus, ShieldCheck } from "lucide-svelte";

  let {
    searchQuery = $bindable(),
    statusFilter = $bindable(),
    onAddTicket,
  }: {
    searchQuery: string;
    statusFilter: string;
    onAddTicket: () => void;
  } = $props();
</script>

<div class="flex flex-col md:flex-row justify-between items-center gap-4 mb-6">
  <div class="flex flex-col sm:flex-row items-center w-full md:w-auto gap-4">
    <!-- Search Input -->
    <div class="relative w-full sm:w-80">
      <span class="absolute inset-y-0 left-0 flex items-center pl-3.5 pointer-events-none text-slate-400">
        <Search size={18} />
      </span>
      <input
        type="text"
        bind:value={searchQuery}
        placeholder="Search by customer, model, issue..."
        class="w-full bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl py-2.5 pl-10 pr-4 text-sm text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all placeholder-slate-400"
      />
    </div>

    <!-- Status Filter -->
    <select
      bind:value={statusFilter}
      class="w-full sm:w-44 bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl py-2.5 px-3.5 text-sm text-slate-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 transition-all cursor-pointer"
    >
      <option value="all">Semua Status</option>
      <option value="service_in">Masuk</option>
      <option value="on_process">Sedang Diproses</option>
      <option value="waiting_confirmation">Menunggu Konfirmasi</option>
      <option value="cancelled">Dibatalkan</option>
      <option value="fixed">Selesai Diperbaiki</option>
      <option value="picked_up">Sudah Diambil</option>
    </select>
  </div>

  <!-- Actions -->
  <div class="flex items-center gap-2 w-full md:w-auto">
    <a
      href="/warranty"
      class="w-full md:w-auto px-6 py-2.5 border border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-300 font-bold text-sm rounded-xl hover:bg-slate-50 dark:hover:bg-slate-900 transition-all inline-flex items-center justify-center gap-2"
    >
      <ShieldCheck size={16} class="text-indigo-600" />
      Klaim Garansi
    </a>
    <button
      onclick={onAddTicket}
      class="w-full md:w-auto px-6 py-2.5 bg-blue-600 hover:bg-blue-700 text-white font-bold text-sm rounded-xl transition-all shadow-sm hover:shadow active:scale-95 inline-flex items-center justify-center gap-2"
    >
      <Plus size={16} />
      New Repair
    </button>
  </div>
</div>
