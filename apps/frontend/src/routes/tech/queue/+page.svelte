<script lang="ts">
  import { Clock, CheckCircle2, UserPlus } from 'lucide-svelte';
  
  let activeTab = 'unassigned'; // 'unassigned' or 'mine'
</script>

<div class="max-w-7xl mx-auto px-4 py-12">
  <div class="mb-12">
    <h1 class="font-['Fira_Code'] text-4xl lg:text-5xl font-black text-slate-900 dark:text-white uppercase tracking-tight">Technician Lab</h1>
    <p class="font-['Fira_Sans'] text-lg text-slate-500 dark:text-slate-400 font-medium mt-2">Manage diagnostics and active repair workflows</p>
  </div>

  <div class="flex bg-slate-100 dark:bg-slate-900 p-1.5 rounded-2xl border border-slate-200 dark:border-slate-800 backdrop-blur-md mb-8 inline-flex">
    <button 
      class="px-8 py-2.5 font-['Fira_Sans'] font-black text-xs uppercase tracking-widest rounded-xl transition-all duration-300 {activeTab === 'unassigned' ? 'bg-white dark:bg-slate-800 text-blue-600 shadow-xl shadow-blue-900/5' : 'text-slate-500 hover:text-slate-700'}"
      on:click={() => activeTab = 'unassigned'}
    >
      <Clock class="w-4 h-4 mr-2 inline" /> Unassigned (12)
    </button>
    <button 
      class="px-8 py-2.5 font-['Fira_Sans'] font-black text-xs uppercase tracking-widest rounded-xl transition-all duration-300 {activeTab === 'mine' ? 'bg-white dark:bg-slate-800 text-blue-600 shadow-xl shadow-blue-900/5' : 'text-slate-500 hover:text-slate-700'}"
      on:click={() => activeTab = 'mine'}
    >
      <CheckCircle2 class="w-4 h-4 mr-2 inline" /> My Tickets (3)
    </button>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
    <!-- Unassigned Ticket Card Example -->
    <div class="group bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200 dark:border-slate-800 p-8 rounded-3xl shadow-sm hover:shadow-2xl hover:shadow-blue-600/10 transition-all duration-500 flex flex-col h-full border-t-8 border-t-slate-300">
      <div class="flex justify-between items-start mb-6">
        <span class="font-['Fira_Code'] font-black text-blue-600 text-sm tracking-tight">TKT-89A2</span>
        <span class="text-[10px] font-black uppercase tracking-[0.2em] text-amber-600 bg-amber-50 dark:bg-amber-900/20 px-3 py-1 rounded-full border border-amber-100 dark:border-amber-800">45m Wait</span>
      </div>
      
      <div class="mb-8 flex-grow">
        <h3 class="font-['Fira_Sans'] font-black text-2xl text-slate-900 dark:text-white mb-2 uppercase tracking-tight">Apple iPhone 13 Pro</h3>
        <p class="font-['Fira_Sans'] text-slate-500 dark:text-slate-400 font-bold leading-relaxed line-clamp-3">Screen cracked, touch unresponsive. Customer also mentioned battery drains quickly.</p>
      </div>
      
      <div class="mt-auto">
        <button class="w-full bg-blue-600 text-white font-black text-xs uppercase tracking-widest py-4 rounded-xl hover:bg-blue-700 transition-all shadow-lg shadow-blue-600/20 active:scale-95 flex items-center justify-center gap-2 cursor-pointer">
          <UserPlus class="w-4 h-4" /> Claim Ticket
        </button>
      </div>
    </div>
  </div>
</div>
