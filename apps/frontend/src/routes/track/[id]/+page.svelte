<script lang="ts">
  import { 
    Smartphone, 
    Clock, 
    CheckCircle2, 
    AlertCircle, 
    Calendar, 
    ShieldCheck, 
    Camera, 
    ArrowLeft,
    Wrench,
    ChevronRight,
    History,
    ClipboardCheck,
    Image as ImageIcon
  } from 'lucide-svelte';
  import { page } from '$app/state';

  // Mock data based on PRD schema
  const ticketId = page.params.id;
  const ticket = {
    id: ticketId,
    device: {
      brand: 'Apple',
      model: 'iPhone 15 Pro',
      type: 'Smartphone'
    },
    status: 'repairing', // received, diagnosing, waiting_parts, repairing, ready, completed
    created_at: '2024-04-29T10:00:00Z',
    estimated_ready_at: '2024-04-30T17:00:00Z',
    issue: 'Cracked screen and battery health degradation.',
    diagnosis: 'OLED panel shattered. Battery capacity at 78%. Requires full screen assembly replacement and genuine battery swap.',
    attachments: [
      { type: 'before', url: '/images/hero.png', label: 'Initial Condition' },
      { type: 'after', url: '', label: 'After Repair' }
    ],
    warranty_days: 30,
    parts_grade: 'Original'
  };

  const steps = [
    { id: 'received', label: 'Device Received', description: 'Intake and inspection completed.' },
    { id: 'diagnosing', label: 'Diagnosis', description: 'Technician is identifying the issue.' },
    { id: 'repairing', label: 'Repair in Progress', description: 'Precision repair being performed.' },
    { id: 'ready', label: 'Quality Check', description: 'Final testing before pickup.' },
    { id: 'completed', label: 'Ready for Pickup', description: 'Repair completed successfully.' }
  ];

  const currentStepIndex = steps.findIndex(s => s.id === ticket.status);

  function formatDate(dateStr: string) {
    return new Date(dateStr).toLocaleString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      minute: '2-digit'
    });
  }
</script>

<svelte:head>
  <title>Track Repair #{ticketId} | OpenBench</title>
</svelte:head>

