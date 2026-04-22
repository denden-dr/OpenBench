<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>OpenBench - Your Device, Unveiled.</title>
</svelte:head>

<div class="min-h-screen flex flex-col">
	<!-- TopNavBar -->
	<nav class="fixed top-0 w-full z-50 bg-white/90 backdrop-blur-md border-b border-slate-200">
		<div class="max-w-7xl mx-auto px-8 h-16 flex items-center justify-between">
			<div class="text-xl font-black text-black tracking-tighter">
				OpenBench
			</div>
			<div class="hidden md:flex items-center gap-8">
				<a class="font-['Inter'] uppercase tracking-wider text-xs font-bold text-slate-600 hover:text-slate-900 transition-colors hover:bg-slate-50 py-2 px-3 rounded" href="#">Documentation</a>
				<a class="font-['Inter'] uppercase tracking-wider text-xs font-bold text-slate-600 hover:text-slate-900 transition-colors hover:bg-slate-50 py-2 px-3 rounded" href="#">Services</a>
				<a class="font-['Inter'] uppercase tracking-wider text-xs font-bold text-slate-600 hover:text-slate-900 transition-colors hover:bg-slate-50 py-2 px-3 rounded" href="#">Repair Kits</a>
				<a class="font-['Inter'] uppercase tracking-wider text-xs font-bold text-slate-600 hover:text-slate-900 transition-colors hover:bg-slate-50 py-2 px-3 rounded" href="#">Community</a>
			</div>
			<div class="flex items-center gap-4">
				<div class="hidden lg:flex items-center text-slate-400">
					<span class="material-symbols-outlined text-xl">search</span>
				</div>
				<a class="font-['Inter'] uppercase tracking-wider text-xs font-bold text-orange-600 hover:bg-slate-50 transition-all px-4 py-2 rounded" href="#">
					Book a Repair
				</a>
			</div>
		</div>
	</nav>

	<main class="flex-grow">
		{@render children()}
	</main>

	<!-- Footer -->
	<footer class="w-full py-20 bg-slate-50 border-t border-slate-200">
		<div class="max-w-7xl mx-auto px-8 grid grid-cols-1 md:grid-cols-2 gap-8">
			<div class="flex flex-col gap-4">
				<div class="font-mono font-bold text-slate-900 text-xl">
					OpenBench
				</div>
				<p class="font-['Inter'] text-sm text-slate-500">
					© 2026 OpenBench. Built for the Right to Repair.
				</p>
				<div class="mt-4 p-4 border border-slate-200 rounded bg-white inline-block">
					<span class="font-['Inter'] text-xs font-bold uppercase tracking-widest text-slate-400 block mb-2">Certified Techs on Duty</span>
					<div class="flex gap-2">
						<span class="material-symbols-outlined text-orange-600" style="font-variation-settings: 'FILL' 1;">badge</span>
						<span class="material-symbols-outlined text-slate-400">badge</span>
						<span class="material-symbols-outlined text-slate-400">badge</span>
					</div>
				</div>
			</div>
			<div class="flex flex-col md:items-end gap-2">
				<a class="font-['Inter'] text-sm text-slate-500 hover:text-orange-600 transition-colors hover:underline" href="#">Technical Manuals</a>
				<a class="font-['Inter'] text-sm text-slate-500 hover:text-orange-600 transition-colors hover:underline" href="#">Privacy Policy</a>
				<a class="font-['Inter'] text-sm text-slate-500 hover:text-orange-600 transition-colors hover:underline" href="#">Terms of Service</a>
				<a class="font-['Inter'] text-sm text-slate-500 hover:text-orange-600 transition-colors hover:underline" href="#">API Status</a>
				<a class="font-['Inter'] text-sm text-slate-500 hover:text-orange-600 transition-colors hover:underline" href="#">Lab Locations</a>
				<a class="mt-6 font-['Inter'] text-sm font-bold text-slate-900 flex items-center gap-1 hover:text-orange-600 transition-colors" href="#">
					Read the Right to Repair Manifesto <span class="material-symbols-outlined text-sm">arrow_forward</span>
				</a>
			</div>
		</div>
	</footer>
</div>
