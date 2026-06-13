<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { authService } from '$lib/services/auth';

  interface Props {
    children: any;
  }

  let { children }: Props = $props();
  let authorized = $state(false);

  onMount(async () => {
    const session = await authService.checkSession();
    if (!session || session.role !== 'admin') {
      // User is not authenticated; redirect to the admin sign-in page
      goto('/auth/signin');
    } else {
      authorized = true;
    }
  });
</script>

{#if authorized}
  {@render children()}
{:else}
  <!-- Simple loading fallback during redirection check to prevent layout flashes -->
  <div class="min-h-screen flex items-center justify-center bg-neubrutalism-bg">
    <div class="font-mono text-sm uppercase tracking-widest text-neubrutalism-charcoal animate-pulse">
      Verifying Credentials...
    </div>
  </div>
{/if}
