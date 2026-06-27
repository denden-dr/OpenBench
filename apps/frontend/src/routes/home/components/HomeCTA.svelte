<script lang="ts">
  import { Button } from "$lib";
</script>

<section class="bg-neubrutalism-yellow border-t-4 border-neubrutalism-charcoal py-12 text-center">
  <div class="max-w-4xl mx-auto px-4 flex flex-col md:flex-row items-center justify-between gap-6">
    <div class="text-left flex flex-col gap-1.5">
      <h2 class="font-display font-black text-2xl uppercase tracking-tight">
        READY TO RESTORE YOUR DEVICE?
      </h2>
      <p class="font-sans text-xs sm:text-sm text-zinc-800 font-bold uppercase">
        Sign up or sign in to submit repair requests and trace diagnostic
        logs.
      </p>
    </div>

    <a href="/repair-request" class="w-full md:w-auto">
      <Button bgColor="bg-neubrutalism-pink" class="w-full md:w-auto py-3 px-8 text-sm text-white font-display font-bold shadow-neubrutalism-sm">
        START REPAIR REQUEST
      </Button>
    </a>
  </div>
</section>
