<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    type?: 'button' | 'submit' | 'reset';
    disabled?: boolean;
    bgColor?: string;
    onclick?: (event: MouseEvent) => void;
    class?: string;
    children: Snippet;
  }

  let {
    type = 'button',
    disabled = false,
    bgColor = 'bg-neubrutalism-yellow',
    onclick,
    class: className = '',
    children
  }: Props = $props();
</script>

<button
  {type}
  {disabled}
  {onclick}
  class="font-display font-bold text-lg px-6 py-3 border-4 border-neubrutalism-charcoal rounded-none text-neubrutalism-charcoal transition-all duration-150 select-none cursor-pointer focus:outline-none focus:ring-4 focus:ring-neubrutalism-charcoal
  {bgColor}
  {disabled ? 'opacity-50 cursor-not-allowed shadow-none translate-x-0 translate-y-0' : 'shadow-neubrutalism-md hover:-translate-x-1 hover:-translate-y-1 hover:shadow-neubrutalism-lg active:translate-x-0 active:translate-y-0 active:shadow-neubrutalism-md'}
  {className}"
>
  {@render children()}
</button>
