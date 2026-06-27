<script lang="ts">
  import { Card } from "$lib";
  import { Search, Zap, ShieldCheck } from "lucide-svelte";
</script>

<section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16 sm:py-24 w-full">
  <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
    <Card bgColor="bg-white" class="p-6 flex flex-col gap-3">
      <div class="bg-neubrutalism-pink border-2 border-neubrutalism-charcoal p-2 w-fit text-white">
        <Search class="w-5 h-5 text-white" />
      </div>
      <h3 class="font-display font-bold text-lg uppercase">
        Zero Hidden Fees
      </h3>
      <p class="font-sans text-sm text-zinc-500 font-semibold leading-relaxed">
        Get complete breakdown lists of labor costs and replacement
        components. You approve the diagnostic report before we repair.
      </p>
    </Card>

    <Card bgColor="bg-white" class="p-6 flex flex-col gap-3">
      <div class="bg-neubrutalism-green border-2 border-neubrutalism-charcoal p-2 w-fit">
        <Zap class="w-5 h-5 text-neubrutalism-charcoal" />
      </div>
      <h3 class="font-display font-bold text-lg uppercase">
        Instant Tracker
      </h3>
      <p class="font-sans text-sm text-zinc-500 font-semibold leading-relaxed">
        Every device has a unique live-tracker link. Check status logs, read
        comments from the technicians, and download receipts in real time.
      </p>
    </Card>

    <Card bgColor="bg-white" class="p-6 flex flex-col gap-3">
      <div class="bg-neubrutalism-yellow border-2 border-neubrutalism-charcoal p-2 w-fit">
        <ShieldCheck class="w-5 h-5 text-neubrutalism-charcoal" />
      </div>
      <h3 class="font-display font-bold text-lg uppercase">
        Warranty Protection
      </h3>
      <p class="font-sans text-sm text-zinc-500 font-semibold leading-relaxed">
        Every eligible repair includes a digital warranty tracking system. We
        stand behind our workmanship.
      </p>
    </Card>
  </div>
</section>
