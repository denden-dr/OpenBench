<script lang="ts">
  import { ArrowLeft } from 'lucide-svelte';
</script>

<!-- Back Link -->
<div>
  <a href="/admin/sales" class="inline-flex items-center gap-1.5 font-mono text-xs font-bold uppercase hover:underline">
    <ArrowLeft class="w-4 h-4" />
    <span>BACK TO POS CHECKOUT</span>
  </a>
</div>

<!-- Header Area -->
<div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 border-b-4 border-neubrutalism-charcoal pb-4">
  <div>
    <h2 class="font-display font-extrabold text-2xl md:text-3xl uppercase tracking-tight">
      Sales Transaction Archive
    </h2>
    <p class="font-sans text-xs sm:text-sm text-zinc-500 font-semibold mt-1">
      Logs of completed retail transactions, total turnover metrics, and receipts.
    </p>
  </div>
</div>
