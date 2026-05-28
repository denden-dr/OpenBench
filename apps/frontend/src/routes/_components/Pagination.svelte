<script lang="ts">
  import { ChevronLeft, ChevronRight } from "lucide-svelte";

  // Svelte 5 Props
  let {
    currentPage = $bindable(1),
    totalPages = 1,
    totalItems = 0,
    limit = 20,
    onPageChange = (page: number) => {}
  }: {
    currentPage: number;
    totalPages: number;
    totalItems: number;
    limit: number;
    onPageChange?: (page: number) => void;
  } = $props();

  // Hitung rentang data saat ini
  let startRange = $derived(totalItems === 0 ? 0 : (currentPage - 1) * limit + 1);
  let endRange = $derived(Math.min(currentPage * limit, totalItems));
  let visiblePages = $derived.by(() => {
    const pages = new Set<number>();
    pages.add(1);
    pages.add(totalPages);
    for (let page = currentPage - 1; page <= currentPage + 1; page += 1) {
      if (page >= 1 && page <= totalPages) pages.add(page);
    }
    return Array.from(pages).sort((a, b) => a - b);
  });

  function goToPage(page: number) {
    if (page >= 1 && page <= totalPages && page !== currentPage) {
      currentPage = page;
      onPageChange(page);
    }
  }
</script>

<div class="flex flex-col sm:flex-row justify-between items-center gap-4 mt-6 px-2 py-4 border-t border-slate-100 dark:border-slate-900">
  <!-- Info Status Rentang -->
  <div class="text-xs text-slate-500 dark:text-slate-400 font-medium">
    Menampilkan <span class="font-semibold text-slate-800 dark:text-slate-200">{startRange}-{endRange}</span> 
    dari <span class="font-semibold text-slate-800 dark:text-slate-200">{totalItems}</span> data
  </div>

  <!-- Tombol Kontrol Halaman -->
  <div class="flex items-center gap-1.5">
    <!-- Prev Button -->
    <button
      onclick={() => goToPage(currentPage - 1)}
      disabled={currentPage <= 1}
      class="inline-flex items-center justify-center p-2 rounded-xl border border-slate-200 dark:border-slate-800 text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-900 transition-all duration-200 disabled:opacity-45 disabled:pointer-events-none active:scale-95 cursor-pointer"
      aria-label="Halaman sebelumnya"
    >
      <ChevronLeft size={16} />
    </button>

    <!-- Halaman Utama: render hanya halaman yang terlihat, bukan semua totalPages -->
    {#each visiblePages as pageNum, index}
      {#if index > 0 && pageNum - visiblePages[index - 1] > 1}
        <span class="text-slate-400 px-1 text-xs select-none">...</span>
      {/if}

      <button
        onclick={() => goToPage(pageNum)}
        class="min-w-[34px] h-[34px] flex items-center justify-center rounded-xl text-xs font-semibold border transition-all duration-200 active:scale-95 cursor-pointer
          {currentPage === pageNum
            ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white border-blue-600 shadow-sm shadow-blue-500/20 dark:shadow-indigo-500/10'
            : 'border-slate-200 dark:border-slate-800 text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-900 bg-white dark:bg-slate-950'}"
      >
        {pageNum}
      </button>
    {/each}

    <!-- Next Button -->
    <button
      onclick={() => goToPage(currentPage + 1)}
      disabled={currentPage >= totalPages}
      class="inline-flex items-center justify-center p-2 rounded-xl border border-slate-200 dark:border-slate-800 text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-900 transition-all duration-200 disabled:opacity-45 disabled:pointer-events-none active:scale-95 cursor-pointer"
      aria-label="Halaman berikutnya"
    >
      <ChevronRight size={16} />
    </button>
  </div>
</div>
