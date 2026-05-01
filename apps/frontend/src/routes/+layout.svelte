<!-- apps/frontend/src/routes/+layout.svelte -->
<script lang="ts">
  import '../app.css';
  import { Smartphone, Mail, Globe, Code, Share2, ArrowRight } from 'lucide-svelte';
  let { children } = $props();
</script>

<div class="relative flex min-h-screen flex-col bg-background selection:bg-blue-100 selection:text-blue-900">
  <!-- Navbar -->
  <header class="fixed top-0 left-0 right-0 z-50 bg-white/80 dark:bg-slate-950/80 backdrop-blur-md border-b border-slate-200 dark:border-slate-800">
    <div class="container mx-auto px-4">
      <nav class="h-16 flex items-center justify-between">
        <a href="/" class="flex items-center gap-2.5 group">
          <div class="w-9 h-9 bg-blue-600 text-white rounded-lg flex items-center justify-center group-hover:bg-blue-700 transition-colors shadow-sm">
            <Smartphone size={20} />
          </div>
          <span class="text-xl font-bold tracking-tight text-slate-900 dark:text-white">OpenBench</span>
        </a>

        <div class="hidden md:flex items-center gap-8 font-medium text-slate-600 dark:text-slate-400 text-sm">
          <a href="#features" class="hover:text-blue-600 dark:hover:text-blue-400 transition-colors">Features</a>
          <a href="#solutions" class="hover:text-blue-600 dark:hover:text-blue-400 transition-colors">Solutions</a>
          <a href="#pricing" class="hover:text-blue-600 dark:hover:text-blue-400 transition-colors">Pricing</a>
        </div>

        <div class="flex items-center gap-4">
          <button class="hidden sm:block px-4 py-2 font-semibold text-sm text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white transition-colors">Log in</button>
          <button class="px-5 py-2.5 bg-blue-600 text-white font-bold text-sm rounded-lg hover:bg-blue-700 transition-all active:scale-95 shadow-sm">
            Sign up
          </button>
        </div>
      </nav>
    </div>
  </header>

  <main class="flex-1">
    {@render children()}
  </main>

  <!-- Footer -->
  <footer class="bg-slate-950 text-slate-400 py-24 border-t border-slate-900">
    <div class="container mx-auto px-4">
      <div class="grid grid-cols-1 md:grid-cols-12 gap-12 mb-20">
        <!-- Brand Section -->
        <div class="md:col-span-4">
          <a href="/" class="flex items-center gap-2.5 mb-6">
            <div class="w-10 h-10 bg-blue-600 text-white rounded-xl flex items-center justify-center shadow-lg shadow-blue-600/20">
              <Smartphone size={22} />
            </div>
            <span class="text-2xl font-bold tracking-tight text-white">OpenBench</span>
          </a>
          <p class="text-slate-400 text-lg leading-relaxed mb-8 max-w-sm">
            Modernizing the repair industry with professional tools and automated workflows. Built for shops that scale.
          </p>
          <div class="flex gap-4">
            <a href="#" class="w-10 h-10 rounded-lg bg-slate-900 flex items-center justify-center hover:bg-blue-600 hover:text-white transition-all text-slate-500">
              <Globe size={20} />
            </a>
            <a href="#" class="w-10 h-10 rounded-lg bg-slate-900 flex items-center justify-center hover:bg-blue-600 hover:text-white transition-all text-slate-500">
              <Code size={20} />
            </a>
            <a href="#" class="w-10 h-10 rounded-lg bg-slate-900 flex items-center justify-center hover:bg-blue-600 hover:text-white transition-all text-slate-500">
              <Share2 size={20} />
            </a>
          </div>
        </div>

        <!-- Links Sections -->
        <div class="md:col-span-2 md:col-start-6">
          <h4 class="font-bold text-white uppercase tracking-widest text-xs mb-8">Product</h4>
          <ul class="space-y-4 text-sm font-medium">
            <li><a href="#" class="hover:text-white transition-colors">Platform</a></li>
            <li><a href="/track" class="hover:text-white transition-colors">Public Tracker</a></li>
            <li><a href="#" class="hover:text-white transition-colors">Inventory</a></li>
            <li><a href="#" class="hover:text-white transition-colors">API</a></li>
          </ul>
        </div>

        <div class="md:col-span-2">
          <h4 class="font-bold text-white uppercase tracking-widest text-xs mb-8">Support</h4>
          <ul class="space-y-4 text-sm font-medium">
            <li><a href="#" class="hover:text-white transition-colors">Guides</a></li>
            <li><a href="#" class="hover:text-white transition-colors">Security</a></li>
            <li><a href="#" class="hover:text-white transition-colors">Contact</a></li>
            <li><a href="#" class="hover:text-white transition-colors">Privacy</a></li>
          </ul>
        </div>

        <!-- Newsletter Section -->
        <div class="md:col-span-3">
          <h4 class="font-bold text-white uppercase tracking-widest text-xs mb-8">Stay Updated</h4>
          <p class="text-sm text-slate-500 mb-6 leading-relaxed">
            Get the latest updates on features and repair shop best practices.
          </p>
          <div class="relative group">
            <input 
              type="email" 
              placeholder="Enter your email" 
              class="w-full bg-slate-900 border border-slate-800 rounded-xl py-3 px-4 text-sm text-white focus:outline-none focus:border-blue-600 transition-colors"
            />
            <button class="absolute right-2 top-2 bottom-2 bg-blue-600 text-white p-2 rounded-lg hover:bg-blue-700 transition-colors">
              <ArrowRight size={16} />
            </button>
          </div>
        </div>
      </div>

      <!-- Bottom Bar -->
      <div class="pt-8 border-t border-slate-900 flex flex-col md:flex-row justify-between items-center gap-6 text-xs font-medium text-slate-600 uppercase tracking-widest">
        <p>© 2026 OpenBench Technologies Inc. All rights reserved.</p>
        <div class="flex gap-8">
          <a href="#" class="hover:text-white transition-colors">Terms of Service</a>
          <a href="#" class="hover:text-white transition-colors">Cookie Policy</a>
        </div>
      </div>
    </div>
  </footer>
</div>
