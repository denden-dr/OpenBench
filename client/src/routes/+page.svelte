<script lang="ts">
	import { fetchHealth, type HealthData } from '$lib/api';

	let status = $state<'loading' | 'connected' | 'error'>('loading');
	let healthData = $state<HealthData | null>(null);
	let errorMessage = $state<string | null>(null);

	$effect(() => {
		fetchHealth()
			.then((res) => {
				healthData = res.data;
				status = 'connected';
			})
			.catch((err) => {
				errorMessage = err.message;
				status = 'error';
			});
	});
</script>

<div class="min-h-screen bg-neutral-950 text-neutral-100 flex items-center justify-center p-6 font-sans">
	<div class="max-w-md w-full bg-neutral-900/50 backdrop-blur-xl border border-neutral-800 rounded-3xl p-8 shadow-2xl overflow-hidden relative group">
		<!-- Background Glow -->
		<div class="absolute -inset-24 bg-gradient-to-br from-indigo-500/10 to-purple-500/10 opacity-50 blur-3xl group-hover:opacity-100 transition duration-1000"></div>

		<div class="relative z-10">
			<header class="flex items-center justify-between mb-8">
				<h1 class="text-2xl font-bold tracking-tight bg-gradient-to-r from-indigo-400 to-purple-400 bg-clip-text text-transparent">
					OpenBench
				</h1>
				<div class="flex items-center gap-2">
					{#if status === 'loading'}
						<div class="w-2 h-2 rounded-full bg-neutral-400 animate-pulse"></div>
						<span class="text-xs text-neutral-400 font-medium">Checking...</span>
					{:else if status === 'connected'}
						<div class="w-2 h-2 rounded-full bg-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.5)]"></div>
						<span class="text-xs text-emerald-400 font-medium tracking-wide uppercase">Backend Online</span>
					{:else}
						<div class="w-2 h-2 rounded-full bg-rose-500"></div>
						<span class="text-xs text-rose-400 font-medium uppercase tracking-wide">Connection Error</span>
					{/if}
				</div>
			</header>

			<main class="space-y-6">
				<p class="text-neutral-400 leading-relaxed text-sm">
					System initialization complete. Connected to high-performance Fiber v3 backend.
				</p>

				<div class="bg-black/40 rounded-2xl p-6 border border-neutral-800/50">
					<div class="flex justify-between items-center text-sm">
						<span class="text-neutral-500">Platform Version</span>
						<span class="font-mono text-indigo-400">
							{healthData?.version || 'v0.0.0'}
						</span>
					</div>
				</div>

				{#if status === 'error'}
					<div class="bg-rose-500/10 border border-rose-500/20 rounded-xl p-4 text-xs text-rose-400 text-center">
						{errorMessage}
					</div>
				{/if}
			</main>

			<footer class="mt-8 pt-8 border-t border-neutral-800/50">
				<button class="w-full bg-white text-black font-semibold py-3 px-6 rounded-xl hover:bg-neutral-200 transition duration-300 active:scale-95 text-sm">
					Open Dashboard
				</button>
			</footer>
		</div>
	</div>
</div>

<style>
	:global(body) {
		margin: 0;
		background-color: #0a0a0a;
	}
</style>
