<script lang="ts">
  import { Clock, Wrench, CheckCircle, Smartphone, SmartphoneIcon } from 'lucide-svelte';
  
  const columns = [
    { id: 'received', title: 'Received', color: 'bg-slate-400', textColor: 'text-slate-500', hoverBorder: 'hover:border-slate-400', glow: 'from-slate-400/5', bgLight: 'bg-slate-400/10' },
    { id: 'diagnosing', title: 'Diagnosing', color: 'bg-amber-500', textColor: 'text-amber-600', pulse: true, hoverBorder: 'hover:border-amber-500', glow: 'from-amber-500/5', bgLight: 'bg-amber-500/10' },
    { id: 'repairing', title: 'Repairing', color: 'bg-blue-600', textColor: 'text-blue-600', hoverBorder: 'hover:border-blue-600', glow: 'from-blue-600/5', bgLight: 'bg-blue-600/10' },
    { id: 'ready', title: 'Ready', color: 'bg-emerald-500', textColor: 'text-emerald-600', hoverBorder: 'hover:border-emerald-500', glow: 'from-emerald-500/5', bgLight: 'bg-emerald-500/10' }
  ];

  const tickets = [
    // Received
    { id: 'TKT-89A2', brand: 'Apple', model: 'iPhone 13 Pro', issue: 'Screen Replacement', status: 'received', time: '15m ago' },
    { id: 'TKT-22B9', brand: 'Google', model: 'Pixel 8', issue: 'Charging Port', status: 'received', time: '45m ago' },
    { id: 'TKT-55X1', brand: 'Sony', model: 'Xperia 1 V', issue: 'Back Glass', status: 'received', time: '1h ago' },
    
    // Diagnosing
    { id: 'TKT-44B1', brand: 'Samsung', model: 'S22 Ultra', issue: 'Water Damage', status: 'diagnosing', label: 'Active' },
    { id: 'TKT-77R4', brand: 'Apple', model: 'iPad Pro M2', issue: 'FaceID Failure', status: 'diagnosing', label: 'Queued' },
    
    // Repairing
    { id: 'TKT-12C8', brand: 'Google', model: 'Pixel 6', issue: 'Battery Service', status: 'repairing', label: 'Parts In' },
    { id: 'TKT-90K2', brand: 'Apple', model: 'iPhone 14', issue: 'Rear Camera', status: 'repairing', label: 'In Progress' },
    { id: 'TKT-33M9', brand: 'Samsung', model: 'Fold 5', issue: 'Hinge Alignment', status: 'repairing', label: 'Complex' },
    
    // Ready
    { id: 'TKT-99F5', brand: 'Apple', model: 'MacBook Pro M1', issue: 'Keyboard Repair', status: 'ready', label: 'Tested' },
    { id: 'TKT-11Z0', brand: 'Xiaomi', model: '14 Ultra', issue: 'Lens Polish', status: 'ready', label: 'Ready' },
    { id: 'TKT-66Q3', brand: 'Apple', model: 'iPhone 12', issue: 'Audio IC', status: 'ready', label: 'Tested' },
    { id: 'TKT-88P2', brand: 'Google', model: 'Pixel Fold', issue: 'Inner Screen', status: 'ready', label: 'Quality Check' }
  ];

  const getTicketsByStatus = (status: string) => tickets.filter(t => t.status === status);
</script>

