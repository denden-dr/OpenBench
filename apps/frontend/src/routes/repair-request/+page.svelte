<script lang="ts">
  import { onMount } from "svelte";
  import { Card, Button } from "$lib";
  import { ArrowLeft, Clock, MapPin, Truck, Wrench } from "lucide-svelte";
  import { authService } from "$lib/services/auth";

  let userSession = $state<any | null>(null);
  let isChecking = $state(true);

  onMount(async () => {
    userSession = await authService.checkSession();
    isChecking = false;
  });
</script>

<svelte:head>
  <title>Online Booking Coming Soon - OpenBench</title>
</svelte:head>

<div
  class="min-h-screen bg-neubrutalism-bg flex flex-col font-sans p-4 md:p-8 text-neubrutalism-charcoal"
>
  <div class="max-w-2xl w-full mx-auto flex flex-col gap-6">
    <!-- Header/Navigation Back -->
    <div class="flex justify-between items-center">
      {#if isChecking}
        <span class="w-24 h-4 bg-zinc-200 animate-pulse"></span>
      {:else}
        <a
          href={userSession ? "/portal" : "/home"}
          class="inline-flex items-center gap-1.5 font-mono text-xs font-bold uppercase hover:underline"
        >
          <ArrowLeft class="w-4 h-4" />
          <span>{userSession ? "BACK TO PORTAL" : "BACK TO HOME"}</span>
        </a>
      {/if}
      <span
        class="font-mono text-xs font-bold bg-neubrutalism-charcoal text-white px-2 py-0.5 shadow-neubrutalism-sm"
      >
        FEATURE UPDATE
      </span>
    </div>

    <!-- Main coming soon card -->
    <Card
      bgColor="bg-white"
      class="border-4 border-neubrutalism-charcoal p-8 shadow-neubrutalism-lg flex flex-col items-center text-center gap-6"
    >
      <!-- Feature Icon badge -->
      <div
        class="bg-neubrutalism-yellow border-4 border-neubrutalism-charcoal p-4 rounded-none shadow-neubrutalism-md w-fit"
      >
        <Truck class="w-12 h-12 text-neubrutalism-charcoal" />
      </div>

      <div class="flex flex-col gap-3">
        <span
          class="w-fit mx-auto font-mono text-xs font-black bg-neubrutalism-pink text-white border-2 border-neubrutalism-charcoal px-3 py-1 shadow-neubrutalism-sm uppercase tracking-wider transform -rotate-1"
        >
          🚧 COMING SOON
        </span>
        <h1
          class="font-display font-black text-3xl sm:text-4xl uppercase tracking-tight leading-none mt-2"
        >
          ONLINE BOOKING & <br />
          <span
            class="bg-neubrutalism-green border-4 border-neubrutalism-charcoal px-3 py-1 shadow-neubrutalism-sm inline-block transform rotate-1 mt-1 text-neubrutalism-charcoal"
          >
            MAIL-IN REPAIR
          </span>
        </h1>
        <p
          class="font-sans text-sm text-zinc-600 font-semibold max-w-md mt-4 leading-relaxed"
        >
          We are currently building our online ticketing flow in tandem with a
          complete shipping and courier delivery integration system. Soon, you
          will be able to book repairs online and have your device picked up and
          returned via courier!
        </p>
      </div>

      <!-- Drop off details banner -->
      <div
        class="w-full border-t-4 border-b-4 border-dashed border-neubrutalism-charcoal/30 py-6 my-2 text-left flex flex-col gap-4"
      >
        <h3
          class="font-display font-black text-lg uppercase tracking-tight flex items-center gap-2"
        >
          <Wrench class="w-5 h-5 text-neubrutalism-pink" />
          <span>IN THE MEANTIME: WALK-IN SERVICE</span>
        </h3>
        <p
          class="font-sans text-xs text-zinc-500 font-bold uppercase tracking-wider"
        >
          You don't have to wait! Bring your device directly to our physical
          workshop for diagnostics:
        </p>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-2">
          <div
            class="bg-zinc-50 border-4 border-neubrutalism-charcoal p-4 shadow-neubrutalism-sm flex flex-col gap-2"
          >
            <span
              class="flex items-center gap-1.5 font-display font-extrabold text-xs text-zinc-500 uppercase"
            >
              <MapPin class="w-4 h-4 text-neubrutalism-green" />
              <span>OUR WORKSHOP</span>
            </span>
            <p
              class="font-sans text-xs font-bold leading-normal text-neubrutalism-charcoal"
            >
              Jakarta Pusat, DKI Jakarta
            </p>
          </div>

          <div
            class="bg-zinc-50 border-4 border-neubrutalism-charcoal p-4 shadow-neubrutalism-sm flex flex-col gap-2"
          >
            <span
              class="flex items-center gap-1.5 font-display font-extrabold text-xs text-zinc-500 uppercase"
            >
              <Clock class="w-4 h-4 text-neubrutalism-yellow" />
              <span>BENCH HOURS</span>
            </span>
            <p
              class="font-sans text-xs font-bold leading-normal text-neubrutalism-charcoal"
            >
              Mon - Fri: 09:00 - 18:00<br />
              Saturday: 09:00 - 15:00<br />
              Sunday: CLOSED
            </p>
          </div>
        </div>
      </div>

      <!-- Back Action -->
      <a href={userSession ? "/portal" : "/home"} class="w-full sm:w-auto">
        <Button
          bgColor="bg-neubrutalism-pink"
          class="w-full sm:w-auto py-3 px-8 text-sm text-white font-display font-bold shadow-neubrutalism-sm"
        >
          {userSession ? "RETURN TO PORTAL" : "GO TO HOME PAGE"}
        </Button>
      </a>
    </Card>
  </div>
</div>
