<script lang="ts">
  import { Mail, Key, User, Phone } from 'lucide-svelte';
</script>

<div class="min-h-[calc(100vh-5rem)] flex items-center justify-center py-12 px-4">
  <div class="w-full max-w-md p-8 border border-slate-200 dark:border-slate-800 bg-white/80 dark:bg-slate-900/80 backdrop-blur-xl rounded-2xl shadow-xl shadow-blue-900/5">
    <div class="text-center mb-8">
      <h1 class="text-3xl font-black text-slate-900 dark:text-white tracking-tight">Sign Up</h1>
      <p class="text-slate-500 dark:text-slate-400 mt-2 font-medium">Create an account to submit your booking</p>
    </div>

    <button class="w-full mb-6 py-3 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl flex items-center justify-center gap-2 hover:bg-slate-50 dark:hover:bg-slate-700 transition-all duration-200 cursor-pointer font-bold text-slate-900 dark:text-white shadow-sm active:scale-[0.98]">
    <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor" xmlns="http://www.w3.org/2000/svg"><path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/><path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.16v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/><path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.16C1.43 8.55 1 10.22 1 12s.43 3.45 1.16 4.93l3.68-2.84z" fill="#FBBC05"/><path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.16 7.07l3.68 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/></svg>
    Sign up with Google
  </button>

  <div class="relative flex py-5 items-center">
    <div class="flex-grow border-t border-slate-200 dark:border-slate-800"></div>
    <span class="flex-shrink-0 mx-4 text-xs font-black uppercase tracking-widest text-slate-400">or</span>
    <div class="flex-grow border-t border-slate-200 dark:border-slate-800"></div>
  </div>

  <form class="space-y-4" on:submit|preventDefault>
    <div>
      <label class="block text-xs font-black uppercase tracking-widest text-slate-500 mb-2" for="name">Full Name</label>
      <div class="relative group">
        <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400 group-focus-within:text-blue-600 transition-colors">
          <User size={18} />
        </div>
        <input type="text" id="name" class="block w-full pl-11 pr-4 py-3 bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 outline-none transition-all dark:text-white" placeholder="John Doe" />
      </div>
    </div>
    
    <div>
      <label class="block text-xs font-black uppercase tracking-widest text-slate-500 mb-2" for="email">Email</label>
      <div class="relative group">
        <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400 group-focus-within:text-blue-600 transition-colors">
          <Mail size={18} />
        </div>
        <input type="email" id="email" class="block w-full pl-11 pr-4 py-3 bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 outline-none transition-all dark:text-white" placeholder="you@example.com" />
      </div>
    </div>

    <div>
      <label class="block text-xs font-black uppercase tracking-widest text-slate-500 mb-2" for="phone">Phone Number</label>
      <div class="relative group">
        <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400 group-focus-within:text-blue-600 transition-colors">
          <Phone size={18} />
        </div>
        <input type="tel" id="phone" class="block w-full pl-11 pr-4 py-3 bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 outline-none transition-all dark:text-white" placeholder="+1 (555) 000-0000" />
      </div>
    </div>
    
    <div>
      <label class="block text-xs font-black uppercase tracking-widest text-slate-500 mb-2" for="password">Password</label>
      <div class="relative group">
        <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-400 group-focus-within:text-blue-600 transition-colors">
          <Key size={18} />
        </div>
        <input type="password" id="password" class="block w-full pl-11 pr-4 py-3 bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded-xl focus:ring-2 focus:ring-blue-600/20 focus:border-blue-600 outline-none transition-all dark:text-white" placeholder="••••••••" />
      </div>
    </div>

    <button type="submit" class="w-full mt-6 py-4 bg-blue-600 text-white font-bold rounded-xl hover:bg-blue-700 transition-all duration-200 cursor-pointer shadow-lg shadow-blue-600/20 active:scale-[0.98]">
      Create Account
    </button>
  </form>

  <p class="mt-8 text-center text-sm text-slate-500 dark:text-slate-400 font-medium">
    Already have an account? <a href="/login" class="text-blue-600 font-bold hover:underline">Log in</a>
  </p>
  </div>
</div>
