<script lang="ts">
  import { Button } from "$lib";
  import { Menu, X } from "lucide-svelte";
  import { onMount } from "svelte";
  import { authService } from "$lib/services/auth";
  import { goto } from "$app/navigation";

  let userSession = $state<any | null>(null);
  let isChecking = $state(true);
  let mobileMenuOpen = $state(false);

  onMount(async () => {
    userSession = await authService.checkSession();
    isChecking = false;
  });

  async function handleSignOut() {
    await authService.signOut();
    userSession = null;
    await goto("/home");
  }
</script>

<header class="border-b-4 border-neubrutalism-charcoal bg-white sticky top-0 z-50">
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-20 flex justify-between items-center">
    <!-- Brand Logo -->
    <a href="/home" class="flex items-center gap-2.5 group">
      <div class="bg-neubrutalism-yellow border-4 border-neubrutalism-charcoal p-1.5 font-display font-black text-xl shadow-neubrutalism-sm group-hover:-translate-x-0.5 group-hover:-translate-y-0.5 group-hover:shadow-neubrutalism-md transition-all">
        OB
      </div>
      <span class="font-display font-black text-2xl tracking-tight uppercase">
        OPEN<span class="text-neubrutalism-pink">BENCH</span>
      </span>
    </a>

    <!-- Desktop & Mobile CTA Button -->
    <div class={userSession ? "block" : "hidden md:block"}>
      {#if isChecking}
        <div class="w-32 h-10 bg-zinc-200 animate-pulse"></div>
      {:else if userSession}
        {#if userSession.role === "user"}
          <a href="/portal">
            <Button bgColor="bg-neubrutalism-pink" class="py-2 px-5 text-sm text-white shadow-neubrutalism-sm font-display font-bold">
              MY PORTAL
            </Button>
          </a>
        {:else}
          <a href="/admin">
            <Button bgColor="bg-neubrutalism-yellow" class="py-2 px-5 text-sm shadow-neubrutalism-sm font-display font-bold">
              ADMIN CONSOLE
            </Button>
          </a>
        {/if}
      {:else}
        <div class="flex items-center gap-3">
          <a href="/auth/signin">
            <Button bgColor="bg-white" class="py-2 px-5 text-sm border-2 border-neubrutalism-charcoal shadow-neubrutalism-sm font-display font-bold hover:bg-zinc-50">
              SIGN IN
            </Button>
          </a>
          <a href="/auth/signup">
            <Button bgColor="bg-neubrutalism-pink" class="py-2 px-5 text-sm text-white shadow-neubrutalism-sm font-display font-bold">
              SIGN UP
            </Button>
          </a>
        </div>
      {/if}
    </div>

    <!-- Mobile menu button (only visible when logged out) -->
    {#if !isChecking && !userSession}
      <button
        class="md:hidden p-2 border-2 border-neubrutalism-charcoal bg-white shadow-neubrutalism-sm active:translate-x-0.5 active:translate-y-0.5 active:shadow-none"
        onclick={() => (mobileMenuOpen = !mobileMenuOpen)}
        aria-label="Toggle Navigation Menu"
      >
        {#if mobileMenuOpen}
          <X class="w-6 h-6" />
        {:else}
          <Menu class="w-6 h-6" />
        {/if}
      </button>
    {/if}
  </div>

  <!-- Mobile Drawer Dropdown (only active when logged out) -->
  {#if mobileMenuOpen && !userSession}
    <div class="md:hidden border-t-4 border-neubrutalism-charcoal bg-white flex flex-col p-6 gap-4 font-display font-bold text-base shadow-neubrutalism-md">
      <a
        href="/auth/signin"
        onclick={() => (mobileMenuOpen = false)}
        class="p-2 border-2 border-transparent hover:bg-neubrutalism-yellow/10 hover:border-neubrutalism-charcoal transition-all"
      >
        SIGN IN
      </a>
      <a
        href="/auth/signup"
        onclick={() => (mobileMenuOpen = false)}
        class="p-2 border-2 border-transparent hover:bg-neubrutalism-pink/10 hover:border-neubrutalism-charcoal transition-all"
      >
        SIGN UP
      </a>
    </div>
  {/if}
</header>
