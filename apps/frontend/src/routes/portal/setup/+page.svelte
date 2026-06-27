<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, Button, Input } from '$lib';
  import { authService } from '$lib/services/auth';
  import { goto } from '$app/navigation';
  import { User, Phone, UserCheck, ShieldAlert } from 'lucide-svelte';

  let username = $state('');
  let fullName = $state('');
  let phoneNumber = $state('');
  let error = $state('');
  let loading = $state(false);
  let hydrated = $state(false);

  onMount(async () => {
    hydrated = true;
    
    // Guard route: user must be logged in to access profile completion
    const session = authService.getSession();
    if (!session) {
      await goto('/auth/signin');
    }
  });

  async function handleProfileSubmit(e: SubmitEvent) {
    e.preventDefault();
    error = '';

    if (!username.trim()) {
      error = 'Username is required.';
      return;
    }
    if (!fullName.trim()) {
      error = 'Full name is required.';
      return;
    }

    loading = true;
    try {
      await authService.updateProfile({
        username: username.trim(),
        full_name: fullName.trim(),
        phone_number: phoneNumber.trim() || undefined
      });
      // Redirect to main portal dashboard
      await goto('/portal');
    } catch (err: any) {
      error = err.message || 'Failed to update profile. Please try again.';
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>Setup Profile - OpenBench</title>
  <meta name="description" content="Complete your OpenBench customer profile." />
</svelte:head>

<main class="min-h-screen flex items-center justify-center p-4 bg-neubrutalism-bg select-none" data-hydrated={hydrated}>
  <div class="w-full max-w-md">
    
    <!-- Brand / Header Section -->
    <div class="text-center mb-8">
      <h1 class="font-display font-black text-4xl tracking-tight text-neubrutalism-charcoal uppercase">
        OPEN<span class="bg-neubrutalism-yellow px-2 py-1 border-4 border-neubrutalism-charcoal shadow-neubrutalism-sm ml-1">BENCH</span>
      </h1>
      <p class="font-mono text-xs mt-3 text-neubrutalism-charcoal opacity-80 uppercase tracking-widest">
        Onboarding Portal
      </p>
    </div>

    <!-- Setup Card -->
    <Card class="relative overflow-visible" bgColor="bg-white">
      <!-- Corner tag -->
      <div class="absolute -top-4 -right-4 bg-neubrutalism-yellow text-neubrutalism-charcoal font-mono font-bold text-xs uppercase px-3 py-1.5 border-4 border-neubrutalism-charcoal shadow-neubrutalism-sm">
        STEP 2/2
      </div>

      <form onsubmit={handleProfileSubmit} class="flex flex-col gap-6" novalidate>
        <div class="flex flex-col gap-2">
          <h2 class="font-display font-bold text-2xl text-neubrutalism-charcoal">
            COMPLETE PROFILE
          </h2>
          <p class="font-sans text-sm text-neubrutalism-charcoal opacity-70">
            Tell us a bit about yourself so we can link and trace your repair records.
          </p>
        </div>

        {#if error}
          <div class="flex gap-3 items-start border-4 border-neubrutalism-charcoal bg-rose-100 p-4 rounded-none shadow-neubrutalism-sm" role="alert">
            <ShieldAlert class="w-5 h-5 text-neubrutalism-pink shrink-0 mt-0.5" />
            <div class="flex flex-col gap-1">
              <span class="font-sans font-bold text-sm text-neubrutalism-charcoal">Setup Failed</span>
              <span class="font-sans text-xs text-neubrutalism-charcoal">{error}</span>
            </div>
          </div>
        {/if}

        <div class="flex flex-col gap-4">
          <Input
            id="fullName"
            type="text"
            label="Full Name *"
            placeholder="e.g. Tony Stark"
            required
            bind:value={fullName}
            disabled={loading}
          />

          <Input
            id="username"
            type="text"
            label="Username *"
            placeholder="e.g. ironman"
            required
            bind:value={username}
            disabled={loading}
          />

          <Input
            id="phoneNumber"
            type="tel"
            label="Phone Number (Optional)"
            placeholder="e.g. +628123456789"
            bind:value={phoneNumber}
            disabled={loading}
          />
        </div>

        <div class="mt-2 flex flex-col gap-4">
          <Button
            type="submit"
            bgColor="bg-neubrutalism-pink"
            class="w-full flex items-center justify-center gap-2 text-white"
            disabled={loading}
          >
            {#if loading}
              <span class="font-mono text-sm animate-pulse">SAVING PROFILE...</span>
            {:else}
              <UserCheck class="w-5 h-5" />
              <span>SAVE & ENTER PORTAL</span>
            {/if}
          </Button>
        </div>
      </form>
    </Card>

  </div>
</main>
