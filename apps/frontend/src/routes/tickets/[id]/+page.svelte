<script lang="ts">
  import { Smartphone, CheckCircle, CheckCircle2, Clock, Download, AlertTriangle, ChevronLeft, Wrench } from 'lucide-svelte';
  
  // Mock data for UI development
  const ticket = {
    id: 'TKT-89A2',
    status: 'diagnosing',
    device: { brand: 'Apple', model: 'iPhone 13 Pro' },
    issue: 'Screen cracked, touch unresponsive. Customer also mentioned battery drains quickly.',
    createdAt: 'May 05, 2026 • 10:45 AM',
    estimatedReady: 'May 07, 2026',
    diagnosis: {
      findings: 'Screen panel is shattered and digitizer is damaged. Needs full screen replacement. Internal flex cable requires reseating.',
      diagnosisFee: 25.00,
      laborFee: 45.00,
      parts: [
        { name: 'OLED Display Panel', grade: 'Original', price: 180.00 }
      ],
      total: 250.00
    }
  };
</script>

<div class="max-w-6xl mx-auto px-4 py-12">
  <!-- Navigation -->
  <div class="mb-8">
    <a href="/tickets" class="inline-flex items-center gap-2 text-blue-600 font-['Fira_Sans'] font-black text-xs uppercase tracking-widest hover:translate-x-[-4px] transition-all duration-300">
      <ChevronLeft size={16} /> Back to Repairs
    </a>
  </div>

  <!-- Header Card -->
  <div class="bg-slate-900 text-white rounded-[2rem] p-10 lg:p-14 mb-10 relative overflow-hidden shadow-2xl border border-slate-800">
    <!-- Decor -->
    <div class="absolute top-0 right-0 w-64 h-64 bg-blue-600/10 rounded-full translate-x-1/3 -translate-y-1/3 blur-3xl"></div>
    
    <div class="relative z-10 flex flex-col md:flex-row justify-between items-start md:items-center gap-8">
      <div>
        <div class="flex items-center gap-4 mb-4">
          <span class="font-['Fira_Code'] text-blue-400 font-black tracking-tight text-lg">{ticket.id}</span>
          <span class="bg-amber-500/10 text-amber-500 text-[10px] font-black px-3 py-1.5 rounded-full uppercase tracking-[0.2em] border border-amber-500/20">{ticket.status}</span>
        </div>
        <h1 class="font-['Fira_Sans'] text-4xl lg:text-5xl font-black mb-2 uppercase tracking-tight">{ticket.device.brand} {ticket.device.model}</h1>
        <p class="font-['Fira_Sans'] text-slate-400 text-lg font-bold">Estimated completion: {ticket.estimatedReady}</p>
      </div>
      
      <button class="px-8 py-4 bg-white/5 hover:bg-white/10 border border-white/10 rounded-2xl text-white font-black text-xs uppercase tracking-widest transition-all backdrop-blur-md flex items-center gap-2">
        <Download size={16} /> Download Receipt
      </button>
    </div>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-3 gap-10">
    <!-- Main Content -->
    <div class="lg:col-span-2 space-y-10">
      <!-- Status Timeline -->
      <div class="bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200 dark:border-slate-800 rounded-[2rem] p-8 lg:p-12 shadow-xl shadow-blue-900/5">
        <h2 class="font-['Fira_Code'] text-xs font-black text-slate-400 uppercase tracking-[0.3em] mb-12">Repair Timeline</h2>
        
        <div class="space-y-12">
          <!-- Step: Received -->
          <div class="flex gap-8 relative">
            <div class="absolute left-4 top-10 bottom-[-48px] w-0.5 bg-blue-600"></div>
            <div class="relative z-10 w-8 h-8 rounded-full bg-blue-600 flex items-center justify-center text-white shadow-lg shadow-blue-600/40">
              <CheckCircle2 class="w-4 h-4" />
            </div>
            <div>
              <h3 class="font-['Fira_Sans'] font-black text-slate-900 dark:text-white text-lg uppercase tracking-tight">Device Received</h3>
              <p class="font-['Fira_Sans'] text-slate-500 dark:text-slate-400 font-bold mt-1">{ticket.createdAt}</p>
              <div class="mt-4 p-5 bg-slate-50 dark:bg-slate-950/50 rounded-2xl border border-slate-100 dark:border-slate-800">
                <p class="text-sm font-bold text-slate-600 dark:text-slate-400">Initial problem: {ticket.issue}</p>
              </div>
            </div>
          </div>

          <!-- Step: Diagnosing -->
          <div class="flex gap-8">
            <div class="relative z-10 w-8 h-8 rounded-full bg-amber-500 flex items-center justify-center text-white shadow-lg shadow-amber-500/40">
              <Clock class="w-4 h-4" />
            </div>
            <div>
              <h3 class="font-['Fira_Sans'] font-black text-slate-900 dark:text-white text-lg uppercase tracking-tight">Under Diagnosis</h3>
              <p class="font-['Fira_Sans'] text-slate-500 dark:text-slate-400 font-bold mt-1">Expected update in 2-4 hours</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Action Required: Diagnosis Complete -->
      {#if ticket.status === 'diagnosing'}
        <div class="bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-amber-500/30 rounded-[2rem] overflow-hidden shadow-xl shadow-amber-900/5">
          <div class="bg-amber-500/10 p-6 border-b border-amber-500/20 flex items-center gap-4 text-amber-600">
            <AlertTriangle class="w-6 h-6" />
            <h2 class="font-['Fira_Code'] font-black text-lg uppercase tracking-tight">Diagnosis Complete</h2>
          </div>
          
          <div class="p-8 lg:p-10">
            <div class="mb-10">
              <h3 class="font-['Fira_Code'] text-[10px] font-black text-slate-400 uppercase tracking-[0.3em] mb-4">Technician Findings</h3>
              <div class="p-6 bg-slate-50 dark:bg-slate-950/50 rounded-2xl border border-slate-100 dark:border-slate-800">
                <p class="font-['Fira_Sans'] text-slate-900 dark:text-white font-bold leading-relaxed italic">
                  "{ticket.diagnosis.findings}"
                </p>
              </div>
            </div>
            
            <div class="mb-10">
              <h3 class="font-['Fira_Code'] text-[10px] font-black text-slate-400 uppercase tracking-[0.3em] mb-4">Cost Breakdown</h3>
              <div class="space-y-4">
                <div class="flex justify-between items-center px-2">
                  <span class="font-['Fira_Sans'] font-bold text-slate-500">Diagnosis Fee</span>
                  <span class="font-['Fira_Code'] font-black text-slate-900 dark:text-white">${ticket.diagnosis.diagnosisFee.toFixed(2)}</span>
                </div>
                <div class="flex justify-between items-center px-2">
                  <span class="font-['Fira_Sans'] font-bold text-slate-500">Labor</span>
                  <span class="font-['Fira_Code'] font-black text-slate-900 dark:text-white">${ticket.diagnosis.laborFee.toFixed(2)}</span>
                </div>
                {#each ticket.diagnosis.parts as part}
                  <div class="flex justify-between items-center px-2">
                    <div class="flex items-center gap-3">
                      <span class="font-['Fira_Sans'] font-black text-slate-900 dark:text-white">{part.name}</span>
                      <span class="text-[10px] font-black px-2 py-0.5 bg-blue-50 dark:bg-blue-900/20 text-blue-600 rounded uppercase tracking-widest border border-blue-100 dark:border-blue-800">{part.grade}</span>
                    </div>
                    <span class="font-['Fira_Code'] font-black text-slate-900 dark:text-white">${part.price.toFixed(2)}</span>
                  </div>
                {/each}
                <div class="flex justify-between items-center pt-6 mt-4 border-t border-slate-100 dark:border-slate-800 px-2">
                  <span class="font-['Fira_Code'] font-black text-slate-900 dark:text-white text-lg uppercase tracking-tight">Total Estimate</span>
                  <span class="font-['Fira_Code'] font-black text-blue-600 text-2xl">${ticket.diagnosis.total.toFixed(2)}</span>
                </div>
              </div>
            </div>
            
            <div class="flex flex-col sm:flex-row gap-4">
              <button class="flex-1 bg-blue-600 hover:bg-blue-700 text-white py-4 rounded-2xl font-black text-xs uppercase tracking-widest shadow-lg shadow-blue-600/20 transition-all cursor-pointer active:scale-[0.98]">
                Approve Repair
              </button>
              <button class="flex-1 bg-white dark:bg-slate-900 border border-red-200 dark:border-red-900/30 text-red-600 py-4 rounded-2xl font-black text-xs uppercase tracking-widest hover:bg-red-50 dark:hover:bg-red-900/10 transition-all cursor-pointer active:scale-[0.98]">
                Decline & Collect
              </button>
            </div>
          </div>
        </div>
      {/if}
    </div>

    <!-- Sidebar -->
    <div class="space-y-8">
      <div class="bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border border-slate-200 dark:border-slate-800 rounded-[2rem] p-8 shadow-xl shadow-blue-900/5">
        <h2 class="font-['Fira_Code'] text-xs font-black text-slate-400 uppercase tracking-[0.3em] mb-6">Device Specs</h2>
        <div class="space-y-6">
          <div>
            <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1 block">Brand / Model</span>
            <p class="font-['Fira_Sans'] font-black text-slate-900 dark:text-white uppercase tracking-tight">{ticket.device.brand} {ticket.device.model}</p>
          </div>
          <div>
            <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1 block">Accessories Left</span>
            <div class="flex flex-wrap gap-2 mt-2">
              <span class="px-3 py-1 bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 text-[10px] font-black rounded-lg uppercase tracking-widest">SIM Card</span>
              <span class="px-3 py-1 bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 text-[10px] font-black rounded-lg uppercase tracking-widest">Clear Case</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
