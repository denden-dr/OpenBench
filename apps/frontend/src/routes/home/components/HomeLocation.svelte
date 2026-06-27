<script lang="ts">
  import { Card } from "$lib";
  import { MapPin, Clock, Phone, Mail } from "lucide-svelte";
</script>

<section class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16 w-full">
  <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-stretch">
    <!-- Left: Location Card -->
    <div class="lg:col-span-7 flex">
      <Card
        bgColor="bg-white"
        class="border-4 border-neubrutalism-charcoal shadow-neubrutalism-lg p-6 flex flex-col gap-6 w-full"
      >
        <div
          class="flex items-center gap-3 border-b-2 border-dashed border-zinc-200 pb-4"
        >
          <div
            class="bg-neubrutalism-green border-2 border-neubrutalism-charcoal p-2.5"
          >
            <MapPin class="w-6 h-6 text-neubrutalism-charcoal" />
          </div>
          <div>
            <h3
              class="font-display font-black text-xl uppercase tracking-tight"
            >
              VISIT OUR WORKSHOP
            </h3>
            <p
              class="font-sans text-xs text-zinc-500 font-bold uppercase tracking-wider"
            >
              WALK-IN SERVICE & DEVICE DROP-OFF
            </p>
          </div>
        </div>

        <div class="flex flex-col gap-4">
          <p
            class="font-sans text-base font-semibold leading-relaxed text-zinc-700"
          >
            Since online booking is temporarily offline for shipping integration
            updates, you can bring your gadget directly to our main service
            bench.
          </p>

          <div
            class="bg-zinc-50 border-4 border-neubrutalism-charcoal p-5 shadow-neubrutalism-sm"
          >
            <h4
              class="font-display font-extrabold text-sm uppercase text-zinc-500 mb-1"
            >
              ADDRESS
            </h4>
            <p
              class="font-sans text-base font-bold text-neubrutalism-charcoal leading-relaxed"
            >
              Jakarta Pusat, DKI Jakarta
            </p>
          </div>
        </div>
      </Card>
    </div>

    <!-- Right: Hours & Contact Card -->
    <div class="lg:col-span-5 flex">
      <Card
        bgColor="bg-neubrutalism-yellow"
        class="border-4 border-neubrutalism-charcoal shadow-neubrutalism-lg p-6 flex flex-col gap-6 w-full text-neubrutalism-charcoal"
      >
        <div
          class="flex items-center gap-3 border-b-2 border-dashed border-neubrutalism-charcoal/30 pb-4"
        >
          <div class="bg-white border-2 border-neubrutalism-charcoal p-2.5">
            <Clock class="w-6 h-6 text-neubrutalism-charcoal" />
          </div>
          <div>
            <h3
              class="font-display font-black text-xl uppercase tracking-tight"
            >
              HOURS & CONTACT
            </h3>
            <p
              class="font-sans text-xs text-zinc-700 font-bold uppercase tracking-wider"
            >
              WHEN WE ARE BENCH-READY
            </p>
          </div>
        </div>

        <div class="flex flex-col gap-4 font-mono text-sm">
          <div
            class="flex justify-between items-center border-b border-neubrutalism-charcoal/20 pb-2"
          >
            <span class="font-bold">MON - FRI</span>
            <span class="font-extrabold text-neubrutalism-charcoal"
              >09:00 - 18:00</span
            >
          </div>
          <div
            class="flex justify-between items-center border-b border-neubrutalism-charcoal/20 pb-2"
          >
            <span class="font-bold">SATURDAY</span>
            <span class="font-extrabold text-neubrutalism-charcoal"
              >09:00 - 15:00</span
            >
          </div>
          <div class="flex justify-between items-center pb-2">
            <span class="font-bold">SUNDAY</span>
            <span class="font-extrabold text-zinc-600">CLOSED</span>
          </div>
        </div>

        <div
          class="mt-auto border-t-2 border-dashed border-neubrutalism-charcoal/30 pt-6 flex flex-col gap-3"
        >
          <div class="flex items-center gap-3">
            <Phone class="w-5 h-5" />
            <span class="font-sans font-bold text-sm">+62 812-3456-7890</span>
          </div>
          <div class="flex items-center gap-3">
            <Mail class="w-5 h-5" />
            <span class="font-sans font-bold text-sm"
              >support@openbench.dev</span
            >
          </div>
        </div>
      </Card>
    </div>
  </div>
</section>
