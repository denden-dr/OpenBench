<script lang="ts">
  import { Search, Plus, AlertTriangle } from 'lucide-svelte';
</script>

<div class="max-w-7xl mx-auto px-4 py-12">
  <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-12 gap-8">
    <div>
      <h1 class="font-['Fira_Code'] text-4xl lg:text-5xl font-black text-slate-900 dark:text-white uppercase tracking-tight">Inventory Lab</h1>
      <p class="font-['Fira_Sans'] text-lg text-slate-500 dark:text-slate-400 font-medium mt-2">Precision parts management & stock control</p>
    </div>
    <button class="bg-emerald-500 hover:bg-emerald-600 text-white font-['Fira_Sans'] font-black text-xs uppercase tracking-widest py-4 px-8 rounded-2xl transition-all shadow-lg shadow-emerald-600/20 cursor-pointer flex items-center gap-3 active:scale-95">
      <Plus class="w-5 h-5" /> Add Part
    </button>
  </div>

  <div class="mb-10 relative w-full md:w-96 group">
    <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-slate-400 group-focus-within:text-blue-600 transition-colors">
      <Search class="h-5 w-5" />
    </div>
    <input type="text" class="block w-full pl-12 pr-4 py-4 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200 dark:border-slate-800 rounded-2xl focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 outline-none transition-all font-bold dark:text-white" placeholder="Search parts..." />
  </div>

  <!-- Inventory Table Card -->
  <div class="bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200 dark:border-slate-800 rounded-[2.5rem] overflow-hidden shadow-xl shadow-blue-900/5">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-slate-50 dark:bg-slate-950 border-b border-slate-200 dark:border-slate-800 text-slate-400 font-['Fira_Code'] text-[10px] font-black uppercase tracking-[0.2em]">
            <th class="p-8">Part Name</th>
            <th class="p-8">Compatibility</th>
            <th class="p-8">Grade</th>
            <th class="p-8">Unit Price</th>
            <th class="p-8 text-right">Available Stock</th>
          </tr>
        </thead>
        <tbody class="font-['Fira_Sans'] text-sm">
          <tr class="border-b border-slate-100 dark:border-slate-800 hover:bg-blue-50/30 dark:hover:bg-blue-900/5 transition-all group">
            <td class="p-8 font-black text-slate-900 dark:text-white uppercase tracking-tight">OLED Display Panel</td>
            <td class="p-8 text-slate-500 dark:text-slate-400 font-bold">iPhone 13 Pro</td>
            <td class="p-8"><span class="bg-blue-50 dark:bg-blue-900/20 text-blue-600 px-3 py-1.5 rounded-full text-[10px] font-black uppercase tracking-widest border border-blue-100 dark:border-blue-800">Original</span></td>
            <td class="p-8 text-slate-900 dark:text-white font-black">$180.00</td>
            <td class="p-8 text-right">
              <span class="flex items-center justify-end gap-3 font-black text-red-600">
                <AlertTriangle class="w-5 h-5" /> 2 Units
              </span>
            </td>
          </tr>
          <tr class="border-b border-slate-100 dark:border-slate-800 hover:bg-blue-50/30 dark:hover:bg-blue-900/5 transition-all group">
            <td class="p-8 font-black text-slate-900 dark:text-white uppercase tracking-tight">Battery Cell 3095mAh</td>
            <td class="p-8 text-slate-500 dark:text-slate-400 font-bold">iPhone 13 Pro</td>
            <td class="p-8"><span class="bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 px-3 py-1.5 rounded-full text-[10px] font-black uppercase tracking-widest border border-slate-200 dark:border-slate-700">ODM</span></td>
            <td class="p-8 text-slate-900 dark:text-white font-black">$45.00</td>
            <td class="p-8 text-right font-black text-emerald-600">14 Units</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</div>