<div class="min-h-screen bg-slate-50 dark:bg-slate-950 pt-24 pb-12">
  <div class="container mx-auto px-4 max-w-4xl">
    <!-- Back Navigation -->
    <div class="mb-8">
      <a href="/" class="inline-flex items-center gap-2 text-sm font-semibold text-slate-500 hover:text-blue-600 transition-colors">
        <ArrowLeft size={16} />
        Back to Home
      </a>
    </div>

    <!-- Header Section -->
    <div class="bg-white dark:bg-slate-900 rounded-3xl p-8 shadow-soft border border-slate-200 dark:border-slate-800 mb-8">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div class="flex items-center gap-4">
          <div class="w-16 h-16 rounded-2xl bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 flex items-center justify-center shadow-sm">
            <Smartphone size={32} />
          </div>
          <div>
            <div class="flex items-center gap-3 mb-1">
              <h1 class="text-2xl font-bold text-slate-900 dark:text-white">{ticket.device.brand} {ticket.device.model}</h1>
              <span class="px-3 py-1 rounded-full bg-blue-100 dark:bg-blue-900/50 text-blue-700 dark:text-blue-300 text-xs font-bold uppercase tracking-wider">
                {ticket.status.replace('_', ' ')}
              </span>
            </div>
            <p class="text-sm font-medium text-slate-500">Ticket ID: <span class="text-slate-900 dark:text-slate-300 font-mono">#{ticket.id}</span></p>
          </div>
        </div>
        
        <div class="flex items-center gap-8 border-l border-slate-100 dark:border-slate-800 md:pl-8">
          <div>
            <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-1">Est. Completion</p>
            <div class="flex items-center gap-2 text-slate-900 dark:text-white font-bold">
              <Clock size={16} class="text-blue-600" />
              <span>{formatDate(ticket.estimated_ready_at)}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="grid lg:grid-cols-3 gap-8">
      <!-- Left Column: Progress Timeline -->
      <div class="lg:col-span-2 space-y-8">
        <!-- Progress Timeline -->
        <div class="bg-white dark:bg-slate-900 rounded-3xl p-8 shadow-soft border border-slate-200 dark:border-slate-800">
          <h2 class="text-lg font-bold text-slate-900 dark:text-white mb-8 flex items-center gap-2">
            <History size={20} class="text-blue-600" />
            Repair Timeline
          </h2>

          <div class="relative space-y-8">
            {#each steps as step, i}
              <div class="flex gap-4 relative">
                <!-- Line -->
                {#if i < steps.length - 1}
                  <div class="absolute left-4 top-8 bottom-[-2rem] w-0.5 bg-slate-100 dark:bg-slate-800">
                    {#if i < currentStepIndex}
                      <div class="h-full w-full bg-blue-600"></div>
                    {/if}
                  </div>
                {/if}

                <!-- Dot -->
                <div class="relative z-10">
                  {#if i < currentStepIndex}
                    <div class="w-8 h-8 rounded-full bg-blue-600 text-white flex items-center justify-center shadow-lg shadow-blue-600/20">
                      <CheckCircle2 size={16} />
                    </div>
                  {:else if i === currentStepIndex}
                    <div class="w-8 h-8 rounded-full bg-blue-50 dark:bg-blue-900/30 border-2 border-blue-600 flex items-center justify-center animate-pulse">
                      <div class="w-2.5 h-2.5 rounded-full bg-blue-600"></div>
                    </div>
                  {:else}
                    <div class="w-8 h-8 rounded-full bg-slate-50 dark:bg-slate-800 border-2 border-slate-200 dark:border-slate-700"></div>
                  {/if}
                </div>

                <!-- Content -->
                <div class="pb-2">
                  <h3 class="font-bold text-sm {i <= currentStepIndex ? 'text-slate-900 dark:text-white' : 'text-slate-400'}">
                    {step.label}
                  </h3>
                  <p class="text-xs text-slate-500 mt-1">{step.description}</p>
                </div>
              </div>
            {/each}
          </div>
        </div>

        <!-- Repair Details -->
        <div class="bg-white dark:bg-slate-900 rounded-3xl p-8 shadow-soft border border-slate-200 dark:border-slate-800">
          <h2 class="text-lg font-bold text-slate-900 dark:text-white mb-6 flex items-center gap-2">
            <ClipboardCheck size={20} class="text-blue-600" />
            Technician Notes
          </h2>
          
          <div class="space-y-6">
            <div>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-2">Reported Issue</p>
              <p class="text-sm text-slate-600 dark:text-slate-400 leading-relaxed bg-slate-50 dark:bg-slate-950/50 p-4 rounded-xl border border-slate-100 dark:border-slate-800">
                {ticket.issue}
              </p>
            </div>
            
            <div>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-2">Diagnosis & Plan</p>
              <p class="text-sm text-slate-600 dark:text-slate-400 leading-relaxed">
                {ticket.diagnosis}
              </p>
            </div>
          </div>
        </div>

        <!-- Visual Evidence -->
        <div class="bg-white dark:bg-slate-900 rounded-3xl p-8 shadow-soft border border-slate-200 dark:border-slate-800">
          <h2 class="text-lg font-bold text-slate-900 dark:text-white mb-6 flex items-center gap-2">
            <Camera size={20} class="text-blue-600" />
            Visual Evidence
          </h2>
          
          <div class="grid sm:grid-cols-2 gap-6">
            {#each ticket.attachments as attachment}
              <div class="group">
                <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-3">{attachment.label}</p>
                <div class="aspect-video rounded-2xl bg-slate-50 dark:bg-slate-950 border border-slate-100 dark:border-slate-800 overflow-hidden relative">
                  {#if attachment.url}
                    <img src={attachment.url} alt={attachment.label} class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
                  {:else}
                    <div class="absolute inset-0 flex flex-col items-center justify-center text-slate-400 gap-2">
                      <ImageIcon size={32} strokeWidth={1} />
                      <span class="text-[10px] font-bold uppercase tracking-widest">Pending Upload</span>
                    </div>
                  {/if}
                </div>
              </div>
            {/each}
          </div>
        </div>
      </div>

      <!-- Right Column: Info Cards -->
      <div class="space-y-8">
        <!-- Warranty Card -->
        <div class="bg-gradient-to-br from-blue-600 to-indigo-700 rounded-3xl p-8 text-white shadow-xl shadow-blue-600/20 relative overflow-hidden">
          <div class="absolute top-0 right-0 w-32 h-32 bg-white/10 rounded-full translate-x-1/2 -translate-y-1/2 blur-2xl"></div>
          <div class="relative z-10">
            <ShieldCheck size={40} class="mb-6 opacity-80" />
            <h3 class="text-xl font-bold mb-2">Service Warranty</h3>
            <p class="text-blue-100 text-sm mb-6 leading-relaxed">This repair is covered by our professional service guarantee.</p>
            
            <div class="space-y-4 pt-6 border-t border-white/10">
              <div class="flex justify-between items-center text-sm">
                <span class="text-blue-200">Duration</span>
                <span class="font-bold">{ticket.warranty_days} Days</span>
              </div>
              <div class="flex justify-between items-center text-sm">
                <span class="text-blue-200">Part Grade</span>
                <span class="font-bold">{ticket.parts_grade}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Help Card -->
        <div class="bg-white dark:bg-slate-900 rounded-3xl p-8 shadow-soft border border-slate-200 dark:border-slate-800">
          <h3 class="font-bold text-slate-900 dark:text-white mb-4">Need Help?</h3>
          <p class="text-sm text-slate-500 dark:text-slate-400 mb-6 leading-relaxed">
            If you have questions about this repair, please contact our support team.
          </p>
          <button class="w-full py-3 rounded-xl border border-slate-200 dark:border-slate-700 text-sm font-bold text-slate-900 dark:text-white hover:bg-slate-50 dark:hover:bg-slate-800 transition-all flex items-center justify-center gap-2">
            Contact Support
            <ChevronRight size={16} />
          </button>
        </div>

        <!-- Security Notice -->
        <div class="flex gap-3 p-4 bg-amber-50 dark:bg-amber-900/10 border border-amber-100 dark:border-amber-900/30 rounded-2xl">
          <AlertCircle size={20} class="text-amber-600 flex-shrink-0" />
          <p class="text-xs text-amber-800 dark:text-amber-200 leading-relaxed">
            This dashboard displays non-sensitive repair data for tracking purposes only. Personal contact details are hidden for your security.
          </p>
        </div>
      </div>
    </div>
  </div>
</div>

<style>
  .shadow-soft {
    box-shadow: 0 10px 30px -10px rgba(0, 0, 0, 0.05);
  }
</style>
