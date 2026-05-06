<script lang="ts">
  import { Save, AlertTriangle, CheckCircle, Package, Camera, Wrench } from 'lucide-svelte';
  
  let currentStatus = 'diagnosing'; // Mock state for UI dev
</script>

<div class="max-w-7xl mx-auto px-4 py-12">
  <!-- Header -->
  <div class="bg-slate-900 text-white rounded-[2rem] p-8 lg:p-12 flex flex-wrap justify-between items-center gap-6 relative overflow-hidden shadow-2xl border border-slate-800">
    <!-- Decor -->
    <div class="absolute top-0 right-0 w-64 h-64 bg-blue-600/10 rounded-full translate-x-1/3 -translate-y-1/3 blur-3xl"></div>
    
    <div class="relative z-10">
      <div class="flex items-center gap-4 mb-2">
        <h1 class="font-['Fira_Code'] text-2xl font-black tracking-tight">TKT-89A2</h1>
        <span class="bg-amber-500/10 text-amber-500 text-[10px] font-black px-3 py-1.5 rounded-full uppercase tracking-[0.2em] border border-amber-500/20">Diagnosing</span>
      </div>
      <p class="font-['Fira_Sans'] text-slate-400 text-lg font-bold">Apple iPhone 13 Pro • John Doe</p>
    </div>
    
    <div class="flex gap-4 relative z-10">
      <button class="bg-white/5 hover:bg-white/10 px-6 py-3 rounded-xl text-xs font-black uppercase tracking-widest transition-all backdrop-blur-md border border-white/10">
        Mark Unrepairable
      </button>
      <button class="bg-blue-600 hover:bg-blue-700 px-8 py-3 rounded-xl text-xs font-black uppercase tracking-widest transition-all shadow-lg shadow-blue-600/20 flex items-center gap-2">
        <Save class="w-4 h-4" /> Submit Diagnosis
      </button>
    </div>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-10 mt-10">
    <!-- Left Column: Customer & Device Info -->
    <div class="lg:col-span-1 space-y-10">
      <div class="bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200 dark:border-slate-800 rounded-[2rem] p-8 shadow-xl shadow-blue-900/5">
        <h2 class="font-['Fira_Code'] text-xs font-black text-slate-400 uppercase tracking-[0.3em] mb-8">Issue Details</h2>
        <p class="font-['Fira_Sans'] text-slate-900 dark:text-white bg-slate-50 dark:bg-slate-950 p-6 rounded-2xl mb-8 font-bold border border-slate-100 dark:border-slate-800">
          Screen cracked, touch unresponsive. Customer also mentioned battery drains quickly.
        </p>
        
        <div class="space-y-6">
          <div>
            <h3 class="text-[10px] font-black text-slate-400 mb-2 uppercase tracking-widest">Access Info</h3>
            <p class="text-blue-600 font-black bg-blue-50 dark:bg-blue-900/20 p-3 rounded-xl text-sm border border-blue-100 dark:border-blue-800">Passcode: 123456</p>
          </div>
          
          <div>
            <h3 class="text-[10px] font-black text-slate-400 mb-2 uppercase tracking-widest">Accessories</h3>
            <ul class="text-sm text-slate-600 dark:text-slate-400 space-y-2 font-bold">
              <li class="flex items-center gap-3"><CheckCircle class="w-4 h-4 text-emerald-500" /> SIM Card</li>
              <li class="flex items-center gap-3"><CheckCircle class="w-4 h-4 text-emerald-500" /> Clear Case</li>
            </ul>
          </div>
        </div>
      </div>

      <div class="bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200 dark:border-slate-800 rounded-[2rem] p-8 shadow-xl shadow-blue-900/5">
        <h2 class="font-['Fira_Code'] text-xs font-black text-slate-400 uppercase tracking-[0.3em] mb-8 flex items-center justify-between">
          Photos
          <button class="text-blue-600 bg-blue-50 dark:bg-blue-900/20 p-2 rounded-xl hover:bg-blue-100 transition-all cursor-pointer border border-blue-100 dark:border-blue-800">
            <Camera class="w-4 h-4" />
          </button>
        </h2>
        <div class="grid grid-cols-2 gap-4">
          <div class="aspect-square bg-slate-100 dark:bg-slate-950 rounded-2xl border border-slate-200 dark:border-slate-800 flex items-center justify-center text-[10px] font-black uppercase text-slate-400 tracking-widest">Before 1</div>
          <div class="aspect-square bg-slate-50 dark:bg-slate-950 rounded-2xl border-2 border-dashed border-slate-200 dark:border-slate-800 flex items-center justify-center text-blue-600 cursor-pointer hover:bg-blue-50 transition-all font-black text-xs uppercase tracking-widest">+ Add</div>
        </div>
      </div>
    </div>

    <!-- Right Column: Workspace -->
    <div class="lg:col-span-2">
      <div class="bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200 dark:border-slate-800 rounded-[2rem] p-10 shadow-xl shadow-blue-900/5 h-full">
        <h2 class="font-['Fira_Code'] text-2xl font-black text-slate-900 dark:text-white mb-10 flex items-center gap-4 uppercase tracking-tight">
          <Wrench class="w-8 h-8 text-blue-600" /> Diagnosis Workspace
        </h2>
        
        <form class="space-y-8" on:submit|preventDefault>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
            <div class="space-y-6">
              <div>
                <label class="block text-[10px] font-black text-slate-400 mb-2 uppercase tracking-[0.2em]">Technical Notes (Internal)</label>
                <textarea class="w-full bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl p-4 focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 outline-none text-sm font-bold h-32 transition-all dark:text-white" placeholder="Flex cable damaged..."></textarea>
              </div>
              
              <div>
                <label class="block text-[10px] font-black text-slate-400 mb-2 uppercase tracking-[0.2em]">Customer Findings</label>
                <textarea class="w-full bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl p-4 focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 outline-none text-sm font-bold h-32 transition-all dark:text-white" placeholder="What the customer will see..."></textarea>
              </div>
            </div>

            <div class="space-y-6">
              <div>
                <label class="block text-[10px] font-black text-slate-400 mb-2 uppercase tracking-[0.2em]">Diagnosis Fee ($)</label>
                <input type="number" class="w-full bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl p-4 focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 outline-none font-black text-blue-600 transition-all" value="25.00" />
              </div>
              <div>
                <label class="block text-[10px] font-black text-slate-400 mb-2 uppercase tracking-[0.2em]">Estimated Labor ($)</label>
                <input type="number" class="w-full bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl p-4 focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 outline-none font-black text-blue-600 transition-all" placeholder="0.00" />
              </div>
              <div>
                <label class="block text-[10px] font-black text-slate-400 mb-2 uppercase tracking-[0.2em]">Est. Ready Date</label>
                <input type="datetime-local" class="w-full bg-slate-50 dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-2xl p-4 focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 outline-none font-bold text-slate-600 dark:text-slate-400 transition-all" />
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  </div>
</div>
