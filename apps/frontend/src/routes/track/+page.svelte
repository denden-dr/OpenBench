<script lang="ts">
  import { Search, Smartphone, ShieldCheck, ArrowRight, History, QrCode, ChevronLeft, ChevronRight } from 'lucide-svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/state';
  import { onMount } from 'svelte';

  let ticketId = $state(page.url.searchParams.get('id') || '');
  let phoneLast4 = $state(page.url.searchParams.get('phone') || '');

  onMount(() => {
    if (ticketId && phoneLast4) {
      handleTrack();
    }
  });

  function handleTrack() {
    if (ticketId) {
      goto(`/track/${ticketId}`);
    }
  }

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
</style>
