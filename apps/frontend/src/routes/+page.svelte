<script lang="ts">
  import { Card, Button } from '$lib';
  import { authService } from '$lib/services/auth';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { Shield, ArrowRight, Wrench, Settings, ClipboardCheck } from 'lucide-svelte';

  let isAuthenticated = $state(false);
  let isChecking = $state(true);

  onMount(async () => {
    const session = await authService.checkSession();
    isAuthenticated = session !== null && session.role === 'admin';
    isChecking = false;
  });

  function handleNavigate() {
    if (isAuthenticated) {
      goto('/admin');
    } else {
      goto('/auth/signin');
    }
  }
</script>

<svelte:head>
  <title>OpenBench - Repair Shop Management</title>
</svelte:head>

<main class="min-h-screen bg-neubrutalism-bg flex flex-col font-sans select-none">
  
  <!-- Outer Hero Container -->
  <div class="flex-grow flex items-center justify-center p-4">
    <div class="w-full max-w-2xl text-center flex flex-col gap-8">
      
      <!-- Big Bold Header -->
      <div class="flex flex-col gap-4">
        <h1 class="font-display font-black text-5xl md:text-7xl uppercase text-neubrutalism-charcoal tracking-tight">
          OPEN<span class="bg-neubrutalism-yellow border-4 border-neubrutalism-charcoal px-3 py-1 shadow-neubrutalism-md inline-block transform rotate-1">BENCH</span>
        </h1>
        <p class="font-mono text-sm md:text-base text-neubrutalism-charcoal opacity-80 uppercase tracking-widest mt-2">
          State-of-the-art repair shop orchestration engine
        </p>
      </div>

      <!-- Feature Highlight Cards -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-left">
        <Card bgColor="bg-white" class="p-4 flex flex-col gap-2">
          <div class="bg-neubrutalism-green border-2 border-neubrutalism-charcoal p-1.5 w-fit">
            <Wrench class="w-4 h-4 text-neubrutalism-charcoal" />
          </div>
          <h3 class="font-display font-bold text-sm uppercase">Track Repairs</h3>
          <p class="font-sans text-xs text-neubrutalism-charcoal opacity-70">Log devices, diagnostics, and components in real time.</p>
        </Card>

        <Card bgColor="bg-white" class="p-4 flex flex-col gap-2">
          <div class="bg-neubrutalism-pink border-2 border-neubrutalism-charcoal p-1.5 w-fit text-white">
            <ClipboardCheck class="w-4 h-4" />
          </div>
          <h3 class="font-display font-bold text-sm uppercase">Quick Estimates</h3>
          <p class="font-sans text-xs text-neubrutalism-charcoal opacity-70">Generate precise customer invoices and digital approvals.</p>
        </Card>

        <Card bgColor="bg-white" class="p-4 flex flex-col gap-2">
          <div class="bg-neubrutalism-yellow border-2 border-neubrutalism-charcoal p-1.5 w-fit">
            <Settings class="w-4 h-4 text-neubrutalism-charcoal" />
          </div>
          <h3 class="font-display font-bold text-sm uppercase">Smooth Workflow</h3>
          <p class="font-sans text-xs text-neubrutalism-charcoal opacity-70">Automated transitions and updates keep teams synchronized.</p>
        </Card>
      </div>

      <!-- Main Action Callout -->
      <Card bgColor="bg-white" class="p-6 md:p-8 flex flex-col md:flex-row items-center justify-between gap-6">
        <div class="text-left flex flex-col gap-1">
          <h2 class="font-display font-bold text-xl uppercase tracking-tight text-neubrutalism-charcoal">
            Administration Workbench
          </h2>
          <p class="font-sans text-sm text-neubrutalism-charcoal opacity-70">
            Sign in as an administrator to edit tickets, record inventory, and handle estimates.
          </p>
        </div>

        {#if isChecking}
          <div class="h-14 w-full md:w-auto px-8 border-4 border-neubrutalism-charcoal bg-zinc-200 animate-pulse flex items-center justify-center font-mono text-xs font-bold">
            CHECKING SESSION...
          </div>
        {:else}
          <Button 
            bgColor={isAuthenticated ? 'bg-neubrutalism-green' : 'bg-neubrutalism-yellow'}
            onclick={handleNavigate}
            class="w-full md:w-auto shrink-0 flex items-center justify-center gap-2 group py-3 px-6"
          >
            {#if isAuthenticated}
              <Shield class="w-5 h-5" />
              <span>GOTO DASHBOARD</span>
            {:else}
              <span>ACCESS WORKBENCH</span>
              <ArrowRight class="w-5 h-5 group-hover:translate-x-1 transition-transform" />
            {/if}
          </Button>
        {/if}
      </Card>

      <!-- Copyright Footer -->
      <div class="mt-4">
        <p class="font-mono text-xs text-neubrutalism-charcoal opacity-50 uppercase">
          OpenBench &bull; Clean Slate Phase 1 Development
        </p>
      </div>

    </div>
  </div>

</main>