<div class="min-h-screen bg-white dark:bg-slate-950 pt-32 pb-20">
  <div class="max-w-7xl mx-auto px-4">
    <div class="mb-12 flex flex-col md:flex-row items-center justify-between gap-6 bg-slate-50 dark:bg-slate-900/50 p-6 rounded-3xl border border-slate-200 dark:border-slate-800">
    <div>
      <h1 class="text-3xl lg:text-4xl font-black text-slate-900 dark:text-white mb-2 uppercase tracking-tight">Live Repair Queue</h1>
      <p class="text-sm text-slate-500 dark:text-slate-400 font-medium">Real-time status of all devices currently in our clinical environment.</p>
    </div>
    
    <div class="flex items-center gap-4">
      <div class="text-right">
        <p class="text-[10px] font-black uppercase tracking-widest text-slate-400">Last Updated</p>
        <p class="text-sm font-bold text-slate-900 dark:text-white">Just now</p>
      </div>
      <div class="h-10 w-px bg-slate-200 dark:bg-slate-700"></div>
      <div class="inline-flex items-center gap-2 px-4 py-2 rounded-xl bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-200 dark:border-emerald-800/50 text-emerald-600 dark:text-emerald-400">
        <div class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></div>
        <span class="text-xs font-black uppercase tracking-widest">System Active</span>
      </div>
    </div>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 lg:gap-10 items-start">
    {#each columns as column}
      <div class="flex flex-col max-h-none lg:max-h-[calc(100vh-16rem)]">
        <!-- Column Header (Fixed) -->
        <div class="flex items-center gap-3 px-2 mb-6 shrink-0">
          <div class="w-2.5 h-2.5 rounded-full {column.color} {column.pulse ? 'animate-pulse' : ''}"></div>
          <h2 class="font-black text-[11px] uppercase tracking-[0.3em] {column.textColor}">{column.title}</h2>
          <span class="ml-auto bg-slate-100 dark:bg-slate-900 text-slate-500 dark:text-slate-400 px-2.5 py-1 rounded-lg text-[9px] font-black tracking-widest uppercase border border-slate-200 dark:border-slate-800">
            {getTicketsByStatus(column.id).length} units
          </span>
        </div>
        
        <!-- Column Content (Scrollable) -->
        <div class="space-y-4 overflow-y-auto pr-2 pb-8 flex-1 scrollbar-thin scrollbar-thumb-slate-200 dark:scrollbar-thumb-slate-800 scrollbar-track-transparent">
          {#each getTicketsByStatus(column.id) as ticket}
            <div class="group bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 p-4 rounded-2xl shadow-sm hover:shadow-md transition-all duration-300 cursor-pointer hover:-translate-y-1 {column.hoverBorder} relative overflow-hidden">
              <!-- Glow accent on hover -->
              <div class="absolute inset-0 bg-gradient-to-br {column.glow} to-transparent opacity-0 group-hover:opacity-100 transition-opacity"></div>
              
              <div class="relative z-10">
                <div class="flex justify-between items-start mb-3">
                  <div class="flex items-center gap-2">
                    <div class="w-1.5 h-1.5 rounded-full {column.color}"></div>
                    <span class="font-mono font-bold text-slate-700 dark:text-slate-300 text-xs uppercase">{ticket.id}</span>
                  </div>
                  <span class="text-[9px] font-black uppercase tracking-widest {column.textColor} {column.bgLight} px-1.5 py-0.5 rounded flex items-center justify-center">
                    {ticket.time || ticket.label}
                  </span>
                </div>
                
                <h3 class="font-bold text-slate-900 dark:text-white mb-1 leading-tight text-sm truncate" title="{ticket.brand} {ticket.model}">{ticket.brand} {ticket.model}</h3>
                
                <div class="flex items-center justify-between mt-2 pt-2 border-t border-slate-100 dark:border-slate-800/50">
                  <div class="flex items-center gap-1.5 text-slate-500 dark:text-slate-400">
                    <Wrench size={10} />
                    <p class="text-[10px] font-semibold uppercase tracking-wider truncate max-w-[120px]">{ticket.issue}</p>
                  </div>
                </div>
              </div>
            </div>
          {/each}

          {#if getTicketsByStatus(column.id).length === 0}
            <div class="py-8 flex flex-col items-center justify-center text-center border-2 border-dashed border-slate-200 dark:border-slate-800 rounded-2xl bg-slate-50/50 dark:bg-slate-900/20">
              <div class="w-8 h-8 mb-2 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center">
                <CheckCircle size={14} class="text-slate-400" />
              </div>
              <p class="text-[10px] font-black text-slate-400 dark:text-slate-500 uppercase tracking-widest">Station Clear</p>
            </div>
          {/if}
        </div>
      </div>
    {/each}
  </div>
</div>
</div>

<style>
  /* Custom Scrollbar Utilities (if tailwind-scrollbar plugin is not present) */
  .scrollbar-thin::-webkit-scrollbar {
    width: 6px;
  }
  .scrollbar-thin::-webkit-scrollbar-track {
    background: transparent;
  }
  .scrollbar-thin::-webkit-scrollbar-thumb {
    background-color: #cbd5e1;
    border-radius: 20px;
  }
  :global(.dark) .scrollbar-thin::-webkit-scrollbar-thumb {
    background-color: #1e293b;
  }
</style>
