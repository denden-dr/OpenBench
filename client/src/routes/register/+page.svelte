<script lang="ts">
	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');

	function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (password !== confirmPassword) {
			alert('Passkeys do not match!');
			return;
		}
		console.log('Registration attempt:', { email, password });
		// API call placeholder
	}
</script>

<svelte:head>
	<title>OpenBench - Registration</title>
</svelte:head>

<div class="min-h-[calc(100vh-64px)] flex flex-col items-center justify-center font-body-main text-body-main text-on-surface p-6 bg-surface">
	<div class="w-full max-w-[480px]">
		<div class="text-center mb-8">
			<h1 class="font-h1 text-4xl font-bold text-on-background mb-2">OpenBench</h1>
			<p class="font-body-main text-lg text-on-surface-variant">Account Creation Protocol</p>
		</div>

		<div class="bg-surface-stark border border-border-subtle rounded-DEFAULT p-8 shadow-sm">
			<div class="flex border-b border-border-subtle mb-6">
				<a href="/login" class="flex-1 pb-4 text-center font-button-label text-base font-semibold text-outline hover:text-on-background transition-colors">
					Sign In
				</a>
				<div class="flex-1 pb-4 text-center font-button-label text-base font-semibold text-on-background border-b-2 border-secondary">
					Register
				</div>
			</div>

			<form class="space-y-6" onsubmit={handleSubmit}>
				<div>
					<label class="block font-tech-data text-xs font-medium text-on-surface-variant uppercase mb-2" for="email">
						Identifier (Email)
					</label>
					<div class="relative">
						<span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-outline">mail</span>
						<input
							bind:value={email}
							class="w-full pl-10 pr-4 py-3 bg-surface-stark border border-border-subtle rounded-DEFAULT font-body-main text-base text-on-surface focus:outline-none focus:border-secondary focus:ring-1 focus:ring-secondary"
							id="email"
							placeholder="user@domain.com"
							type="email"
							required
						/>
					</div>
				</div>

				<div>
					<label class="block font-tech-data text-xs font-medium text-on-surface-variant uppercase mb-2" for="password">
						Passkey
					</label>
					<div class="relative">
						<span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-outline">lock</span>
						<input
							bind:value={password}
							class="w-full pl-10 pr-4 py-3 bg-surface-stark border border-border-subtle rounded-DEFAULT font-body-main text-base text-on-surface focus:outline-none focus:border-secondary focus:ring-1 focus:ring-secondary"
							id="password"
							placeholder="••••••••"
							type="password"
							required
						/>
					</div>
				</div>

				<div>
					<label class="block font-tech-data text-xs font-medium text-on-surface-variant uppercase mb-2" for="confirmPassword">
						Confirm Passkey
					</label>
					<div class="relative">
						<span class="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-outline">shield_lock</span>
						<input
							bind:value={confirmPassword}
							class="w-full pl-10 pr-4 py-3 bg-surface-stark border border-border-subtle rounded-DEFAULT font-body-main text-base text-on-surface focus:outline-none focus:border-secondary focus:ring-1 focus:ring-secondary"
							id="confirmPassword"
							placeholder="••••••••"
							type="password"
							required
						/>
					</div>
				</div>

				<button
					class="w-full py-4 bg-secondary text-on-primary font-button-label text-base font-semibold rounded-DEFAULT hover:bg-on-secondary-fixed transition-colors flex items-center justify-center gap-2"
					type="submit"
				>
					Establish Identity
					<span class="material-symbols-outlined">person_add</span>
				</button>
			</form>

			<div class="mt-6 flex items-center gap-4">
				<div class="flex-1 h-px bg-border-subtle"></div>
				<span class="font-tech-data text-xs text-outline uppercase">OR</span>
				<div class="flex-1 h-px bg-border-subtle"></div>
			</div>

			<button
				class="mt-6 w-full py-3 bg-surface-stark border border-border-subtle text-on-surface font-button-label text-base font-semibold rounded-DEFAULT hover:bg-surface-container-low transition-colors flex items-center justify-center gap-3"
				onclick={() => console.log('Google registration placeholder')}
			>
				<svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
					<path
						d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
						fill="#4285F4"
					></path>
					<path
						d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
						fill="#34A853"
					></path>
					<path
						d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
						fill="#FBBC05"
					></path>
					<path
						d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
						fill="#EA4335"
					></path>
				</svg>				Register with Google
			</button>
		</div>

		<div class="mt-8 bg-technical-dark rounded-DEFAULT p-4 flex items-start gap-4 border border-outline">
			<span class="material-symbols-outlined text-secondary" style="font-variation-settings: 'FILL' 1;">
				shield
			</span>
			<div>
				<h3 class="font-button-label text-base font-semibold text-on-primary mb-1">Security Protocol Active</h3>
				<p class="font-body-sm text-sm text-on-primary-container">
					Your data is yours. We use high-entropy encryption standards to secure your technical profiles.
				</p>
			</div>
		</div>
	</div>
</div>
