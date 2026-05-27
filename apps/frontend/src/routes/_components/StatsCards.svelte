<script lang="ts">
  import { Wrench, CheckCircle2, AlertTriangle } from "lucide-svelte";
  import { formatCurrency } from "$lib/utils/format";

  let {
    waitingCount,
    onProcessCount,
    fixedCount,
    todaysRevenue,
  }: {
    waitingCount: number;
    onProcessCount: number;
    fixedCount: number;
    todaysRevenue: number;
  } = $props();
</script>

<div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
  <div class="flex-1 lg:flex-none bg-white dark:bg-slate-950 p-6 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm relative overflow-hidden transition-all hover:shadow-md h-32 flex flex-col justify-between">
    <div class="flex justify-between items-start">
      <div>
        <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Menunggu Konfirmasi</p>
      </div>
      <div class="p-3 bg-rose-50 dark:bg-rose-900/20 text-rose-600 dark:text-rose-400 rounded-xl">
        <AlertTriangle size={20} />
      </div>
    </div>
    <h3 class="text-3xl font-bold text-slate-900 dark:text-white mb-1">{waitingCount}</h3>
    <div class="h-1 bg-rose-500 w-full absolute bottom-0 left-0 rounded-b-2xl"></div>
  </div>

  <div class="flex-1 lg:flex-none bg-white dark:bg-slate-950 p-6 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm relative overflow-hidden transition-all hover:shadow-md h-32 flex flex-col justify-between">
    <div class="flex justify-between items-start">
      <div>
        <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Sedang Diproses</p>
      </div>
      <div class="p-3 bg-amber-50 dark:bg-amber-900/20 text-amber-600 dark:text-amber-400 rounded-xl">
        <Wrench size={20} />
      </div>
    </div>
    <h3 class="text-3xl font-bold text-slate-900 dark:text-white mb-1">{onProcessCount}</h3>
    <div class="h-1 bg-amber-500 w-full absolute bottom-0 left-0 rounded-b-2xl"></div>
  </div>

  <div class="flex-1 lg:flex-none bg-white dark:bg-slate-950 p-6 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm relative overflow-hidden transition-all hover:shadow-md h-32 flex flex-col justify-between">
    <div class="flex justify-between items-start">
      <div>
        <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Siap Diambil</p>
      </div>
      <div class="p-3 bg-emerald-50 dark:bg-emerald-900/20 text-emerald-600 dark:text-emerald-400 rounded-xl">
        <CheckCircle2 size={20} />
      </div>
    </div>
    <h3 class="text-3xl font-bold text-slate-900 dark:text-white mb-1">{fixedCount}</h3>
    <div class="h-1 bg-emerald-500 w-full absolute bottom-0 left-0 rounded-b-2xl"></div>
  </div>

  <div class="flex-1 lg:flex-none bg-white dark:bg-slate-950 p-6 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm relative overflow-hidden transition-all hover:shadow-md h-32 flex flex-col justify-between">
    <div class="flex justify-between items-start">
      <div>
        <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Pendapatan Hari Ini</p>
      </div>
      <div class="w-11 h-11 flex items-center justify-center bg-indigo-50 dark:bg-indigo-900/20 text-indigo-600 dark:text-indigo-400 rounded-xl">
        <span class="font-black text-sm tracking-tighter select-none">RP</span>
      </div>
    </div>
    <h3 class="text-2xl font-bold text-slate-900 dark:text-white mb-1">{formatCurrency(todaysRevenue)}</h3>
    <div class="h-1 bg-indigo-500 w-full absolute bottom-0 left-0 rounded-b-2xl"></div>
  </div>
</div>
