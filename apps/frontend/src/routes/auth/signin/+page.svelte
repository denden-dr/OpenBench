<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, Button, Input } from '$lib';
  import { authService } from '$lib/services/auth';
  import { goto } from '$app/navigation';
  import { KeyRound, ShieldAlert } from 'lucide-svelte';

  let email = $state('');
  let password = $state('');
  let error = $state('');
  let loading = $state(false);
  let hydrated = $state(false);

  onMount(() => {
    hydrated = true;
  });


  async function handleSignIn(e: SubmitEvent) {
    e.preventDefault();
    error = '';
    
    if (!email) {
      error = 'Email address is required.';
      return;
    }
    if (!password) {
      error = 'Password is required.';
      return;
    }

    loading = true;
    try {
      await authService.signIn(email, password);
      // On success, redirect to the admin dashboard
      await goto('/admin');
    } catch (err: any) {
      error = err.message || 'An unexpected error occurred.';
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>Admin Sign In - OpenBench</title>
  <meta name="description" content="Sign in to the OpenBench Repair Shop administration portal." />
</svelte:head>

<main class="min-h-screen flex items-center justify-center p-4 bg-neubrutalism-bg select-none" data-hydrated={hydrated}>
  <div class="w-full max-w-md">
    
    <!-- Brand / Header Section -->
    <div class="text-center mb-8">
      <h1 class="font-display font-bold text-4xl tracking-tight text-neubrutalism-charcoal">
        OPEN<span class="bg-neubrutalism-yellow px-2 py-1 border-4 border-neubrutalism-charcoal shadow-neubrutalism-sm ml-1">BENCH</span>
      </h1>
      <p class="font-mono text-sm mt-3 text-neubrutalism-charcoal opacity-80 uppercase tracking-widest">
        Repair Shop Administration Portal
      </p>
    </div>

    <!-- Login Card -->
    <Card class="relative overflow-visible" bgColor="bg-white">
      <!-- Decorative corner tag -->
      <div class="absolute -top-4 -right-4 bg-neubrutalism-pink text-white font-mono font-bold text-xs uppercase px-3 py-1.5 border-4 border-neubrutalism-charcoal shadow-neubrutalism-sm">
        Admin Portal
      </div>

      <form onsubmit={handleSignIn} class="flex flex-col gap-6" novalidate>
        <div class="flex flex-col gap-2">
          <h2 class="font-display font-bold text-2xl text-neubrutalism-charcoal">
            SIGN IN
          </h2>
          <p class="font-sans text-sm text-neubrutalism-charcoal opacity-70">
            Enter your credentials below to access the workbench dashboard.
          </p>
        </div>

        {#if error}
          <div class="flex gap-3 items-start border-4 border-neubrutalism-charcoal bg-rose-100 p-4 rounded-none shadow-neubrutalism-sm" role="alert">
            <ShieldAlert class="w-5 h-5 text-neubrutalism-pink shrink-0 mt-0.5" />
            <div class="flex flex-col gap-1">
              <span class="font-sans font-bold text-sm text-neubrutalism-charcoal">Authentication Failed</span>
              <span class="font-sans text-xs text-neubrutalism-charcoal">{error}</span>
            </div>
          </div>
        {/if}

        <div class="flex flex-col gap-4">
          <Input
            id="email"
            type="email"
            label="Email Address"
            placeholder="admin@openbench.dev"
            required
            bind:value={email}
            disabled={loading}
          />

          <Input
            id="password"
            type="password"
            label="Password"
            placeholder="••••••••"
            required
            bind:value={password}
            disabled={loading}
          />
        </div>

        <div class="mt-2 flex flex-col gap-4">
          <Button
            type="submit"
            bgColor="bg-neubrutalism-yellow"
            class="w-full flex items-center justify-center gap-2"
            disabled={loading}
          >
            {#if loading}
              <span class="font-mono text-sm animate-pulse">AUTHENTICATING...</span>
            {:else}
              <KeyRound class="w-5 h-5" />
              <span>ACCESS WORKBENCH</span>
            {/if}
          </Button>
          
          <div class="text-center">
            <p class="font-mono text-xs text-neubrutalism-charcoal opacity-60">
              Demo Credentials: admin@openbench.dev / SecureAdminPassword123!
            </p>
          </div>
        </div>
      </form>
    </Card>

    <!-- Footer links -->
    <div class="text-center mt-6">
      <a href="/" class="font-mono text-xs font-bold text-neubrutalism-charcoal hover:underline">
        &larr; BACK TO PUBLIC LANDING
      </a>
    </div>

  </div>
</main>
