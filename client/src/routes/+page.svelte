<script lang="ts">
	import { fetchHealth, type HealthData } from "$lib/api";

	let status = $state<"loading" | "connected" | "error">("loading");
	let healthData = $state<HealthData | null>(null);
	let errorMessage = $state<string | null>(null);

	$effect(() => {
		fetchHealth()
			.then((res) => {
				healthData = res.data;
				status = "connected";
			})
			.catch((err) => {
				errorMessage = err.message;
				status = "error";
			});
	});
</script>

<div class="bg-surface-stark">
	<!-- Hero Section -->
	<section class="pt-32 pb-section-padding px-gutter bg-surface-stark">
		<div
			class="max-w-container-max mx-auto grid grid-cols-1 lg:grid-cols-12 gap-gutter items-center"
		>
			<div class="lg:col-span-5 flex flex-col gap-element-gap">
				<h1 class="font-h1 text-h1 text-on-background">
					Your Device, Unveiled.
				</h1>
				<p
					class="font-body-main text-body-main text-on-surface-variant max-w-lg"
				>
					Professional repair engineered with radical transparency.
					Every diagnostic, every component, every step—documented and
					mapped against open standards. We don't just fix it; we show
					you how.
				</p>
				<div class="flex flex-wrap gap-4 mt-4">
					<a
						class="inline-flex items-center justify-center bg-secondary-container text-surface-stark font-button-label text-button-label px-6 py-4 rounded-DEFAULT hover:brightness-90 transition-all shadow-sm"
						href="#"
					>
						Book a Repair
					</a>
					<a
						class="inline-flex items-center justify-center border border-border-subtle text-on-surface-variant font-button-label text-button-label px-6 py-4 rounded-DEFAULT hover:bg-surface-alt transition-colors"
						href="#docs"
					>
						View Documentation
					</a>
				</div>
			</div>
			<div class="lg:col-span-7 mt-12 lg:mt-0 relative">
				<div
					class="aspect-4/3 bg-surface-alt rounded-lg overflow-hidden border border-border-subtle relative group"
				>
					<!-- Placeholder for visually rich 3D render -->
					<div
						class="absolute inset-0 bg-linear-to-br from-surface-alt to-surface-dim opacity-50"
					></div>
					<img
						alt="Exploded view"
						class="w-full h-full object-cover mix-blend-multiply opacity-90 transition-transform duration-700 group-hover:scale-105"
						src="https://lh3.googleusercontent.com/aida-public/AB6AXuDfV6zQXcyIwPAlddNO-qM8dlD729PAf09Ep2RrtjszoyPkqwAafAZBBxFzxkTUYqBQruVTJ5czk3_xhry0nYwuEwdrrOSA2HaK8Jc1SMy9-xWa8HlEuGiuASMGOkHT46uYE9GXxeSIIM6J9_nRF-aLX2yUCdARPn-NutMpGdl19C7BcSFQqhjjYQOApsEk4exrFWUkQJi5UBADVCWsuYPhvNO2n5yJtbg3khM6wC5H97exmYDY2xh0fR279uAbxdulxqHK7ZbnWpZ8"
					/>
					<!-- Tech overlay elements -->
					<div
						class="absolute top-4 left-4 font-tech-data text-tech-data text-on-surface-variant bg-surface-stark/80 backdrop-blur-sm px-2 py-1 rounded border border-border-subtle"
					>
						FIG. 01 / LOGIC_BOARD_ASY
					</div>
					<div class="absolute bottom-4 right-4 flex gap-2">
						<span
							class="material-symbols-outlined text-secondary-container"
							style="font-variation-settings: 'FILL' 1;"
							>add_circle</span
						>
						<span
							class="material-symbols-outlined text-on-surface-variant"
							>architecture</span
						>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Transparency Engine (Features) -->
	<section
		class="py-section-padding px-gutter bg-surface-alt border-y border-border-subtle"
	>
		<div class="max-w-container-max mx-auto">
			<div class="text-center mb-16 max-w-2xl mx-auto">
				<h2 class="font-h2 text-h2 text-on-background mb-4">
					The Transparency Engine
				</h2>
				<p
					class="font-body-main text-body-main text-on-surface-variant"
				>
					Our service centers operate like clean rooms, driven by
					strict adherence to public technical manuals.
				</p>
			</div>
			<div class="grid grid-cols-1 md:grid-cols-3 gap-gutter">
				<!-- Feature 1 -->
				<div
					class="bg-surface-stark border border-border-subtle rounded-DEFAULT p-8 hover:shadow-sm transition-shadow"
				>
					<div
						class="w-12 h-12 rounded-full bg-surface-alt flex items-center justify-center mb-6 border border-border-subtle"
					>
						<span
							class="material-symbols-outlined text-primary-container"
							>architecture</span
						>
					</div>
					<h3
						class="font-button-label text-button-label text-on-background mb-3"
					>
						Precision Diagnosis
					</h3>
					<p
						class="font-body-sm text-body-sm text-on-surface-variant"
					>
						Every device undergoes a 42-point hardware scan.
						Anomalies are mapped directly to corresponding
						schematics before any physical intervention begins.
					</p>
				</div>
				<!-- Feature 2 -->
				<div
					class="bg-surface-stark border border-border-subtle rounded-DEFAULT p-8 hover:shadow-sm transition-shadow"
				>
					<div
						class="w-12 h-12 rounded-full bg-surface-alt flex items-center justify-center mb-6 border border-border-subtle"
					>
						<span
							class="material-symbols-outlined text-primary-container"
							>schema</span
						>
					</div>
					<h3
						class="font-button-label text-button-label text-on-background mb-3"
					>
						Documented Assembly
					</h3>
					<p
						class="font-body-sm text-body-sm text-on-surface-variant"
					>
						Repairs follow open-source teardown guides step-by-step.
						We provide you with the exact manual our technicians
						use, ensuring standardized quality.
					</p>
				</div>
				<!-- Feature 3 -->
				<div
					class="bg-surface-stark border border-border-subtle rounded-DEFAULT p-8 hover:shadow-sm transition-shadow relative overflow-hidden"
				>
					<div
						class="w-12 h-12 rounded-full bg-surface-alt flex items-center justify-center mb-6 border border-border-subtle"
					>
						<span
							class="material-symbols-outlined text-primary-container"
							>verified_user</span
						>
					</div>
					<h3
						class="font-button-label text-button-label text-on-background mb-3"
					>
						Verified Testing
					</h3>
					<p
						class="font-body-sm text-body-sm text-on-surface-variant"
					>
						Post-repair calibration is validated against OEM
						specifications. You receive a cryptographic verification
						hash of your device's final test results.
					</p>
					<div class="absolute -bottom-4 -right-4 opacity-5">
						<span class="material-symbols-outlined text-9xl"
							>verified_user</span
						>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Live Documentation Preview (The 'Open' Bench) -->
	<section
		class="py-section-padding px-gutter bg-technical-dark text-surface-stark"
		id="docs"
	>
		<div class="max-w-container-max mx-auto">
			<!-- Search Interaction -->
			<div class="max-w-3xl mx-auto mb-16 relative">
				<div
					class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none"
				>
					<span class="material-symbols-outlined text-outline-variant"
						>search</span
					>
				</div>
				<input
					class="w-full bg-inverse-surface border border-outline-variant text-surface-stark font-body-main text-body-main rounded-DEFAULT py-4 pl-12 pr-16 focus:outline-none focus:border-inverse-primary focus:ring-1 focus:ring-inverse-primary placeholder:text-outline-variant shadow-[0_8px_32px_rgba(0,0,0,0.4)]"
					placeholder="Search technical manuals, parts, or error codes..."
					type="text"
				/>
				<div
					class="absolute inset-y-0 right-0 pr-4 flex items-center pointer-events-none"
				>
					<span
						class="font-tech-data text-tech-data text-outline-variant bg-technical-dark px-2 py-1 rounded border border-outline-variant text-xs"
						>Press /</span
					>
				</div>
			</div>
			<!-- The Documentation Block -->
			<div
				class="border border-outline-variant rounded-DEFAULT bg-inverse-surface overflow-hidden shadow-[0_12px_40px_rgba(0,0,0,0.5)]"
			>
				<!-- Header -->
				<div
					class="border-b border-outline-variant px-6 py-4 flex items-center justify-between bg-[#1f2d3a]"
				>
					<div
						class="flex items-center gap-2 font-tech-data text-tech-data text-inverse-primary tracking-widest text-xs"
					>
						<span class="material-symbols-outlined text-sm"
							>folder_open</span
						>
						DOCS / IPHONE 15 / LOGIC_BOARD_REPLACEMENT
					</div>
					<div class="flex items-center gap-4 text-outline-variant">
						<span
							class="material-symbols-outlined hover:text-surface-stark cursor-pointer transition-colors text-sm"
							>share</span
						>
						<span
							class="material-symbols-outlined hover:text-surface-stark cursor-pointer transition-colors text-sm"
							>print</span
						>
					</div>
				</div>
				<!-- Content Grid -->
				<div class="grid grid-cols-1 lg:grid-cols-2">
					<!-- Left: Technical Illustration -->
					<div
						class="border-r border-outline-variant bg-[#111820] relative p-8 min-h-[400px] flex items-center justify-center"
					>
						<div
							class="absolute top-4 left-4 font-tech-data text-tech-data text-outline-variant text-xs"
						>
							VIEW: TOP_DOWN_XRAY
						</div>
						<!-- Decorative Blueprint Grid -->
						<div
							class="absolute inset-0"
							style="background-image: radial-gradient(#303032 1px, transparent 1px); background-size: 20px 20px; opacity: 0.2;"
						></div>
						<img
							alt="Logic board schematic"
							class="relative z-10 max-w-full h-auto mix-blend-screen opacity-80"
							src="https://lh3.googleusercontent.com/aida-public/AB6AXuBxgWtRjtUTd9Yu0FOOliI3LzWEUeZbc92vM_MHR4ikRaedeAfoQv7R1V5PhSFGxzvfeWsYOsy3F-hEJdCSFNg1trlsGK8eU9cPGvHYcI0VxrFEln6cwi1-6jpHe-BKzQRLlbEFsdc3iV3qyegfkUKNTUoNRcaZ5p04FogTssBuJuvE7sCuUgvvcXLFT6GmSqMgO2RLfXeilKKI9sCXfgdIQp-t5jNSHP0c5HGToiaQLS7NpjT7dH-atPlGrQSyBl2v3snfaYmbI2rG"
						/>
						<div
							class="absolute bottom-8 right-8 w-8 h-8 border border-secondary-container rounded-full flex items-center justify-center animate-pulse"
						>
							<div
								class="w-2 h-2 bg-secondary-container rounded-full"
							></div>
						</div>
					</div>
					<!-- Right: Steps -->
					<div class="p-8 bg-inverse-surface">
						<div class="flex items-start gap-4 mb-8">
							<div
								class="w-8 h-8 rounded bg-technical-dark border border-outline-variant flex items-center justify-center font-tech-data text-tech-data text-inverse-primary shrink-0"
							>
								14
							</div>
							<div>
								<h3
									class="font-button-label text-button-label mb-2"
								>
									Re-seating the logic board
								</h3>
								<p
									class="font-body-sm text-body-sm text-outline-variant mb-4"
								>
									Carefully align the logic board with the
									standoffs in the rear case. Ensure no ribbon
									cables are trapped beneath the board before
									proceeding.
								</p>
								<div
									class="bg-[#1f2d3a] border border-secondary-container/30 border-l-2 border-l-secondary-container p-4 rounded-r-DEFAULT text-sm"
								>
									<strong
										class="text-secondary-fixed block mb-1 flex items-center gap-2"
										><span
											class="material-symbols-outlined text-sm"
											>warning</span
										> Critical Alignment</strong
									>
									<span class="text-outline-variant"
										>Do not force the board into place. If
										resistance is felt, check the display
										flex cables near the top edge.</span
									>
								</div>
							</div>
						</div>
						<div class="flex items-start gap-4 opacity-50">
							<div
								class="w-8 h-8 rounded bg-technical-dark border border-outline-variant flex items-center justify-center font-tech-data text-tech-data text-outline-variant shrink-0"
							>
								15
							</div>
							<div>
								<h3
									class="font-button-label text-button-label mb-1"
								>
									Securing grounding screws
								</h3>
								<p
									class="font-body-sm text-body-sm text-outline-variant"
								>
									Install the three 1.2mm Phillips screws
									along the grounding rail.
								</p>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Total Visibility Trust Section -->
	<section class="py-section-padding px-gutter bg-surface-stark">
		<div
			class="max-w-container-max mx-auto grid grid-cols-1 lg:grid-cols-2 gap-16 items-center"
		>
			<div>
				<h2 class="font-h2 text-h2 text-on-background mb-6">
					The End of the Black Box
				</h2>
				<p
					class="font-body-main text-body-main text-on-surface-variant mb-8"
				>
					Historically, device repair has been shrouded in mystery—a
					closed ecosystem designed to obscure the process. We're
					breaking that open. Every repair conducted on our benches is
					digitally logged and mapped against public documentation
					databases.
				</p>
				<ul class="space-y-4">
					<li class="flex items-start gap-3">
						<span
							class="material-symbols-outlined text-primary-container mt-1"
							>check_circle</span
						>
						<div>
							<strong
								class="font-button-label text-button-label text-on-background block"
								>OEM Sourced Components</strong
							>
							<span
								class="font-body-sm text-body-sm text-on-surface-variant"
								>We prioritize original parts when available, or
								mathematically equivalent high-grade
								alternatives.</span
							>
						</div>
					</li>
					<li class="flex items-start gap-3">
						<span
							class="material-symbols-outlined text-primary-container mt-1"
							>check_circle</span
						>
						<div>
							<strong
								class="font-button-label text-button-label text-on-background block"
								>No Proprietary Glue</strong
							>
							<span
								class="font-body-sm text-body-sm text-on-surface-variant"
								>We refuse to use structural adhesives that
								prevent future reparability.</span
							>
						</div>
					</li>
				</ul>
			</div>
			<div
				class="bg-surface-alt p-8 rounded-DEFAULT border border-border-subtle font-tech-data text-tech-data text-on-background flex flex-col gap-4"
			>
				<div class="text-xs text-on-surface-variant mb-2 uppercase">
					Live Diagnostic Feed // Platform: {healthData?.version ||
						"v0.0.0"}
				</div>

				<div
					class="flex items-center justify-between border-b border-border-subtle pb-2"
				>
					<span>BACKEND_CONNECTIVITY</span>
					{#if status === "loading"}
						<span class="text-neutral-400 animate-pulse uppercase"
							>Checking...</span
						>
					{:else if status === "connected"}
						<span
							class="text-emerald-600 font-bold uppercase tracking-widest"
							>Connected</span
						>
					{:else}
						<span
							class="text-rose-600 font-bold uppercase tracking-widest"
							>Error: {errorMessage || "Unknown"}</span
						>
					{/if}
				</div>

				<div
					class="flex items-center justify-between border-b border-border-subtle pb-2"
				>
					<span>BATTERY_HEALTH</span>
					<span class="text-secondary-container">FAIL (72% CAP)</span>
				</div>
				<div
					class="flex items-center justify-between border-b border-border-subtle pb-2"
				>
					<span>DISPLAY_DIGITIZER</span>
					<span class="text-on-background">PASS</span>
				</div>
				<div
					class="flex items-center justify-between border-b border-border-subtle pb-2"
				>
					<span>THERMAL_SENSOR_ARRAY</span>
					<span class="text-on-background">PASS</span>
				</div>
				<div
					class="flex items-center justify-between border-b border-border-subtle pb-2"
				>
					<span>CHARGING_PORT_AMP</span>
					<span class="text-secondary-container"
						>ANOMALY DETECTED</span
					>
				</div>
				<div
					class="mt-4 pt-4 border-t-2 border-primary-container flex items-center gap-2 text-xs"
				>
					<span class="material-symbols-outlined text-sm"
						>terminal</span
					>
					{#if status === "connected"}
						READY FOR REPAIR PROTOCOL GEN
					{:else}
						WAITING FOR BACKGROUND CONNECTIVITY...
					{/if}
				</div>
			</div>
		</div>
	</section>
</div>
