<script lang="ts">
  import { Search, Smartphone, ShieldCheck, ArrowRight, History, QrCode, ChevronLeft, ChevronRight } from 'lucide-svelte';
  import { goto } from '$app/navigation';

  let ticketId = $state('');
  let phoneLast4 = $state('');
  let currentPage = $state(0);
  const itemsPerPage = 5;

  const allActivity = [
    { id: 'FIX-98231', brand: 'Apple', model: 'iPhone 15 Pro', issue: 'Screen Replacement', status: 'repairing', time: '2 mins ago' },
    { id: 'FIX-98230', brand: 'Samsung', model: 'Galaxy S24 Ultra', issue: 'Battery Service', status: 'ready', time: '15 mins ago' },
    { id: 'FIX-98229', brand: 'Google', model: 'Pixel 8 Pro', issue: 'Charging Port', status: 'diagnosing', time: '1 hour ago' },
    { id: 'FIX-98228', brand: 'Apple', model: 'iPhone 13', issue: 'Water Damage', status: 'received', time: '3 hours ago' },
    { id: 'FIX-98227', brand: 'Apple', model: 'iPad Pro M2', issue: 'FaceID Repair', status: 'waiting_parts', time: '5 hours ago' },
    { id: 'FIX-98226', brand: 'Sony', model: 'Xperia 1 V', issue: 'Back Glass', status: 'repairing', time: '6 hours ago' },
    { id: 'FIX-98225', brand: 'Xiaomi', model: '14 Ultra', issue: 'Camera Lens', status: 'ready', time: '8 hours ago' },
    { id: 'FIX-98224', brand: 'Apple', model: 'iPhone 14 Plus', issue: 'Software Fix', status: 'diagnosing', time: '12 hours ago' },
    { id: 'FIX-98223', brand: 'Samsung', model: 'Galaxy Z Fold 5', issue: 'Hinge Service', status: 'waiting_parts', time: '1 day ago' },
    { id: 'FIX-98222', brand: 'Google', model: 'Pixel Fold', issue: 'Inner Screen', status: 'repairing', time: '1 day ago' }
  ];

  const totalPages = Math.ceil(allActivity.length / itemsPerPage);
  const displayedActivity = $derived(allActivity.slice(currentPage * itemsPerPage, (currentPage + 1) * itemsPerPage));

  function handleTrack() {
    if (ticketId) {
      goto(`/track/${ticketId}`);
    }
  }

  function nextPage() {
    if (currentPage < totalPages - 1) currentPage++;
  }

  function prevPage() {
    if (currentPage > 0) currentPage--;
  }

  const getStatusColorClass = (status: string) => {
    switch (status) {
      case 'repairing': return 'status-blue';
      case 'ready': return 'status-green';
      case 'diagnosing': return 'status-amber';
      case 'received': return 'status-slate';
      case 'waiting_parts': return 'status-purple';
      default: return '';
    }
  };
</script>

<svelte:head>
  <title>Track Your Repair | OpenBench</title>
</svelte:head>

