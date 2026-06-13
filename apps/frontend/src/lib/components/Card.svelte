<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    class?: string;
    bgColor?: string;
    children: Snippet;
  }

  let { class: className = '', bgColor = 'bg-white', children }: Props = $props();
</script>

<div class="border-4 border-neubrutalism-charcoal {bgColor} shadow-neubrutalism-md p-6 md:p-8 rounded-none transition-all duration-150 {className}">
  {@render children()}
</div>
