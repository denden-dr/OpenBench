<script lang="ts">
  import { ArrowLeft } from 'lucide-svelte';
</script>

<div>
  <a href="/" class="inline-flex items-center gap-1.5 font-mono text-xs font-bold uppercase hover:underline">
    <ArrowLeft class="w-4 h-4" />
    <span>BACK TO HOME</span>
  </a>
</div>

<!-- Header Title -->
<div class="text-center flex flex-col gap-2">
  <h1 class="font-display font-black text-3xl md:text-5xl uppercase tracking-tight text-neubrutalism-charcoal">
    PUBLIC REPAIR <span class="bg-neubrutalism-green border-4 border-neubrutalism-charcoal px-2.5 py-0.5 shadow-neubrutalism-md inline-block transform -rotate-1">TRACKER</span>
  </h1>
  <p class="font-mono text-xs text-zinc-500 uppercase mt-1">
    Track gadget repair progress instantly & securely using your UUID code.
  </p>
</div>
