<script lang="ts">
  import { Card, Button } from "$lib";
  import { ArrowRight, Search, Zap } from "lucide-svelte";
</script>

<section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16 md:py-24 w-full flex-grow flex flex-col justify-center gap-12">
  <div class="grid grid-cols-1 lg:grid-cols-12 gap-12 items-center">
    <!-- Left: Catchy Copy & CTAs -->
    <div class="lg:col-span-7 flex flex-col gap-6 text-left">
      <!-- Badge Tag -->
      <span class="w-fit font-mono text-xs font-bold bg-neubrutalism-green border-2 border-neubrutalism-charcoal px-3 py-1 shadow-neubrutalism-sm uppercase tracking-wider transform -rotate-1">
        ⚙️ REPAIR SHOP ORCHESTRATION ENGINE
      </span>

      <h1 class="font-display font-black text-4xl sm:text-6xl uppercase tracking-tight leading-none">
        FAST, TRANSPARENT <br />
        <span class="bg-neubrutalism-yellow border-4 border-neubrutalism-charcoal px-4 py-1.5 shadow-neubrutalism-md inline-block transform rotate-1 mt-2">
          GADGET REPAIRS
        </span>
      </h1>

      <p class="font-sans text-base sm:text-lg text-zinc-600 font-semibold max-w-xl leading-relaxed mt-2">
        No guesswork. No hidden fees. Submit a repair request, follow our live
        progress logs, and only approve charges after diagnosis.
      </p>

      <div class="flex flex-col sm:flex-row gap-4 mt-4">
        <a href="/repair-request" class="w-full sm:w-auto">
          <Button bgColor="bg-neubrutalism-pink" class="w-full sm:w-auto flex items-center justify-center gap-2 group py-3 px-8 text-white">
            <span>REQUEST REPAIR</span>
            <ArrowRight class="w-5 h-5 group-hover:translate-x-1 transition-transform text-white" />
          </Button>
        </a>

        <a href="/tracker" class="w-full sm:w-auto">
          <Button bgColor="bg-white" class="w-full sm:w-auto flex items-center justify-center gap-2 py-3 px-8 border-2 border-neubrutalism-charcoal shadow-neubrutalism-sm hover:bg-zinc-50">
            <Search class="w-5 h-5" />
            <span>TRACK STATUS</span>
          </Button>
        </a>
      </div>
    </div>

    <!-- Right: Neubrutalist Graphic / Stat Cards -->
    <div class="lg:col-span-5 flex flex-col gap-6 relative">
      <Card bgColor="bg-white" class="border-4 border-neubrutalism-charcoal shadow-neubrutalism-lg relative overflow-hidden">
        <div class="absolute top-0 right-0 bg-neubrutalism-yellow border-b-4 border-l-4 border-neubrutalism-charcoal font-mono text-xs font-bold px-3 py-1 uppercase">
          LIVE STATS
        </div>

        <h3 class="font-display font-bold text-xl uppercase tracking-tight mb-4 flex items-center gap-2">
          <Zap class="w-5 h-5 text-neubrutalism-yellow fill-neubrutalism-yellow" />
          <span>ORCHESTRATION STATUS</span>
        </h3>

        <div class="flex flex-col gap-4 font-mono text-xs font-semibold">
          <div class="flex justify-between items-center border-b-2 border-dashed border-zinc-200 pb-2">
            <span>DIAGNOSTICS PENDING</span>
            <span class="bg-amber-100 text-amber-800 border border-amber-300 px-2 py-0.5">READY</span>
          </div>
          <div class="flex justify-between items-center border-b-2 border-dashed border-zinc-200 pb-2">
            <span>ACTIVE REPAIRS</span>
            <span class="bg-blue-100 text-blue-800 border border-blue-300 px-2 py-0.5 animate-pulse">ON-BENCH</span>
          </div>
          <div class="flex justify-between items-center border-b-2 border-dashed border-zinc-200 pb-2">
            <span>COMPLETED TODAY</span>
            <span class="bg-emerald-100 text-emerald-800 border border-emerald-300 px-2 py-0.5">PASSED</span>
          </div>
          <div class="flex justify-between items-center">
            <span>AVERAGE TURNAROUND</span>
            <span class="bg-zinc-100 border border-zinc-300 px-2 py-0.5">24 HOURS</span>
          </div>
        </div>
      </Card>

      <!-- Decorative floating badge -->
      <div class="absolute -bottom-4 -left-4 bg-neubrutalism-pink text-white font-display font-black text-sm uppercase px-4 py-2 border-4 border-neubrutalism-charcoal shadow-neubrutalism-sm transform -rotate-3 select-none">
        100% TRANSPARENT
      </div>
    </div>
  </div>
</section>
