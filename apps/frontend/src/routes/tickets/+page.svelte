<script lang="ts">
  import { Smartphone, ChevronRight, CheckCircle2, Clock, ArrowRight, Wrench, Search, Plus } from 'lucide-svelte';
  
  let activeTab = $state('active'); // 'active' or 'completed'
  
  const repairs = [
    { id: 'TKT-89A2', device: 'iPhone 13 Pro', status: 'diagnosing', issue: 'Screen Replacement', date: 'May 05, 2026' },
    { id: 'TKT-89A1', device: 'Galaxy S24 Ultra', status: 'ready', issue: 'Battery Service', date: 'May 03, 2026' }
  ];
</script>

<div class="max-w-7xl mx-auto px-6 py-12">
  <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-8 mb-16">
    <div>
      <h1 class="font-['Fira_Code'] text-4xl lg:text-5xl font-black text-slate-900 dark:text-white uppercase tracking-tight">Your Repairs</h1>
      <p class="font-['Fira_Sans'] text-lg text-slate-500 dark:text-slate-400 font-medium mt-2">Track and manage your device service history</p>
    </div>
    
    <a href="/book" class="bg-blue-600 hover:bg-blue-700 text-white font-['Fira_Sans'] font-black text-xs uppercase tracking-widest py-4 px-8 rounded-2xl transition-all shadow-lg shadow-blue-600/20 cursor-pointer flex items-center gap-3 active:scale-95">
      <Plus size={18} /> Book New Repair
    </a>
  </div>

  <!-- Tabs -->
  <div class="flex items-center gap-1 bg-slate-100 dark:bg-slate-900/50 p-1.5 rounded-2xl mb-12 w-fit border border-slate-200 dark:border-slate-800">
    <button 
      onclick={() => activeTab = 'active'}
      class="px-8 py-3 rounded-xl text-xs font-black uppercase tracking-widest transition-all
        {activeTab === 'active' 
          ? 'bg-white dark:bg-slate-800 text-blue-600 shadow-sm border border-slate-200 dark:border-slate-700' 
          : 'text-slate-400 hover:text-slate-600 dark:hover:text-slate-300'}"
    >
      Active Repairs
    </button>
    <button 
      onclick={() => activeTab = 'completed'}
      class="px-8 py-3 rounded-xl text-xs font-black uppercase tracking-widest transition-all
        {activeTab === 'completed' 
          ? 'bg-white dark:bg-slate-800 text-blue-600 shadow-sm border border-slate-200 dark:border-slate-700' 
          : 'text-slate-400 hover:text-slate-600 dark:hover:text-slate-300'}"
    >
      Service History
    </button>
  </div>

  {#if activeTab === 'active'}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
      {#each repairs as repair}
        <a 
          href="/tickets/{repair.id}"
          class="group bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200 dark:border-slate-800 rounded-[2.5rem] p-8 transition-all hover:border-blue-600/30 hover:shadow-xl hover:shadow-blue-900/5 relative overflow-hidden flex flex-col h-full"
        >
          <!-- Status Badge -->
          <div class="flex justify-between items-start mb-8">
            <div class="w-12 h-12 rounded-2xl bg-blue-50 dark:bg-blue-900/20 text-blue-600 flex items-center justify-center group-hover:scale-110 transition-transform">
              <Smartphone size={24} />
            </div>
            <span class="text-[10px] font-black uppercase tracking-widest px-3 py-1.5 rounded-full
              {repair.status === 'diagnosing' ? 'bg-amber-50 dark:bg-amber-900/20 text-amber-600 border border-amber-500/20' : ''}
              {repair.status === 'ready' ? 'bg-emerald-50 dark:bg-emerald-900/20 text-emerald-600 border border-emerald-500/20' : ''}
            ">
              {repair.status}
            </span>
          </div>

          <div class="flex-1">
            <h3 class="font-['Fira_Sans'] text-2xl font-black text-slate-900 dark:text-white mb-2 group-hover:text-blue-600 transition-colors uppercase tracking-tight">{repair.device}</h3>
            <p class="font-['Fira_Code'] text-xs font-black text-slate-400 uppercase tracking-widest mb-6">{repair.id}</p>
            
            <div class="space-y-4 pt-6 border-t border-slate-50 dark:border-slate-800">
              <div class="flex items-center gap-3 text-slate-500 dark:text-slate-400">
                <Wrench size={14} class="text-blue-600" />
                <span class="text-xs font-bold uppercase tracking-wide">{repair.issue}</span>
              </div>
              <div class="flex items-center gap-3 text-slate-500 dark:text-slate-400">
                <Clock size={14} class="text-blue-600" />
                <span class="text-xs font-bold uppercase tracking-wide">Intake: {repair.date}</span>
              </div>
            </div>
          </div>

          <div class="mt-8 flex items-center justify-between group-hover:translate-x-2 transition-transform">
            <span class="font-['Fira_Sans'] text-xs font-black uppercase tracking-[0.2em] text-blue-600">Track Detail</span>
            <ArrowRight size={18} class="text-blue-600" />
          </div>
        </a>
      {/each}
    </div>
  {:else}
    <div class="bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200 dark:border-slate-800 rounded-[3rem] py-24 text-center shadow-xl shadow-blue-900/5">
      <div class="w-20 h-20 bg-slate-50 dark:bg-slate-800 rounded-3xl flex items-center justify-center text-slate-300 dark:text-slate-600 mx-auto mb-8 border border-slate-100 dark:border-slate-700">
        <Search size={40} />
      </div>
      <h3 class="font-['Fira_Code'] text-xl font-black text-slate-900 dark:text-white uppercase tracking-tight mb-4">No completed repairs</h3>
      <p class="font-['Fira_Sans'] text-slate-500 dark:text-slate-400 max-w-sm mx-auto font-bold uppercase tracking-widest text-[10px]">Your service history will appear here once repairs are finalized.</p>
    </div>
  {/if}
</div>
