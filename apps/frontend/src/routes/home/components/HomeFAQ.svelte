<script lang="ts">
  let activeFaqIndex = $state<number | null>(null);

  const faqItems = [
    {
      q: "Do I need to book an appointment before coming?",
      a: "No appointment is needed! You can walk directly into our workshop in Jakarta Pusat. We diagnose all devices on a first-come, first-served basis.",
    },
    {
      q: "How long does a typical repair take?",
      a: "Basic screen or battery replacements usually take 1 to 2 hours. More complex diagnostics (like motherboard troubleshooting or liquid damage checks) may take 24 to 48 hours.",
    },
    {
      q: "Do you offer any warranty on repairs?",
      a: "Yes, absolutely! We provide a digital warranty (typically 30 to 90 days depending on the replacement part used) which is tracked online via your unique ticket ID.",
    },
    {
      q: "Is my personal data safe during repair?",
      a: "We prioritize your privacy. Technicians do not require device passcodes unless testing specific functions. We always advise backing up your data before drop-off if possible.",
    },
  ];
</script>

<section class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-16 w-full">
  <div class="text-center flex flex-col gap-3 mb-12">
    <h2
      class="font-display font-black text-3xl sm:text-5xl uppercase tracking-tight"
    >
      FREQUENTLY ASKED QUESTIONS
    </h2>
    <p
      class="font-sans text-sm sm:text-base text-zinc-500 font-semibold uppercase tracking-wider"
    >
      Answers to common concerns regarding repairs, billing, and diagnostics.
    </p>
  </div>

  <div class="flex flex-col gap-4">
    {#each faqItems as item, idx}
      <div
        class="border-4 border-neubrutalism-charcoal bg-white shadow-neubrutalism-sm transition-all overflow-hidden"
      >
        <button
          onclick={() => (activeFaqIndex = activeFaqIndex === idx ? null : idx)}
          class="w-full flex justify-between items-center p-5 font-display font-black text-base uppercase text-left hover:bg-zinc-50 transition-colors"
        >
          <span>{item.q}</span>
          <span
            class="bg-neubrutalism-yellow border-2 border-neubrutalism-charcoal px-3 py-1 font-mono text-xs font-black shadow-neubrutalism-sm"
          >
            {activeFaqIndex === idx ? "CLOSE" : "OPEN"}
          </span>
        </button>

        {#if activeFaqIndex === idx}
          <div
            class="p-5 border-t-2 border-dashed border-zinc-200 bg-zinc-50 font-sans text-sm text-zinc-600 font-semibold leading-relaxed"
          >
            {item.a}
          </div>
        {/if}
      </div>
    {/each}
  </div>
</section>
