<!-- apps/frontend/src/routes/+layout.svelte -->
<script lang="ts">
  import '../app.css';
  import { Smartphone } from 'lucide-svelte';
  let { children } = $props();
</script>

<div class="relative flex min-h-screen flex-col bg-slate-50 dark:bg-slate-900 selection:bg-blue-100 selection:text-blue-900">
  <!-- Navbar -->
  <header class="fixed top-0 left-0 right-0 z-50 bg-white/70 dark:bg-slate-950/70 backdrop-blur-2xl border-b border-slate-200 dark:border-slate-800 shadow-sm shadow-blue-900/5">
    <div class="container mx-auto px-6">
      <nav class="h-16 flex items-center justify-between">
        <a href="/" class="flex items-center gap-2.5 group">
          <div class="w-8 h-8 bg-blue-600 text-white rounded-lg flex items-center justify-center group-hover:bg-blue-700 transition-colors shadow-sm">
            <Smartphone size={18} />
          </div>
          <span class="text-lg font-bold tracking-tight text-slate-900 dark:text-white">OpenBench <span class="text-xs font-semibold px-2 py-0.5 ml-2 bg-blue-50 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400 rounded-full border border-blue-200/50">Repair Admin</span></span>
        </a>
      </nav>
    </div>
  </header>

  <main class="flex-1 pt-16 pb-12">
    {@render children()}
  </main>

  <!-- Footer -->
  <footer class="bg-white dark:bg-slate-950 border-t border-slate-200 dark:border-slate-800 py-6 text-center text-xs text-slate-500 dark:text-slate-400">
    <div class="container mx-auto px-6 flex flex-col sm:flex-row justify-between items-center gap-4">
      <p>© 2026 OpenBench Technologies. All rights reserved.</p>
      <div class="flex items-center gap-2">
        <div class="w-1.5 h-1.5 rounded-full bg-emerald-500"></div>
        <span class="font-medium text-[10px] uppercase tracking-wider">Local Database Mode</span>
      </div>
    </div>
  </footer>
</div>
