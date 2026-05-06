<script lang="ts">
  import { Search, Filter } from 'lucide-svelte';
</script>

<div class="max-w-7xl mx-auto px-4 py-12">
  <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-12 gap-8">
    <div>
      <h1 class="font-['Fira_Code'] text-4xl lg:text-5xl font-black text-slate-900 dark:text-white uppercase tracking-tight">Ticket Manager</h1>
      <p class="font-['Fira_Sans'] text-lg text-slate-500 dark:text-slate-400 font-medium mt-2">Global oversight of all active and historical repairs</p>
    </div>
  </div>

  <!-- Controls -->
  <div class="flex flex-col sm:flex-row gap-6 mb-10">
    <div class="relative flex-grow group">
      <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-slate-400 group-focus-within:text-blue-600 transition-colors">
        <Search class="h-5 w-5" />
      </div>
      <input type="text" class="block w-full pl-12 pr-4 py-4 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200 dark:border-slate-800 rounded-2xl focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 outline-none transition-all font-bold dark:text-white" placeholder="Search by ID, customer, phone..." />
    </div>
    <button class="bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200 dark:border-slate-800 text-slate-900 dark:text-white font-black text-xs uppercase tracking-widest py-4 px-8 rounded-2xl flex items-center gap-3 hover:bg-slate-50 dark:hover:bg-slate-800 transition-all cursor-pointer shadow-sm shadow-blue-900/5">
      <Filter class="w-4 h-4" /> Filter Status
    </button>
  </div>

  <!-- Table Card -->
  <div class="bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200 dark:border-slate-800 rounded-[2.5rem] overflow-hidden shadow-xl shadow-blue-900/5">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-slate-50 dark:bg-slate-950 border-b border-slate-200 dark:border-slate-800 text-slate-400 font-['Fira_Code'] text-[10px] font-black uppercase tracking-[0.2em]">
            <th class="p-8">Ticket ID</th>
            <th class="p-8">Customer</th>
            <th class="p-8">Device</th>
            <th class="p-8">Status</th>
            <th class="p-8">Technician</th>
            <th class="p-8 text-right">Registered</th>
          </tr>
        </thead>
        <tbody class="font-['Fira_Sans'] text-sm">
          <tr class="border-b border-slate-100 dark:border-slate-800 hover:bg-blue-50/30 dark:hover:bg-blue-900/5 transition-all cursor-pointer group">
            <td class="p-8 font-black text-blue-600">TKT-89A2</td>
            <td class="p-8 text-slate-900 dark:text-white font-bold tracking-tight">John Doe</td>
            <td class="p-8 text-slate-500 dark:text-slate-400 font-bold uppercase tracking-tight text-xs">Apple iPhone 13 Pro</td>
            <td class="p-8">
              <span class="bg-amber-50 dark:bg-amber-900/20 text-amber-600 px-3 py-1.5 rounded-full text-[10px] font-black uppercase tracking-widest border border-amber-100 dark:border-amber-800">Diagnosing</span>
            </td>
            <td class="p-8 text-slate-400 font-bold text-xs uppercase tracking-widest">Unassigned</td>
            <td class="p-8 text-slate-400 font-bold text-right text-xs">Today, 10:45 AM</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</div>
