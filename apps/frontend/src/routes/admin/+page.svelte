<script lang="ts">
  import { Card, Button } from '$lib';
  import { authService, isMockEnabled } from '$lib/services/auth';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { 
    LogOut, Wrench, CheckCircle2, Clock, ClipboardList, 
    AlertTriangle, Shield, User, ArrowUpRight 
  } from 'lucide-svelte';

  let adminEmail = $state('');
  
  onMount(() => {
    const session = authService.getSession();
    if (session) {
      adminEmail = session.email;
    }
  });

  async function handleLogout() {
    await authService.signOut();
    await goto('/auth/signin');
  }

  // Mock list of recent repairs/tickets for dashboard feel
  const recentTickets = [
    { id: 'TX-4092', customer: 'Alice Cooper', device: 'iPhone 15 Pro Max', issue: 'Shattered Rear Glass', status: 'In Progress', date: '2026-06-11' },
    { id: 'TX-4091', customer: 'John Doe', device: 'MacBook Pro 16" M3', issue: 'Liquid Damage Clean', status: 'Pending Estimate', date: '2026-06-11' },
    { id: 'TX-4090', customer: 'Sarah Connor', device: 'Samsung S24 Ultra', issue: 'Battery Replacement', status: 'Completed', date: '2026-06-10' }
  ];

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'Completed': return 'bg-neubrutalism-green';
      case 'In Progress': return 'bg-neubrutalism-yellow';
      case 'Pending Estimate': return 'bg-neubrutalism-pink text-white';
      default: return 'bg-white';
    }
  };
</script>

<svelte:head>
  <title>Admin Dashboard - OpenBench</title>
</svelte:head>

