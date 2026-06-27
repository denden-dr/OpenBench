<script lang="ts">
  import { Wrench, ShieldCheck, ClipboardList, Search } from "lucide-svelte";
</script>

<section class="border-t-4 border-b-4 border-neubrutalism-charcoal bg-white py-16">
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
    <!-- Section Header -->
    <div class="text-center flex flex-col gap-3 mb-16">
      <h2 class="font-display font-black text-3xl sm:text-5xl uppercase tracking-tight">
        HOW WE ORCHESTRATE REPAIRS
      </h2>
      <p class="font-sans text-sm sm:text-base text-zinc-500 font-semibold uppercase tracking-wider">
        Four structured steps to get your device back to peak performance.
      </p>
    </div>

    <!-- Stepper Grid -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-8 relative">
      <!-- Step 1 -->
      <div class="flex flex-col gap-4 relative">
        <div class="w-12 h-12 bg-neubrutalism-pink border-4 border-neubrutalism-charcoal font-display font-black text-xl flex items-center justify-center shadow-neubrutalism-sm text-white select-none">
          1
        </div>
        <h3 class="font-display font-bold text-lg uppercase tracking-tight mt-2 flex items-center gap-2">
          <ClipboardList class="w-5 h-5 text-neubrutalism-pink" />
          <span>Submit Request</span>
        </h3>
        <p class="font-sans text-sm text-zinc-600 leading-relaxed font-medium">
          Fill out our simplified repair form detailing your device
          specifications and problem symptoms.
        </p>
      </div>

      <!-- Step 2 -->
      <div class="flex flex-col gap-4 relative">
        <div class="w-12 h-12 bg-neubrutalism-yellow border-4 border-neubrutalism-charcoal font-display font-black text-xl flex items-center justify-center shadow-neubrutalism-sm select-none">
          2
        </div>
        <h3 class="font-display font-bold text-lg uppercase tracking-tight mt-2 flex items-center gap-2">
          <Search class="w-5 h-5 text-neubrutalism-charcoal" />
          <span>Diagnosis & Quote</span>
        </h3>
        <p class="font-sans text-sm text-zinc-600 leading-relaxed font-medium">
          Bring your device in. Our technicians run diagnostics and submit a
          proposed repair action and cost.
        </p>
      </div>

      <!-- Step 3 -->
      <div class="flex flex-col gap-4 relative">
        <div class="w-12 h-12 bg-neubrutalism-green border-4 border-neubrutalism-charcoal font-display font-black text-xl flex items-center justify-center shadow-neubrutalism-sm select-none">
          3
        </div>
        <h3 class="font-display font-bold text-lg uppercase tracking-tight mt-2 flex items-center gap-2">
          <Wrench class="w-5 h-5 text-neubrutalism-charcoal" />
          <span>Watch Live Repair</span>
        </h3>
        <p class="font-sans text-sm text-zinc-600 leading-relaxed font-medium">
          Monitor real-time updates and progress transitions via the tracker
          console using your unique UUID.
        </p>
      </div>

      <!-- Step 4 -->
      <div class="flex flex-col gap-4 relative">
        <div class="w-12 h-12 bg-zinc-200 border-4 border-neubrutalism-charcoal font-display font-black text-xl flex items-center justify-center shadow-neubrutalism-sm select-none">
          4
        </div>
        <h3 class="font-display font-bold text-lg uppercase tracking-tight mt-2 flex items-center gap-2">
          <ShieldCheck class="w-5 h-5 text-zinc-700" />
          <span>Secure Pickup</span>
        </h3>
        <p class="font-sans text-sm text-zinc-600 leading-relaxed font-medium">
          Settle payment securely, retrieve your fully serviced device, and
          receive your digital warranty parameters.
        </p>
      </div>
    </div>
  </div>
</section>