<div class="min-h-screen bg-white dark:bg-slate-950 pt-32 pb-20 overflow-hidden relative">
  <!-- Background Decor -->
  <div class="absolute top-0 right-0 w-1/3 h-full opacity-5 pointer-events-none hidden lg:block">
    <div class="w-full h-full bg-gradient-to-l from-blue-600 to-transparent"></div>
  </div>

  <div class="container mx-auto px-4 relative z-10">
    <div class="max-w-2xl mx-auto mb-24">
      <div class="text-center mb-12">
        <div class="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-blue-50 dark:bg-blue-900/30 text-sm font-semibold text-blue-600 dark:text-blue-400 mb-6">
          <QrCode size={16} />
          Real-time Transparency
        </div>
        <h1 class="text-4xl lg:text-5xl font-extrabold text-slate-900 dark:text-white mb-6">Track Your Device</h1>
        <p class="text-lg text-slate-600 dark:text-slate-400">
          Enter your Ticket ID and the last 4 digits of your phone number to see the current status of your repair.
        </p>
      </div>

      <!-- Search Card -->
      <div class="bg-white dark:bg-slate-900 p-8 lg:p-12 rounded-[2.5rem] shadow-premium border border-slate-200 dark:border-slate-800 relative group transition-all">
        <div class="absolute -inset-1 bg-gradient-to-r from-blue-600 to-indigo-600 rounded-[2.6rem] blur opacity-0 group-hover:opacity-10 transition duration-1000"></div>
        
        <form onsubmit={(e) => { e.preventDefault(); handleTrack(); }} class="relative space-y-6">
          <div class="space-y-2">
            <label for="ticket-id" class="text-xs font-bold text-slate-400 uppercase tracking-widest block ml-1">Ticket ID</label>
            <div class="relative">
              <div class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400">
                <Search size={18} />
              </div>
              <input 
                id="ticket-id"
                type="text" 
                bind:value={ticketId}
                placeholder="e.g. FIX-98231" 
                class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl py-4 pl-12 pr-4 text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 dark:focus:border-blue-500 transition-all text-lg font-medium"
                required
              />
            </div>
          </div>

          <div class="space-y-2">
            <label for="phone" class="text-xs font-bold text-slate-400 uppercase tracking-widest block ml-1">Phone Number (Last 4 Digits)</label>
            <div class="relative">
              <div class="absolute left-4 top-1/2 -translate-y-1/2 text-slate-400">
                <Smartphone size={18} />
              </div>
              <input 
                id="phone"
                type="tel" 
                bind:value={phoneLast4}
                placeholder="xxxx" 
                maxlength="4"
                class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl py-4 pl-12 pr-4 text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 dark:focus:border-blue-500 transition-all text-lg font-medium"
                required
              />
            </div>
          </div>

          <button 
            type="submit"
            class="w-full bg-blue-600 text-white py-5 rounded-2xl font-bold text-lg shadow-lg shadow-blue-600/20 hover:bg-blue-700 transition-all active:scale-[0.98] flex items-center justify-center gap-3"
          >
            Track Repair Progress
            <ArrowRight size={20} />
          </button>
        </form>
      </div>
    </div>

    <!-- Live Shop Activity Table -->
    <div class="max-w-6xl mx-auto">
      <div class="flex items-center justify-between mb-8">
        <div>
          <h2 class="text-2xl font-bold text-slate-900 dark:text-white">Live Shop Activity</h2>
          <p class="text-sm text-slate-500 dark:text-slate-400">Ongoing repairs across our service hubs.</p>
        </div>
        <div class="hidden sm:flex items-center gap-2 px-3 py-1 bg-green-50 dark:bg-green-900/20 text-green-600 dark:text-green-400 rounded-full text-xs font-bold uppercase tracking-wider">
          <span class="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span>
          Live Updates
        </div>
      </div>

      <div class="bg-white dark:bg-slate-900 rounded-[2.5rem] shadow-soft border border-slate-100 dark:border-slate-800 overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full text-left border-separate border-spacing-0">
            <thead>
              <tr class="bg-slate-50 dark:bg-slate-800/50">
                <th class="px-8 py-5 text-[10px] font-bold text-slate-400 uppercase tracking-widest">Brand</th>
                <th class="px-8 py-5 text-[10px] font-bold text-slate-400 uppercase tracking-widest">Model</th>
                <th class="px-8 py-5 text-[10px] font-bold text-slate-400 uppercase tracking-widest hidden sm:table-cell">Issue</th>
                <th class="px-8 py-5 text-[10px] font-bold text-slate-400 uppercase tracking-widest hidden sm:table-cell">Status</th>
                <th class="px-8 py-5 text-[10px] font-bold text-slate-400 uppercase tracking-widest text-right hidden sm:table-cell">Updated</th>
              </tr>
            </thead>
            <tbody class="bg-white dark:bg-slate-900">
              {#each displayedActivity as item}
                <tr 
                  onclick={() => goto(`/track/${item.id}`)}
                  class="cursor-pointer transition-all duration-300 group {getStatusColorClass(item.status)} glow-row"
                >
                  <td class="px-8 py-6 relative z-10 first-td">
                    <span class="text-sm font-bold text-slate-900 dark:text-white">{item.brand}</span>
                  </td>
                  <td class="px-8 py-6 relative z-10">
                    <span class="text-sm font-medium text-slate-600 dark:text-slate-300">{item.model}</span>
                  </td>
                  <td class="px-8 py-6 hidden sm:table-cell relative z-10">
                    <span class="text-xs text-slate-500 dark:text-slate-400">{item.issue}</span>
                  </td>
                  <td class="px-8 py-6 hidden sm:table-cell relative z-10">
                    <span class="inline-flex px-3 py-1 rounded-full text-[10px] font-bold uppercase tracking-wider
                      {item.status === 'repairing' ? 'bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400' : ''}
                      {item.status === 'ready' ? 'bg-green-50 dark:bg-green-900/30 text-green-600 dark:text-green-400' : ''}
                      {item.status === 'diagnosing' ? 'bg-amber-50 dark:bg-amber-900/30 text-amber-600 dark:text-amber-400' : ''}
                      {item.status === 'received' ? 'bg-slate-50 dark:bg-slate-800 text-slate-600 dark:text-slate-400' : ''}
                      {item.status === 'waiting_parts' ? 'bg-purple-50 dark:bg-purple-900/30 text-purple-600 dark:text-purple-400' : ''}
                    ">
                      {item.status.replace('_', ' ')}
                    </span>
                  </td>
                  <td class="px-8 py-6 text-right hidden sm:table-cell relative z-10 last-td">
                    <span class="text-xs text-slate-400">{item.time}</span>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>

        <!-- Pagination Controls -->
        <div class="px-8 py-6 border-t border-slate-50 dark:border-slate-800 flex items-center justify-between bg-slate-50/30 dark:bg-slate-800/20">
          <p class="text-xs font-medium text-slate-500">
            Showing <span class="text-slate-900 dark:text-white font-bold">{currentPage * itemsPerPage + 1}-{Math.min((currentPage + 1) * itemsPerPage, allActivity.length)}</span> of {allActivity.length} repairs
          </p>
          <div class="flex items-center gap-4">
            <button 
              onclick={prevPage}
              disabled={currentPage === 0}
              class="p-2.5 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-800 text-slate-600 dark:text-slate-300 shadow-sm disabled:opacity-20 hover:border-blue-500 dark:hover:border-blue-500 hover:text-blue-600 dark:hover:text-blue-400 transition-all group"
              aria-label="Previous page"
            >
              <ChevronLeft size={20} class="transition-transform group-hover:-translate-x-0.5" />
            </button>
            
            <div class="flex items-center bg-slate-100 dark:bg-slate-800/50 px-3 py-1.5 rounded-lg border border-slate-200/50 dark:border-slate-700/50">
              <span class="text-xs font-bold text-slate-900 dark:text-white uppercase tracking-tighter">Page {currentPage + 1}</span>
            </div>

            <button 
              onclick={nextPage}
              disabled={currentPage === totalPages - 1}
              class="group relative p-2.5 rounded-xl bg-slate-900 dark:bg-blue-600 text-white disabled:opacity-20 hover:bg-slate-800 dark:hover:bg-blue-500 transition-all shadow-md flex items-center gap-2 pl-4 pr-3"
            >
              <span class="text-xs font-bold uppercase tracking-widest">Next</span>
              <ChevronRight size={20} class="transition-transform group-hover:translate-x-0.5" />
              <!-- Subtle Pulse for Next button -->
              {#if currentPage < totalPages - 1}
                <span class="absolute inset-0 rounded-xl bg-blue-400/20 animate-ping pointer-events-none"></span>
              {/if}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Info Grid -->
    <div class="max-w-4xl mx-auto grid sm:grid-cols-2 gap-8 mt-24">
      <div class="flex gap-4">
        <div class="w-12 h-12 rounded-xl bg-slate-50 dark:bg-slate-900 border border-slate-100 dark:border-slate-800 flex items-center justify-center text-blue-600 flex-shrink-0">
          <History size={24} />
        </div>
        <div>
          <h4 class="font-bold text-slate-900 dark:text-white mb-1">Live History</h4>
          <p class="text-sm text-slate-500 dark:text-slate-400">See every step of the repair process, from intake to quality check.</p>
        </div>
      </div>
      <div class="flex gap-4">
        <div class="w-12 h-12 rounded-xl bg-slate-50 dark:bg-slate-900 border border-slate-100 dark:border-slate-800 flex items-center justify-center text-blue-600 flex-shrink-0">
          <ShieldCheck size={24} />
        </div>
        <div>
          <h4 class="font-bold text-slate-900 dark:text-white mb-1">Privacy First</h4>
          <p class="text-sm text-slate-500 dark:text-slate-400">Your personal data is encrypted and never displayed on public tracking pages.</p>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .shadow-premium {
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.08);
  }

  .shadow-soft {
    box-shadow: 0 10px 30px -10px rgba(0, 0, 0, 0.05);
  }

  @keyframes glow-pulse {
    0%, 100% { 
      outline-color: rgba(var(--glow-rgb), 0.3);
      box-shadow: 0 0 15px -2px rgba(var(--glow-rgb), 0.2);
    }
    50% { 
      outline-color: rgba(var(--glow-rgb), 0.8);
      box-shadow: 0 0 25px 2px rgba(var(--glow-rgb), 0.4);
    }
  }

  .glow-row {
    transition: all 0.3s ease;
    border-bottom: 1px solid rgba(0, 0, 0, 0.02);
    outline: 1px solid transparent;
    outline-offset: -1px;
  }

  .dark .glow-row {
    border-bottom: 1px solid rgba(255, 255, 255, 0.02);
  }

  .glow-row td {
    transition: all 0.3s ease;
  }

  .glow-row:hover {
    position: relative;
    z-index: 20;
    background: rgba(var(--glow-rgb), 0.04);
    animation: glow-pulse 2s ease-in-out infinite;
  }

  .status-blue { --glow-color: #3b82f6; --glow-rgb: 59, 130, 246; }
  .status-green { --glow-color: #22c55e; --glow-rgb: 34, 197, 94; }
  .status-amber { --glow-color: #f59e0b; --glow-rgb: 245, 158, 11; }
  .status-slate { --glow-color: #64748b; --glow-rgb: 100, 116, 139; }
  .status-purple { --glow-color: #a855f7; --glow-rgb: 168, 85, 247; }
</style>
