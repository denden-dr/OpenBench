<script lang="ts">
  import { Button } from '$lib';
  import { ShoppingCart, FileText } from 'lucide-svelte';
</script>

<div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 border-b-4 border-neubrutalism-charcoal pb-4">
  <div>
    <h2 class="font-display font-extrabold text-2xl md:text-3xl uppercase tracking-tight flex items-center gap-2">
      <ShoppingCart class="w-7 h-7" />
      <span>POS Retail Checkout</span>
    </h2>
    <p class="font-sans text-xs sm:text-sm text-zinc-500 font-semibold mt-1">
      Direct checkout for accessories, screen guards, chargers, and general retail sales.
    </p>
  </div>

  <div>
    <a href="/admin/sales/archive" class="inline-block">
      <Button bgColor="bg-white" class="flex items-center gap-1.5 py-2 px-4 text-xs font-bold shadow-neubrutalism-sm border-2 border-neubrutalism-charcoal">
        <FileText class="w-4 h-4" />
        <span>VIEW SALES TRANSACTIONS</span>
      </Button>
    </a>
  </div>
</div>