<div class="min-h-screen bg-neubrutalism-bg flex flex-col font-sans select-none pb-12">
  
  <!-- Header Bar -->
  <header class="border-b-4 border-neubrutalism-charcoal bg-white sticky top-0 z-40">
    <div class="max-w-7xl mx-auto px-4 py-4 sm:px-6 lg:px-8 flex justify-between items-center">
      <div class="flex items-center gap-2">
        <div class="bg-neubrutalism-yellow border-4 border-neubrutalism-charcoal p-1 font-display font-bold text-xl shadow-neubrutalism-sm">
          OB
        </div>
        <span class="font-display font-bold text-2xl tracking-tight text-neubrutalism-charcoal">
          OPENBENCH <span class="text-xs font-mono font-bold bg-neubrutalism-charcoal text-white px-2 py-0.5 ml-1">ADMIN</span>
        </span>
      </div>
      
      <div class="flex items-center gap-4">
        <!-- Logged in status -->
        <div class="hidden md:flex items-center gap-2 border-4 border-neubrutalism-charcoal bg-white py-1 px-3 shadow-neubrutalism-sm font-mono text-xs">
          <Shield class="w-4 h-4 text-neubrutalism-green" />
          <span>{adminEmail || 'admin@openbench.com'}</span>
        </div>

        <Button 
          bgColor="bg-neubrutalism-pink" 
          onclick={handleLogout} 
          class="flex items-center gap-2 py-2 px-4 text-sm"
        >
          <LogOut class="w-4 h-4 text-white" />
          <span class="text-white">LOGOUT</span>
        </Button>
      </div>
    </div>
  </header>

  <!-- Main Workspace -->
  <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mt-8 flex-grow">
    
    <!-- Welcome banner -->
    <div class="mb-8 border-4 border-neubrutalism-charcoal bg-neubrutalism-yellow p-6 shadow-neubrutalism-md">
      <h1 class="font-display font-bold text-3xl md:text-4xl text-neubrutalism-charcoal uppercase leading-none mb-2">
        Welcome to the Workbench, Admin!
      </h1>
      <p class="font-sans text-sm md:text-base text-neubrutalism-charcoal max-w-2xl opacity-90">
        This is the administrative dashboard. You are currently operating inside a secure database-backed environment connected directly to the Go Fiber API.
      </p>
    </div>

    <!-- Metrics Cards Grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      
      <!-- Metric 1 -->
      <Card bgColor="bg-white" class="flex flex-col gap-2 relative">
        <div class="absolute top-4 right-4 bg-neubrutalism-charcoal text-white p-2 border-2 border-neubrutalism-charcoal">
          <Wrench class="w-5 h-5" />
        </div>
        <span class="font-mono text-xs font-bold text-neubrutalism-charcoal opacity-70 uppercase tracking-widest">Active Tickets</span>
        <span class="font-mono text-5xl font-extrabold text-neubrutalism-charcoal mt-2">14</span>
        <span class="font-sans text-xs text-neubrutalism-charcoal opacity-60 mt-2 flex items-center gap-1">
          <AlertTriangle class="w-3.5 h-3.5 text-amber-500" />
          3 critical updates pending action
        </span>
      </Card>

      <!-- Metric 2 -->
      <Card bgColor="bg-white" class="flex flex-col gap-2 relative">
        <div class="absolute top-4 right-4 bg-neubrutalism-charcoal text-white p-2 border-2 border-neubrutalism-charcoal">
          <CheckCircle2 class="w-5 h-5" />
        </div>
        <span class="font-mono text-xs font-bold text-neubrutalism-charcoal opacity-70 uppercase tracking-widest">Completed Today</span>
        <span class="font-mono text-5xl font-extrabold text-neubrutalism-charcoal mt-2">08</span>
        <span class="font-sans text-xs text-neubrutalism-charcoal opacity-60 mt-2 flex items-center gap-1">
          <ArrowUpRight class="w-3.5 h-3.5 text-emerald-500" />
          +12% increase from weekly average
        </span>
      </Card>

      <!-- Metric 3 -->
      <Card bgColor="bg-white" class="flex flex-col gap-2 relative">
        <div class="absolute top-4 right-4 bg-neubrutalism-charcoal text-white p-2 border-2 border-neubrutalism-charcoal">
          <Clock class="w-5 h-5" />
        </div>
        <span class="font-mono text-xs font-bold text-neubrutalism-charcoal opacity-70 uppercase tracking-widest">Pending Estimate</span>
        <span class="font-mono text-5xl font-extrabold text-neubrutalism-charcoal mt-2">03</span>
        <span class="font-sans text-xs text-neubrutalism-charcoal opacity-60 mt-2 flex items-center gap-1">
          <ClipboardList class="w-3.5 h-3.5 text-zinc-500" />
          Awaiting customer confirmation
        </span>
      </Card>
      
    </div>

    <!-- Bottom Section Layout: Main Content (Recent Repairs) & Sidebar (System Info) -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      
      <!-- Left side: Recent Repairs list -->
      <div class="lg:col-span-2 flex flex-col gap-6">
        <h2 class="font-display font-bold text-2xl text-neubrutalism-charcoal uppercase tracking-tight flex items-center gap-2">
          <span>Active Repairs Sandbox</span>
        </h2>
        
        <div class="flex flex-col gap-4">
          {#each recentTickets as ticket}
            <Card bgColor="bg-white" class="hover:shadow-neubrutalism-lg transition-all">
              <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <div class="flex flex-col gap-1">
                  <div class="flex items-center gap-2">
                    <span class="font-mono text-sm font-bold bg-zinc-200 px-2 py-0.5 border-2 border-neubrutalism-charcoal text-neubrutalism-charcoal">
                      {ticket.id}
                    </span>
                    <span class="font-display font-bold text-lg">{ticket.device}</span>
                  </div>
                  <p class="font-sans text-sm text-neubrutalism-charcoal mt-1">
                    <span class="font-semibold">Issue:</span> {ticket.issue}
                  </p>
                  <p class="font-sans text-xs text-neubrutalism-charcoal opacity-60">
                    Customer: {ticket.customer} &bull; Received on {ticket.date}
                  </p>
                </div>
                
                <div class="flex items-center gap-3 w-full sm:w-auto justify-between sm:justify-end">
                  <span class="font-mono text-xs font-bold py-1.5 px-3 border-2 border-neubrutalism-charcoal shadow-neubrutalism-sm uppercase tracking-wide {getStatusColor(ticket.status)}">
                    {ticket.status}
                  </span>
                  
                  <button class="p-2 border-2 border-neubrutalism-charcoal hover:bg-zinc-100 transition duration-150 cursor-pointer shadow-neubrutalism-sm active:translate-x-0.5 active:translate-y-0.5 active:shadow-none bg-white">
                    <ArrowUpRight class="w-4 h-4 text-neubrutalism-charcoal" />
                  </button>
                </div>
              </div>
            </Card>
          {/each}
        </div>
      </div>

      <!-- Right side: Sidebar info -->
      <div class="flex flex-col gap-6">
        <h2 class="font-display font-bold text-2xl text-neubrutalism-charcoal uppercase tracking-tight">
          System Info
        </h2>
        
        <Card bgColor="bg-white" class="flex flex-col gap-4">
          <div class="flex flex-col gap-1 border-b-2 border-dashed border-zinc-300 pb-3">
            <span class="font-mono text-xs text-neubrutalism-charcoal opacity-60 uppercase">Environment</span>
            <span class="font-mono text-sm font-bold text-neubrutalism-pink">{isMockEnabled() ? 'SANDBOX / MOCK MODE' : 'LIVE API INTEGRATION'}</span>
          </div>

          <div class="flex flex-col gap-1 border-b-2 border-dashed border-zinc-300 pb-3">
            <span class="font-mono text-xs text-neubrutalism-charcoal opacity-60 uppercase">App Shell</span>
            <span class="font-sans text-sm font-bold">SvelteKit v2 (SPA Client)</span>
          </div>

          <div class="flex flex-col gap-1 border-b-2 border-dashed border-zinc-300 pb-3">
            <span class="font-mono text-xs text-neubrutalism-charcoal opacity-60 uppercase">Reactivity Engine</span>
            <span class="font-sans text-sm font-bold">Svelte 5 Runes ($state, $props)</span>
          </div>

          <div class="flex flex-col gap-1 pb-1">
            <span class="font-mono text-xs text-neubrutalism-charcoal opacity-60 uppercase">Tailwind Engine</span>
            <span class="font-sans text-sm font-bold">Tailwind CSS v4 (CSS theme)</span>
          </div>
        </Card>
      </div>

    </div>

  </main>
</div>
