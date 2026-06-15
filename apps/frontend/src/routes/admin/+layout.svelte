<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/state'; // SvelteKit v2 / Svelte 5 state-based routing
  import { authService } from '$lib/services/auth';
  import { Button } from '$lib';
  import { 
    LayoutDashboard, Wrench, Package, ShoppingCart, 
    ShieldCheck, Settings, Menu, X, LogOut 
  } from 'lucide-svelte';

  interface Props {
    children: any;
  }

  let { children }: Props = $props();
  let authorized = $state(false);
  let sidebarOpen = $state(false);
  let adminEmail = $state('');

  onMount(async () => {
    const session = await authService.checkSession();
    if (!session || session.role !== 'admin') {
      goto('/auth/signin');
    } else {
      adminEmail = session.email;
      authorized = true;
    }
  });

  async function handleLogout() {
    await authService.signOut();
    await goto('/auth/signin');
  }

  // Helper to determine active nav class
  function isLinkActive(path: string): boolean {
    const currentPath = page.url.pathname;
    if (path === '/admin') {
      return currentPath === '/admin';
    }
    return currentPath.startsWith(path);
  }
</script>

{#if authorized}
  <div class="min-h-screen bg-neubrutalism-bg flex font-sans select-none text-neubrutalism-charcoal">
    
    <!-- Sidebar Navigation -->
    <!-- Mobile overlay background -->
    {#if sidebarOpen}
      <button 
        class="fixed inset-0 bg-black/40 z-40 lg:hidden cursor-default transition-opacity"
        onclick={() => sidebarOpen = false}
        aria-label="Close Sidebar"
      ></button>
    {/if}

    <aside class="fixed inset-y-0 left-0 bg-white border-r-4 border-neubrutalism-charcoal w-64 z-50 flex flex-col transform transition-transform duration-200 ease-in-out lg:translate-x-0 lg:sticky lg:top-0 lg:h-screen lg:flex-shrink-0 {sidebarOpen ? 'translate-x-0' : '-translate-x-full'}">
      
      <!-- Brand Logo Header -->
      <div class="p-6 border-b-4 border-neubrutalism-charcoal flex items-center justify-between bg-neubrutalism-yellow">
        <div class="flex items-center gap-2">
          <div class="bg-white border-4 border-neubrutalism-charcoal p-1 font-display font-bold text-lg shadow-neubrutalism-sm">
            OB
          </div>
          <span class="font-display font-bold text-xl tracking-tight">
            OPENBENCH <span class="text-[10px] font-mono font-bold bg-neubrutalism-charcoal text-white px-1.5 py-0.5 ml-0.5">ADMIN</span>
          </span>
        </div>
        
        <!-- Mobile close button -->
        <button class="lg:hidden p-1 border-2 border-neubrutalism-charcoal bg-white" onclick={() => sidebarOpen = false}>
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Navigation Links -->
      <nav class="flex-grow p-4 flex flex-col gap-2 overflow-y-auto">
        <a 
          href="/admin" 
          onclick={() => sidebarOpen = false}
          class="flex items-center gap-3 p-3 font-display font-bold text-sm border-2 border-transparent transition-all {isLinkActive('/admin') ? 'bg-neubrutalism-yellow border-neubrutalism-charcoal shadow-neubrutalism-sm -translate-x-0.5 -translate-y-0.5 hover:bg-amber-400' : 'hover:bg-zinc-100'}"
        >
          <LayoutDashboard class="w-5 h-5 shrink-0" />
          <span>OVERVIEW</span>
        </a>

        <a 
          href="/admin/tickets" 
          onclick={() => sidebarOpen = false}
          class="flex items-center gap-3 p-3 font-display font-bold text-sm border-2 border-transparent transition-all {isLinkActive('/admin/tickets') ? 'bg-neubrutalism-green border-neubrutalism-charcoal shadow-neubrutalism-sm -translate-x-0.5 -translate-y-0.5 hover:bg-emerald-400' : 'hover:bg-zinc-100'}"
        >
          <Wrench class="w-5 h-5 shrink-0" />
          <span>REPAIR TICKETS</span>
        </a>

        <a 
          href="/admin/inventory" 
          onclick={() => sidebarOpen = false}
          class="flex items-center gap-3 p-3 font-display font-bold text-sm border-2 border-transparent transition-all {isLinkActive('/admin/inventory') ? 'bg-neubrutalism-pink text-white border-neubrutalism-charcoal shadow-neubrutalism-sm -translate-x-0.5 -translate-y-0.5 hover:bg-pink-500' : 'hover:bg-zinc-100'}"
        >
          <Package class="w-5 h-5 shrink-0 {isLinkActive('/admin/inventory') ? 'text-white' : ''}" />
          <span class="{isLinkActive('/admin/inventory') ? 'text-white' : ''}">INVENTORY / STOCK</span>
        </a>

        <a 
          href="/admin/sales" 
          onclick={() => sidebarOpen = false}
          class="flex items-center gap-3 p-3 font-display font-bold text-sm border-2 border-transparent transition-all {isLinkActive('/admin/sales') ? 'bg-neubrutalism-yellow border-neubrutalism-charcoal shadow-neubrutalism-sm -translate-x-0.5 -translate-y-0.5 hover:bg-amber-400' : 'hover:bg-zinc-100'}"
        >
          <ShoppingCart class="w-5 h-5 shrink-0" />
          <span>POINT OF SALES</span>
        </a>

        <a 
          href="/admin/warranties" 
          onclick={() => sidebarOpen = false}
          class="flex items-center gap-3 p-3 font-display font-bold text-sm border-2 border-transparent transition-all {isLinkActive('/admin/warranties') ? 'bg-zinc-200 border-neubrutalism-charcoal shadow-neubrutalism-sm -translate-x-0.5 -translate-y-0.5 hover:bg-zinc-300' : 'hover:bg-zinc-100'}"
        >
          <ShieldCheck class="w-5 h-5 shrink-0" />
          <span>WARRANTIES</span>
        </a>

        <a 
          href="/admin/settings" 
          onclick={() => sidebarOpen = false}
          class="flex items-center gap-3 p-3 font-display font-bold text-sm border-2 border-transparent transition-all {isLinkActive('/admin/settings') ? 'bg-zinc-200 border-neubrutalism-charcoal shadow-neubrutalism-sm -translate-x-0.5 -translate-y-0.5 hover:bg-zinc-300' : 'hover:bg-zinc-100'}"
        >
          <Settings class="w-5 h-5 shrink-0" />
          <span>SETTINGS</span>
        </a>
      </nav>

      <!-- Sidebar Footer / Logout -->
      <div class="p-4 border-t-4 border-neubrutalism-charcoal flex flex-col gap-3">
        <div class="flex items-center gap-2 border-2 border-neubrutalism-charcoal bg-zinc-100 py-1.5 px-3 font-mono text-[10px] truncate">
          <ShieldCheck class="w-3.5 h-3.5 text-neubrutalism-green shrink-0" />
          <span class="truncate">{adminEmail || 'admin@openbench.dev'}</span>
        </div>

        <Button 
          bgColor="bg-neubrutalism-pink" 
          onclick={handleLogout} 
          class="flex items-center justify-center gap-2 py-2.5 px-4 text-xs font-bold text-white shadow-neubrutalism-sm"
        >
          <LogOut class="w-4 h-4 text-white" />
          <span class="text-white">LOGOUT ADMIN</span>
        </Button>
      </div>

    </aside>

    <!-- Main Workspace Content wrapper -->
    <div class="flex-grow flex flex-col min-w-0">
      
      <!-- Top header bar (mobile and general tools) -->
      <header class="border-b-4 border-neubrutalism-charcoal bg-white sticky top-0 z-30">
        <div class="px-4 py-4 sm:px-6 flex justify-between items-center h-16">
          <div class="flex items-center gap-4">
            <!-- Sidebar toggle button for mobile -->
            <button 
              class="lg:hidden p-2 border-2 border-neubrutalism-charcoal bg-white shadow-neubrutalism-sm active:translate-x-0.5 active:translate-y-0.5 active:shadow-none"
              onclick={() => sidebarOpen = !sidebarOpen}
              aria-label="Toggle Sidebar"
            >
              <Menu class="w-5 h-5" />
            </button>
            
            <h1 class="font-display font-extrabold text-lg sm:text-xl uppercase tracking-tight">
              Workstation Console
            </h1>
          </div>
          
          <div class="flex items-center gap-3">
            <span class="font-mono text-xs font-bold bg-neubrutalism-green border-2 border-neubrutalism-charcoal py-1 px-2.5 shadow-neubrutalism-sm">
              MOCK MODE ACTIVE
            </span>
          </div>
        </div>
      </header>

      <!-- Sub-page view -->
      <main class="flex-grow p-4 sm:p-6 lg:p-8 overflow-y-auto max-w-7xl w-full mx-auto pb-16">
        {@render children()}
      </main>

    </div>

  </div>
{:else}
  <!-- Simple loading fallback during redirection check to prevent layout flashes -->
  <div class="min-h-screen flex items-center justify-center bg-neubrutalism-bg">
    <div class="font-mono text-sm uppercase tracking-widest text-neubrutalism-charcoal animate-pulse">
      Verifying Credentials...
    </div>
  </div>
{/if}
