<script lang="ts">
  import { DollarSign, Smartphone, CheckCircle, Clock, AlertCircle } from 'lucide-svelte';
</script>

<div class="max-w-7xl mx-auto px-4 py-12">
  <div class="mb-12">
    <h1 class="font-['Fira_Code'] text-4xl lg:text-5xl font-black text-slate-900 dark:text-white uppercase tracking-tight">System Overview</h1>
    <p class="font-['Fira_Sans'] text-lg text-slate-500 dark:text-slate-400 font-medium mt-2">OpenBench Lab Performance & Metrics</p>
  </div>

  <!-- KPI Cards -->
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 mb-12">
    <div class="bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200 dark:border-slate-800 p-8 rounded-3xl shadow-sm hover:shadow-2xl hover:shadow-blue-600/10 transition-all duration-500 border-t-8 border-t-emerald-500">
      <div class="flex justify-between items-start">
        <div>
          <p class="font-['Fira_Code'] text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] mb-2">Revenue (Today)</p>
          <h2 class="font-['Fira_Sans'] text-3xl font-black text-slate-900 dark:text-white uppercase tracking-tight">$450.00</h2>
        </div>
        <div class="bg-emerald-50 dark:bg-emerald-900/20 text-emerald-500 p-3 rounded-2xl border border-emerald-100 dark:border-emerald-800">
          <DollarSign class="w-6 h-6" />
        </div>
      </div>
    </div>
    
    <div class="bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200 dark:border-slate-800 p-8 rounded-3xl shadow-sm hover:shadow-2xl hover:shadow-blue-600/10 transition-all duration-500 border-t-8 border-t-blue-600">
      <div class="flex justify-between items-start">
        <div>
          <p class="font-['Fira_Code'] text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] mb-2">Active Tickets</p>
          <h2 class="font-['Fira_Sans'] text-3xl font-black text-slate-900 dark:text-white uppercase tracking-tight">15</h2>
        </div>
        <div class="bg-blue-50 dark:bg-blue-900/20 text-blue-600 p-3 rounded-2xl border border-blue-100 dark:border-blue-800">
          <Clock class="w-6 h-6" />
        </div>
      </div>
    </div>
    
    <div class="bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200 dark:border-slate-800 p-8 rounded-3xl shadow-sm hover:shadow-2xl hover:shadow-blue-600/10 transition-all duration-500 border-t-8 border-t-indigo-500">
      <div class="flex justify-between items-start">
        <div>
          <p class="font-['Fira_Code'] text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] mb-2">Completed</p>
          <h2 class="font-['Fira_Sans'] text-3xl font-black text-slate-900 dark:text-white uppercase tracking-tight">8</h2>
        </div>
        <div class="bg-indigo-50 dark:bg-indigo-900/20 text-indigo-500 p-3 rounded-2xl border border-indigo-100 dark:border-indigo-800">
          <CheckCircle class="w-6 h-6" />
        </div>
      </div>
    </div>

    <div class="bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200 dark:border-slate-800 p-8 rounded-3xl shadow-sm hover:shadow-2xl hover:shadow-blue-600/10 transition-all duration-500 border-t-8 border-t-amber-500">
      <div class="flex justify-between items-start">
        <div>
          <p class="font-['Fira_Code'] text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] mb-2">Avg Time</p>
          <h2 class="font-['Fira_Sans'] text-3xl font-black text-slate-900 dark:text-white uppercase tracking-tight">1.2d</h2>
        </div>
        <div class="bg-amber-50 dark:bg-amber-900/20 text-amber-500 p-3 rounded-2xl border border-amber-100 dark:border-amber-800">
          <Smartphone class="w-6 h-6" />
        </div>
      </div>
    </div>
  </div>

  <!-- Alerts Panel -->
  <div class="bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-red-200 dark:border-red-900/30 rounded-[2rem] overflow-hidden shadow-xl shadow-red-900/5">
    <div class="bg-red-50/50 dark:bg-red-900/10 p-6 border-b border-red-200 dark:border-red-900/30 flex items-center gap-4 text-red-800 dark:text-red-400">
      <AlertCircle class="w-6 h-6" />
      <h2 class="font-['Fira_Code'] font-black text-lg uppercase tracking-tight">System Alerts</h2>
    </div>
    <div class="divide-y divide-red-100 dark:divide-red-900/20">
      <div class="flex justify-between items-center p-6 hover:bg-red-50/30 dark:hover:bg-red-900/5 transition-all">
        <div>
          <span class="font-black text-slate-900 dark:text-white uppercase text-xs tracking-widest mr-3">Low Stock:</span>
          <span class="text-slate-500 dark:text-slate-400 font-bold">iPhone 13 Pro OLED Display</span>
        </div>
        <span class="bg-red-100 dark:bg-red-900/40 text-red-700 dark:text-red-400 px-3 py-1 rounded-full text-[10px] font-black tracking-widest uppercase">2 left</span>
      </div>
    </div>
  </div>
</div>
