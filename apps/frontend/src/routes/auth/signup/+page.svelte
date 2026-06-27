<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, Button, Input } from '$lib';
  import { authService } from '$lib/services/auth';
  import { goto } from '$app/navigation';
  import { UserPlus, ShieldAlert } from 'lucide-svelte';

  let email = $state('');
  let password = $state('');
  let error = $state('');
  let loading = $state(false);
  let hydrated = $state(false);

  onMount(() => {
    hydrated = true;
  });

  async function handleSignUp(e: SubmitEvent) {
    e.preventDefault();
    error = '';

    if (!email.trim()) {
      error = 'Email address is required.';
      return;
    }
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email.trim())) {
      error = 'Please enter a valid email address.';
      return;
    }
    if (!password) {
      error = 'Password is required.';
      return;
    }
    if (password.length < 6) {
      error = 'Password must be at least 6 characters long.';
      return;
    }

    loading = true;
    try {
      // Call authService signUp which automatically logs the user in
      await authService.signUp(email.trim(), password);
      // Redirect to profile setup page upon success
      await goto('/portal/setup');
    } catch (err: any) {
      error = err.message || 'An error occurred during registration.';
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>Sign Up - OpenBench</title>
  <meta name="description" content="Create a new OpenBench account to manage and request device repairs." />
</svelte:head>

<main class="min-h-screen flex items-center justify-center p-4 bg-neubrutalism-bg select-none" data-hydrated={hydrated}>
  <div class="w-full max-w-md">
    
    <!-- Brand / Header Section -->
    <div class="text-center mb-8">
      <a href="/home" class="inline-block group">
        <h1 class="font-display font-black text-4xl tracking-tight text-neubrutalism-charcoal uppercase">
          OPEN<span class="bg-neubrutalism-yellow px-2 py-1 border-4 border-neubrutalism-charcoal shadow-neubrutalism-sm ml-1 group-hover:shadow-neubrutalism-md transition-all">BENCH</span>
        </h1>
      </a>
      <p class="font-mono text-xs mt-3 text-neubrutalism-charcoal opacity-80 uppercase tracking-widest">
        Repair Shop Orchestration Portal
      </p>
    </div>

    <!-- Sign Up Card -->
    <Card class="relative overflow-visible" bgColor="bg-white">
      <!-- Decorative corner tag -->
      <div class="absolute -top-4 -right-4 bg-neubrutalism-pink text-white font-mono font-bold text-xs uppercase px-3 py-1.5 border-4 border-neubrutalism-charcoal shadow-neubrutalism-sm">
        REGISTER
      </div>

      <form onsubmit={handleSignUp} class="flex flex-col gap-6" novalidate>
        <div class="flex flex-col gap-2">
          <h2 class="font-display font-bold text-2xl text-neubrutalism-charcoal">
            SIGN UP
          </h2>
          <p class="font-sans text-sm text-neubrutalism-charcoal opacity-70">
            Create an account to submit repair requests and trace diagnostic logs.
          </p>
        </div>

        {#if error}
          <div class="flex gap-3 items-start border-4 border-neubrutalism-charcoal bg-rose-100 p-4 rounded-none shadow-neubrutalism-sm" role="alert">
            <ShieldAlert class="w-5 h-5 text-neubrutalism-pink shrink-0 mt-0.5" />
            <div class="flex flex-col gap-1">
              <span class="font-sans font-bold text-sm text-neubrutalism-charcoal">Registration Failed</span>
              <span class="font-sans text-xs text-neubrutalism-charcoal">{error}</span>
            </div>
          </div>
        {/if}

        <div class="flex flex-col gap-4">
          <Input
            id="email"
            type="email"
            label="Email Address"
            placeholder="e.g. customer@example.com"
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
              <span class="font-mono text-sm animate-pulse">CREATING ACCOUNT...</span>
            {:else}
              <UserPlus class="w-5 h-5" />
              <span>CREATE ACCOUNT</span>
            {/if}
          </Button>

          <div class="border-t-2 border-dashed border-zinc-200 pt-4 text-center">
            <span class="font-sans text-xs font-semibold text-zinc-500">Already have an account? </span>
            <a href="/auth/signin" class="font-sans text-xs font-bold text-neubrutalism-pink hover:underline uppercase tracking-wide">
              Sign In Here
            </a>
          </div>
        </div>
      </form>
    </Card>

    <!-- Footer links -->
    <div class="text-center mt-6">
      <a href="/home" class="font-mono text-xs font-bold text-neubrutalism-charcoal hover:underline">
        &larr; BACK TO PUBLIC LANDING
      </a>
    </div>

  </div>
</main>
