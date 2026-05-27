<script lang="ts">
  let {
    selectedStatus = $bindable(),
    statusCounts = {},
  }: {
    selectedStatus: string;
    statusCounts: Record<string, number>;
  } = $props();

  const statuses = [
    {
      value: "service_in",
      label: "Masuk",
      activeClass: "bg-blue-600 text-white border-blue-600 dark:bg-blue-600 dark:text-white dark:border-blue-600 ring-2 ring-offset-2 ring-blue-500/20 dark:ring-offset-slate-950",
      inactiveClass: "bg-blue-50/50 hover:bg-blue-100/50 dark:bg-blue-950/20 dark:hover:bg-blue-950/40 text-blue-700 dark:text-blue-400 border-blue-200/50 dark:border-blue-900/50",
      badgeActiveClass: "bg-white/20 text-white",
      badgeInactiveClass: "bg-blue-100/80 dark:bg-blue-900/60 text-blue-800 dark:text-blue-300"
    },
    {
      value: "on_process",
      label: "Sedang Diproses",
      activeClass: "bg-amber-500 text-white border-amber-500 dark:bg-amber-500 dark:text-white dark:border-amber-500 ring-2 ring-offset-2 ring-amber-500/20 dark:ring-offset-slate-950",
      inactiveClass: "bg-amber-50/50 hover:bg-amber-100/50 dark:bg-amber-950/20 dark:hover:bg-amber-950/40 text-amber-700 dark:text-amber-400 border-amber-200/50 dark:border-amber-900/50",
      badgeActiveClass: "bg-white/20 text-white",
      badgeInactiveClass: "bg-amber-100/80 dark:bg-amber-900/60 text-amber-800 dark:text-amber-300"
    },
    {
      value: "waiting_confirmation",
      label: "Menunggu Konfirmasi",
      activeClass: "bg-rose-500 text-white border-rose-500 dark:bg-rose-500 dark:text-white dark:border-rose-500 ring-2 ring-offset-2 ring-rose-500/20 dark:ring-offset-slate-950",
      inactiveClass: "bg-rose-50/50 hover:bg-rose-100/50 dark:bg-rose-950/20 dark:hover:bg-rose-950/40 text-rose-700 dark:text-rose-400 border-rose-200/50 dark:border-rose-900/50",
      badgeActiveClass: "bg-white/20 text-white",
      badgeInactiveClass: "bg-rose-100/80 dark:bg-rose-900/60 text-rose-800 dark:text-rose-300"
    },
    {
      value: "fixed",
      label: "Siap Diambil",
      activeClass: "bg-emerald-600 text-white border-emerald-600 dark:bg-emerald-600 dark:text-white dark:border-emerald-600 ring-2 ring-offset-2 ring-emerald-500/20 dark:ring-offset-slate-950",
      inactiveClass: "bg-emerald-50/50 hover:bg-emerald-100/50 dark:bg-emerald-950/20 dark:hover:bg-emerald-950/40 text-emerald-700 dark:text-emerald-400 border-emerald-200/50 dark:border-emerald-900/50",
      badgeActiveClass: "bg-white/20 text-white",
      badgeInactiveClass: "bg-emerald-100/80 dark:bg-emerald-900/60 text-emerald-800 dark:text-emerald-300"
    },
    {
      value: "cancelled",
      label: "Dibatalkan",
      activeClass: "bg-slate-600 text-white border-slate-600 dark:bg-slate-600 dark:text-white dark:border-slate-600 ring-2 ring-offset-2 ring-slate-500/20 dark:ring-offset-slate-950",
      inactiveClass: "bg-slate-100 hover:bg-slate-200/60 dark:bg-slate-900/50 dark:hover:bg-slate-900/80 text-slate-700 dark:text-slate-300 border-slate-200 dark:border-slate-800",
      badgeActiveClass: "bg-white/20 text-white",
      badgeInactiveClass: "bg-slate-200 dark:bg-slate-800 text-slate-800 dark:text-slate-200"
    }
  ];
</script>

<div class="flex items-center gap-2 overflow-x-auto scrollbar-none py-1 -mx-4 px-4 md:mx-0 md:px-0 md:flex-wrap">
  {#each statuses as status}
    {@const count = statusCounts[status.value] ?? 0}
    {#if count > 0}
      {@const isActive = selectedStatus === status.value}
      <button
        onclick={() => {
          if (isActive) {
            selectedStatus = "all";
          } else {
            selectedStatus = status.value;
          }
        }}
        class="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-full text-xs font-semibold border transition-all duration-200 cursor-pointer select-none active:scale-95 whitespace-nowrap
          {isActive ? status.activeClass : status.inactiveClass}"
      >
        <span>{status.label}</span>
        <span class="inline-flex items-center justify-center px-1.5 py-0.5 rounded-full text-[10px] font-bold transition-colors duration-200
          {isActive ? status.badgeActiveClass : status.badgeInactiveClass}">
          {count}
        </span>
      </button>
    {/if}
  {/each}
</div>

<style>
  /* Hide scrollbar for Chrome, Safari and Opera */
  .scrollbar-none::-webkit-scrollbar {
    display: none;
  }
  /* Hide scrollbar for IE, Edge and Firefox */
  .scrollbar-none {
    -ms-overflow-style: none;  /* IE and Edge */
    scrollbar-width: none;  /* Firefox */
  }
</style>
