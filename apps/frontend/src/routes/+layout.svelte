<!-- apps/frontend/src/routes/+layout.svelte -->
<script lang="ts">
  import '../app.css';
  import { page } from '$app/state';
  import { Smartphone, Mail, Globe, Code, Share2, ArrowRight, User, LogOut } from 'lucide-svelte';
  let { children } = $props();

  // Mock login state based on route for UI demo
  let isLoggedIn = $derived(page.url.pathname.startsWith('/profile') || page.url.pathname.startsWith('/dashboard'));
</script>

<div class="relative flex min-h-screen flex-col bg-background selection:bg-blue-100 selection:text-blue-900">
  <!-- Navbar -->
  <header class="fixed top-0 left-0 right-0 z-50 bg-white/70 dark:bg-slate-950/70 backdrop-blur-2xl border-b border-slate-200 dark:border-slate-800 shadow-sm shadow-blue-900/5">
    <div class="container mx-auto px-6">
      <nav class="h-20 flex items-center justify-between">
        <a href="/" class="flex items-center gap-2.5 group">
          <div class="w-9 h-9 bg-blue-600 text-white rounded-lg flex items-center justify-center group-hover:bg-blue-700 transition-colors shadow-sm">
            <Smartphone size={20} />
          </div>
          <span class="text-xl font-bold tracking-tight text-slate-900 dark:text-white">OpenBench</span>
        </a>

        <div class="hidden md:flex items-center gap-8">
          <a href="/pricing" class="text-xs font-black uppercase tracking-widest text-slate-500 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">Services & Fees</a>
          <a href="/track" class="text-xs font-black uppercase tracking-widest text-slate-500 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">Track Status</a>
          <a href="/queue" class="text-xs font-black uppercase tracking-widest text-slate-500 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">Live Queue</a>
        </div>

        <div class="flex items-center gap-6">
          {#if isLoggedIn}
            <a href="/profile" class="flex items-center gap-2.5 px-6 py-2.5 rounded-xl bg-slate-100 dark:bg-slate-900 text-slate-900 dark:text-white font-black text-xs uppercase tracking-widest border border-slate-200 dark:border-slate-800 hover:bg-slate-200 dark:hover:bg-slate-800 transition-all shadow-sm">
              <User size={14} class="text-blue-600" />
              My Profile
            </a>
            <button class="w-11 h-11 flex items-center justify-center rounded-xl border border-slate-200 dark:border-slate-800 text-slate-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/10 transition-all">
              <LogOut size={18} />
            </button>
          {:else}
            <a href="/login" class="hidden sm:block px-6 py-2.5 font-black text-xs uppercase tracking-widest text-slate-500 hover:text-slate-900 dark:hover:text-white transition-colors">Log in</a>
            <a href="/signup" class="px-8 py-3 bg-blue-600 text-white font-black text-xs uppercase tracking-widest rounded-xl hover:bg-blue-700 transition-all active:scale-95 shadow-lg shadow-blue-600/20 inline-flex">
              Sign up
            </a>
          {/if}
        </div>
      </nav>
    </div>
  </header>

  <main class="flex-1 pt-20">
    {@render children()}
  </main>

  <!-- Footer -->
  <footer class="bg-slate-950 text-slate-400 py-24 border-t border-slate-900">
    <div class="container mx-auto px-4">
      <div class="flex flex-col lg:flex-row justify-between gap-16 mb-20">
        <!-- Brand Section -->
        <div class="lg:max-w-sm">
          <a href="/" class="flex items-center gap-2.5 mb-8">
            <div class="w-11 h-11 bg-blue-600 text-white rounded-xl flex items-center justify-center shadow-lg shadow-blue-600/20">
              <Smartphone size={24} />
            </div>
            <span class="text-2xl font-black tracking-tight text-white uppercase">OpenBench</span>
          </a>
          <p class="text-slate-400 text-lg leading-relaxed mb-10">
            The clinical standard for device diagnostics. We combine high-grade parts with transparent, real-time tracking to redefine the repair experience.
          </p>
          <div class="flex gap-3">
            <a href="#" class="w-12 h-12 rounded-xl bg-slate-900 flex items-center justify-center hover:bg-blue-600 hover:text-white transition-all text-slate-500 border border-slate-800">
              <Globe size={20} />
            </a>
            <a href="#" class="w-12 h-12 rounded-xl bg-slate-900 flex items-center justify-center hover:bg-blue-600 hover:text-white transition-all text-slate-500 border border-slate-800">
              <Code size={20} />
            </a>
            <a href="#" class="w-12 h-12 rounded-xl bg-slate-900 flex items-center justify-center hover:bg-blue-600 hover:text-white transition-all text-slate-500 border border-slate-800">
              <Share2 size={20} />
            </a>
          </div>
        </div>

        <!-- Links Grid -->
        <div class="grid grid-cols-2 sm:grid-cols-3 gap-12 lg:gap-24">
          <div>
            <h4 class="font-black text-white uppercase tracking-widest text-[10px] mb-8 opacity-50">Operations</h4>
            <ul class="space-y-4 text-sm font-bold">
              <li><a href="/book" class="hover:text-blue-500 transition-colors">Book Repair</a></li>
              <li><a href="/track" class="hover:text-blue-500 transition-colors">Track Status</a></li>
              <li><a href="/pricing" class="hover:text-blue-500 transition-colors">Fee Schedule</a></li>
              <li><a href="/queue" class="hover:text-blue-500 transition-colors">Live Queue</a></li>
            </ul>
          </div>

          <div>
            <h4 class="font-black text-white uppercase tracking-widest text-[10px] mb-8 opacity-50">Company</h4>
            <ul class="space-y-4 text-sm font-bold">
              <li><a href="#" class="hover:text-blue-500 transition-colors">Our Lab</a></li>
              <li><a href="#" class="hover:text-blue-500 transition-colors">Contact</a></li>
              <li><a href="#" class="hover:text-blue-500 transition-colors">Careers</a></li>
            </ul>
          </div>

          <div class="col-span-2 sm:col-span-1">
            <h4 class="font-black text-white uppercase tracking-widest text-[10px] mb-8 opacity-50">Support</h4>
            <ul class="space-y-4 text-sm font-bold">
              <li><a href="#" class="hover:text-blue-500 transition-colors">Privacy</a></li>
              <li><a href="#" class="hover:text-blue-500 transition-colors">Terms</a></li>
              <li><a href="#" class="hover:text-blue-500 transition-colors">FAQ</a></li>
            </ul>
          </div>
        </div>
      </div>

      <!-- Bottom Bar -->
      <div class="pt-12 border-t border-slate-900 flex flex-col md:flex-row justify-between items-center gap-8 text-[10px] font-black text-slate-600 uppercase tracking-[0.2em]">
        <p>© 2026 OpenBench Technologies. Engineering Trust.</p>
        <div class="flex gap-10">
          <span class="flex items-center gap-2">
            <div class="w-1.5 h-1.5 rounded-full bg-emerald-500"></div>
            Systems Operational
          </span>
          <span class="opacity-50">v1.4.2-clinical</span>
        </div>
      </div>
    </div>
  </footer>
</div>
