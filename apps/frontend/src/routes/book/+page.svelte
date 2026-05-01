<script lang="ts">
  import { 
    Smartphone, 
    Tablet, 
    Laptop, 
    ChevronRight, 
    ChevronLeft, 
    User, 
    Mail, 
    Phone, 
    AlertCircle, 
    CheckCircle2, 
    Camera,
    Wrench,
    Clock,
    ShieldCheck
  } from 'lucide-svelte';
  import { fade, fly, slide } from 'svelte/transition';
  import { goto } from '$app/navigation';

  let step = $state(1);
  let category = $state<'phone' | 'tablet' | ''>('');
  let brand = $state('');
  let model = $state('');
  let customModel = $state('');
  let issues = $state<string[]>([]);
  let description = $state('');
  let passcode = $state('');
  let securityType = $state<'pin' | 'pattern'>('pin');
  let patternSequence = $state<number[]>([]); // [1, 2, 3, ...] for 3x3 grid
  let repairMode = $state(false);
  let accessories = $state<string[]>([]);
  let name = $state('');
  let email = $state('');
  let phone = $state('');
  let termsAgreed = $state(false);
  let isSubmitting = $state(false);
  let showSuccess = $state(false);

  function toggleIssue(issue: string) {
    if (issues.includes(issue)) {
      issues = issues.filter(i => i !== issue);
    } else {
      issues = [...issues, issue];
    }
  }

  function togglePatternNode(node: number) {
    if (patternSequence.includes(node)) {
      patternSequence = patternSequence.filter(n => n !== node);
    } else {
      patternSequence = [...patternSequence, node];
    }
  }

  const categories = [
    { id: 'phone' as const, label: 'Smartphone', icon: Smartphone, color: 'blue' },
    { id: 'tablet' as const, label: 'Tablet', icon: Tablet, color: 'purple' }
  ];

  const brands = {
    phone: ['Apple', 'Samsung', 'Google', 'OnePlus', 'Xiaomi', 'Oppo', 'Vivo', 'Sony'],
    tablet: ['Apple (iPad)', 'Samsung Galaxy Tab', 'Microsoft Surface', 'Google Pixel Tablet']
  };

  const modelsByBrand: Record<string, string[]> = {
    'Apple': ['iPhone 15 Pro Max', 'iPhone 15 Pro', 'iPhone 15', 'iPhone 14 Pro', 'iPhone 14', 'iPhone 13', 'iPhone 12', 'iPhone 11'],
    'Samsung': ['Galaxy S24 Ultra', 'Galaxy S24+', 'Galaxy S24', 'Galaxy S23 Ultra', 'Galaxy S22', 'Galaxy Z Fold 5', 'Galaxy Z Flip 5'],
    'Google': ['Pixel 8 Pro', 'Pixel 8', 'Pixel 7 Pro', 'Pixel 7', 'Pixel Fold', 'Pixel 6a'],
    'Apple (iPad)': ['iPad Pro 12.9 (M2)', 'iPad Pro 11 (M2)', 'iPad Air (M1)', 'iPad mini (6th Gen)', 'iPad (10th Gen)'],
    'Samsung Galaxy Tab': ['Galaxy Tab S9 Ultra', 'Galaxy Tab S9', 'Galaxy Tab S8', 'Galaxy Tab A9'],
    'OnePlus': ['OnePlus 12', 'OnePlus 11', 'OnePlus Open'],
    'Xiaomi': ['Xiaomi 14 Ultra', 'Xiaomi 13T', 'Redmi Note 13'],
  };

  const accessoryOptions = ['Case/Cover', 'SIM Card', 'SD Card', 'Charger/Cable', 'Original Box'];

  const commonIssues = [
    'Cracked Screen',
    'Battery Draining Fast',
    'Charging Port Issue',
    'Water Damage',
    'Software/Boot Loop',
    'Camera Failure',
    'Button/Touch Issue',
    'Other (Describe below)'
  ];

  const canContinue = $derived((() => {
    if (step === 1) {
      const hasModel = (model && model !== 'Other') || (model === 'Other' && customModel.trim().length > 0);
      return !!(category && brand && hasModel);
    }
    if (step === 2) {
      const hasSecurity = repairMode || (securityType === 'pin' && passcode.trim().length > 0) || (securityType === 'pattern' && patternSequence.length > 0);
      return !!(issues.length > 0 && description.trim().length > 0 && hasSecurity);
    }
    if (step === 3) return !!(name.trim().length > 0 && phone.trim().length > 0 && email.includes('@'));
    if (step === 4) return termsAgreed;
    return false;
  })());

  function nextStep() {
    if (step < 4 && canContinue) step++;
  }

  function prevStep() {
    if (step > 1) step--;
  }

  async function handleSubmit() {
    isSubmitting = true;
    // Mock API call
    await new Promise(resolve => setTimeout(resolve, 2000));
    isSubmitting = false;
    showSuccess = true;
  }

  const progressWidth = $derived((step / 4) * 100);
</script>

<div class="min-h-screen bg-white dark:bg-slate-950 pt-32 pb-20 px-4 transition-colors duration-500">
  <div class="max-w-3xl mx-auto">
    
    {#if !showSuccess}
      <!-- Header -->
      <div class="text-center mb-12">
        <h1 class="text-4xl font-black text-slate-900 dark:text-white mb-4 tracking-tight">
          Book Your <span class="text-blue-600">Repair</span>
        </h1>
        <p class="text-slate-500 dark:text-slate-400 text-lg">Precision diagnostics and expert repair for your professional gear.</p>
      </div>

      <!-- Progress Bar -->
      <div class="mb-12 relative">
        <div class="h-2 w-full bg-slate-100 dark:bg-slate-800 rounded-full overflow-hidden">
          <div 
            class="h-full bg-blue-600 transition-all duration-500 ease-out shadow-[0_0_15px_rgba(37,99,235,0.4)]"
            style="width: {progressWidth}%"
          ></div>
        </div>
        <div class="flex justify-between mt-4">
          {#each ['Device', 'Issue', 'Contact', 'Review'] as label, i}
            <span class="text-[10px] font-bold uppercase tracking-widest {step > i ? 'text-blue-600' : 'text-slate-400'}">
              {label}
            </span>
          {/each}
        </div>
      </div>

      <!-- Form Container -->
      <div class="bg-white dark:bg-slate-900 rounded-[2.5rem] shadow-premium border border-slate-100 dark:border-slate-800 overflow-hidden min-h-[500px] flex flex-col">
        
        <div class="flex-1 p-8 sm:p-12">
          {#if step === 1}
            <div in:fly={{ y: 20, duration: 400 }} out:fade>
              <h2 class="text-2xl font-bold text-slate-900 dark:text-white mb-8">What device are we fixing today?</h2>
              
              <!-- Categories -->
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-10">
                {#each categories as cat}
                  <button 
                    onclick={() => category = cat.id}
                    class="p-6 rounded-3xl border-2 transition-all flex flex-col items-center gap-4 group
                      {category === cat.id 
                        ? 'border-blue-600 bg-blue-50/50 dark:bg-blue-900/20' 
                        : 'border-slate-100 dark:border-slate-800 hover:border-slate-200 dark:hover:border-slate-700 bg-slate-50/50 dark:bg-slate-800/30'}"
                  >
                    <div class="w-16 h-16 rounded-2xl flex items-center justify-center transition-all
                      {category === cat.id ? 'bg-blue-600 text-white' : 'bg-white dark:bg-slate-800 text-slate-400 group-hover:text-slate-600 dark:group-hover:text-slate-300 shadow-soft'}">
                      <cat.icon size={32} />
                    </div>
                    <span class="font-bold text-sm uppercase tracking-wider {category === cat.id ? 'text-blue-600' : 'text-slate-500'}">{cat.label}</span>
                  </button>
                {/each}
              </div>

              {#if category}
                <div transition:slide class="space-y-6">
                  <div class="grid sm:grid-cols-2 gap-6">
                    <div class="space-y-2">
                      <label for="brand-select" class="text-xs font-bold text-slate-400 uppercase tracking-widest ml-1">Brand</label>
                      <select 
                        id="brand-select"
                        bind:value={brand}
                        class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl py-4 px-6 text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 dark:focus:border-blue-500 transition-all appearance-none"
                      >
                        <option value="">Select Brand</option>
                        {#if category === 'phone' || category === 'tablet'}
                          {#each (brands as any)[category] as b}
                            <option value={b}>{b}</option>
                          {/each}
                        {/if}
                        <option value="Other">Other Brand</option>
                      </select>
                    </div>
                    <div class="space-y-2">
                      <label for="model-select" class="text-xs font-bold text-slate-400 uppercase tracking-widest ml-1">Popular Models</label>
                      <select 
                        id="model-select"
                        bind:value={model}
                        class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl py-4 px-6 text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 dark:focus:border-blue-500 transition-all appearance-none"
                      >
                        <option value="">Select Model</option>
                        {#if modelsByBrand[brand]}
                          {#each modelsByBrand[brand] as m}
                            <option value={m}>{m}</option>
                          {/each}
                        {/if}
                        <option value="Other">Other / Custom Model</option>
                      </select>
                    </div>
                  </div>

                  {#if model === 'Other'}
                    <div in:slide out:fade class="space-y-2 mt-4">
                      <label for="custom-model" class="text-xs font-bold text-slate-400 uppercase tracking-widest ml-1">Specify Your Model</label>
                      <input 
                        id="custom-model"
                        bind:value={customModel}
                        placeholder="e.g. Rare Edition X1"
                        class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl py-4 px-6 text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 dark:focus:border-blue-500 transition-all"
                      />
                    </div>
                  {/if}
                </div>
              {/if}
            </div>
          
          {:else if step === 2}
            <div in:fly={{ x: 20, duration: 400 }} out:fade>
              <h2 class="text-2xl font-bold text-slate-900 dark:text-white mb-2">Technical Details</h2>
              <p class="text-slate-500 dark:text-slate-400 mb-8 text-sm">Please provide access and tracking details for a faster repair.</p>
              
              <div class="space-y-10">
                <!-- Issue Selection -->
                <div class="space-y-3">
                  <label class="text-xs font-bold text-slate-400 uppercase tracking-widest ml-1">Primary Symptom(s)</label>
                  <div class="grid grid-cols-2 sm:grid-cols-4 gap-2" role="group" aria-label="Symptom Selection">
                    {#each commonIssues as iss}
                      <button 
                        type="button"
                        onclick={() => toggleIssue(iss)}
                        class="px-3 py-3 rounded-xl border text-[10px] font-bold uppercase tracking-tight transition-all
                          {issues.includes(iss) 
                            ? 'bg-blue-600 border-blue-600 text-white shadow-md' 
                            : 'bg-slate-50 dark:bg-slate-800 border-slate-100 dark:border-slate-700 text-slate-500 dark:text-slate-400 hover:border-slate-300'}"
                      >
                        {iss}
                      </button>
                    {/each}
                  </div>
                </div>

                <!-- Device Access (PRD Requirement) -->
                <div class="p-6 rounded-3xl bg-slate-50 dark:bg-slate-800/30 border border-slate-100 dark:border-slate-800 space-y-6">
                  <div class="flex items-center justify-between">
                    <div>
                      <h4 class="text-sm font-bold text-slate-900 dark:text-white">Device Security</h4>
                      <p class="text-xs text-slate-500">Required for hardware testing</p>
                    </div>
                    <div class="flex bg-white dark:bg-slate-900 p-1 rounded-xl border border-slate-200 dark:border-slate-700">
                      <button 
                        onclick={() => { repairMode = false; securityType = 'pin'; }}
                        class="px-3 py-1.5 rounded-lg text-[10px] font-bold uppercase transition-all
                          {!repairMode && securityType === 'pin' ? 'bg-slate-900 dark:bg-blue-600 text-white' : 'text-slate-400'}"
                      >PIN</button>
                      <button 
                        onclick={() => { repairMode = false; securityType = 'pattern'; }}
                        class="px-3 py-1.5 rounded-lg text-[10px] font-bold uppercase transition-all
                          {!repairMode && securityType === 'pattern' ? 'bg-slate-900 dark:bg-blue-600 text-white' : 'text-slate-400'}"
                      >Pattern</button>
                      <button 
                        onclick={() => repairMode = true}
                        class="px-3 py-1.5 rounded-lg text-[10px] font-bold uppercase transition-all
                          {repairMode ? 'bg-slate-900 dark:bg-blue-600 text-white' : 'text-slate-400'}"
                      >Repair Mode</button>
                    </div>
                  </div>

                  {#if repairMode}
                    <div transition:slide class="flex gap-4 p-4 rounded-2xl bg-blue-50/50 dark:bg-blue-900/10 border border-blue-100 dark:border-blue-900/20">
                      <ShieldCheck class="text-blue-600 flex-shrink-0" size={20} />
                      <p class="text-[10px] text-blue-700 dark:text-blue-400 font-medium leading-relaxed uppercase tracking-wider">
                        Repair Mode enabled. Your personal data is protected while technicians perform diagnostics.
                      </p>
                    </div>
                  {:else if !repairMode && securityType === 'pin'}
                    <div transition:slide>
                      <label for="pin-input" class="sr-only">PIN or Passcode</label>
                      <input 
                        id="pin-input"
                        type="text" 
                        bind:value={passcode}
                        placeholder="Enter PIN or Passcode"
                        class="w-full bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-2xl py-4 px-6 text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 transition-all text-center tracking-[0.5em] font-mono text-xl"
                      />
                    </div>
                  {:else if !repairMode && securityType === 'pattern'}
                    <div transition:slide class="flex flex-col items-center gap-6">
                      <div class="grid grid-cols-3 gap-6 p-6 bg-white dark:bg-slate-900 rounded-[2rem] border border-slate-200 dark:border-slate-700 shadow-inner" role="group" aria-label="Pattern Grid">
                        {#each [1, 2, 3, 4, 5, 6, 7, 8, 9] as node}
                          <button 
                            type="button"
                            onclick={() => togglePatternNode(node)}
                            aria-label="Pattern Node {node}"
                            class="w-12 h-12 rounded-full border-2 transition-all flex items-center justify-center relative
                              {patternSequence.includes(node) 
                                ? 'bg-blue-600 border-blue-600 text-white shadow-[0_0_15px_rgba(37,99,235,0.4)]' 
                                : 'bg-slate-50 dark:bg-slate-800 border-slate-200 dark:border-slate-700 text-slate-400 hover:border-slate-300'}"
                          >
                            {#if patternSequence.indexOf(node) !== -1}
                              <span class="text-[10px] font-bold">{patternSequence.indexOf(node) + 1}</span>
                            {:else}
                              <div class="w-2 h-2 rounded-full bg-current opacity-20"></div>
                            {/if}
                          </button>
                        {/each}
                      </div>
                      
                      <div class="flex items-center gap-4">
                        <button 
                          type="button"
                          onclick={() => patternSequence = []}
                          class="text-[10px] font-bold text-slate-400 uppercase tracking-widest hover:text-slate-600 transition-colors"
                        >Reset Pattern</button>
                        <div class="h-4 w-px bg-slate-200 dark:bg-slate-800"></div>
                        <p class="text-[10px] font-bold text-blue-600 uppercase tracking-widest">
                          {patternSequence.length > 0 ? `Sequence: ${patternSequence.join(' → ')}` : 'Draw your pattern above'}
                        </p>
                      </div>
                    </div>
                  {/if}
                </div>

                <!-- Accessories Checklist (PRD Requirement) -->
                <div class="space-y-3">
                  <label class="text-xs font-bold text-slate-400 uppercase tracking-widest ml-1">Items left with device</label>
                  <div class="flex flex-wrap gap-2" role="group" aria-label="Accessories Checklist">
                    {#each accessoryOptions as acc}
                      <button 
                        type="button"
                        onclick={() => {
                          if (accessories.includes(acc)) accessories = accessories.filter(a => a !== acc);
                          else accessories = [...accessories, acc];
                        }}
                        class="px-4 py-2 rounded-full border text-[10px] font-bold uppercase tracking-widest transition-all
                          {accessories.includes(acc) 
                            ? 'bg-blue-600 border-blue-600 text-white' 
                            : 'bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-700 text-slate-500 hover:border-slate-400'}"
                      >
                        {acc}
                      </button>
                    {/each}
                  </div>
                </div>

                <div class="space-y-2">
                  <label for="desc-input" class="text-xs font-bold text-slate-400 uppercase tracking-widest ml-1">Additional Notes / Pre-issue condition</label>
                  <textarea 
                    id="desc-input"
                    bind:value={description}
                    placeholder="e.g. Fallen into water, or issue appeared after software update..."
                    rows="3"
                    class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl py-4 px-6 text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 dark:focus:border-blue-500 transition-all resize-none"
                  ></textarea>
                </div>
              </div>
            </div>

          {:else if step === 3}
            <div in:fly={{ x: 20, duration: 400 }} out:fade>
              <h2 class="text-2xl font-bold text-slate-900 dark:text-white mb-2">Your Contact Information</h2>
              <p class="text-slate-500 dark:text-slate-400 mb-8 text-sm">We'll use this to send you real-time tracking updates.</p>
              
              <div class="space-y-6">
                <div class="space-y-2">
                  <label for="name-input" class="text-xs font-bold text-slate-400 uppercase tracking-widest ml-1">Full Name</label>
                  <div class="relative">
                    <User size={18} class="absolute left-5 top-1/2 -translate-y-1/2 text-slate-400" />
                    <input 
                      id="name-input"
                      bind:value={name}
                      placeholder="John Doe"
                      class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl py-4 pl-14 pr-6 text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 dark:focus:border-blue-500 transition-all"
                    />
                  </div>
                </div>

                <div class="grid sm:grid-cols-2 gap-6">
                  <div class="space-y-2">
                    <label for="email-input" class="text-xs font-bold text-slate-400 uppercase tracking-widest ml-1">Email Address</label>
                    <div class="relative">
                      <Mail size={18} class="absolute left-5 top-1/2 -translate-y-1/2 text-slate-400" />
                      <input 
                        id="email-input"
                        type="email"
                        bind:value={email}
                        placeholder="john@example.com"
                        class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl py-4 pl-14 pr-6 text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 dark:focus:border-blue-500 transition-all"
                      />
                    </div>
                  </div>
                  <div class="space-y-2">
                    <label for="phone" class="text-xs font-bold text-slate-400 uppercase tracking-widest ml-1">Phone Number</label>
                    <div class="relative">
                      <Phone size={18} class="absolute left-5 top-1/2 -translate-y-1/2 text-slate-400" />
                      <input 
                        id="phone"
                        type="tel"
                        bind:value={phone}
                        placeholder="+1 (555) 000-0000"
                        class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-2xl py-4 pl-14 pr-6 text-slate-900 dark:text-white focus:outline-none focus:border-blue-600 dark:focus:border-blue-500 transition-all"
                      />
                    </div>
                  </div>
                </div>
              </div>
            </div>

          {:else if step === 4}
            <div in:fly={{ x: 20, duration: 400 }} out:fade>
              <h2 class="text-2xl font-bold text-slate-900 dark:text-white mb-2">Review Your Booking</h2>
              <p class="text-slate-500 dark:text-slate-400 mb-8 text-sm">Please verify your details before confirming.</p>
              
              <div class="grid gap-4">
                <div class="grid sm:grid-cols-2 gap-4">
                  <div class="p-6 rounded-3xl bg-slate-50 dark:bg-slate-800/50 border border-slate-100 dark:border-slate-800">
                    <span class="text-[10px] font-bold text-slate-400 uppercase tracking-widest block mb-1">Device</span>
                    <h4 class="font-bold text-slate-900 dark:text-white">
                      {brand} {model === 'Other' ? customModel : model}
                    </h4>
                  </div>
                  <div class="p-6 rounded-3xl bg-slate-50 dark:bg-slate-800/50 border border-slate-100 dark:border-slate-800">
                    <span class="text-[10px] font-bold text-slate-400 uppercase tracking-widest block mb-1">Security</span>
                    <h4 class="font-bold text-slate-900 dark:text-white">{repairMode ? 'Repair Mode Enabled' : 'Passcode Provided'}</h4>
                  </div>
                </div>

                <div class="p-6 rounded-3xl bg-slate-50 dark:bg-slate-800/50 border border-slate-100 dark:border-slate-800">
                  <div class="flex justify-between items-start mb-4">
                    <div>
                      <span class="text-[10px] font-bold text-slate-400 uppercase tracking-widest block mb-1">Repair Issue(s)</span>
                      <div class="flex flex-wrap gap-2">
                        {#each issues as iss}
                          <span class="px-2 py-1 bg-blue-600/10 text-blue-600 rounded-lg text-[10px] font-bold uppercase tracking-tight border border-blue-600/20">{iss}</span>
                        {/each}
                      </div>
                    </div>
                    <button onclick={() => step = 2} class="text-xs font-bold text-blue-600 uppercase">Edit</button>
                  </div>
                  <p class="text-sm text-slate-500 dark:text-slate-400 leading-relaxed italic mb-4">"{description}"</p>
                  
                  {#if accessories.length > 0}
                    <div class="flex flex-wrap gap-2">
                      {#each accessories as acc}
                        <span class="px-2 py-1 bg-white dark:bg-slate-900 rounded-lg text-[9px] font-bold text-slate-500 border border-slate-100 dark:border-slate-800">+{acc}</span>
                      {/each}
                    </div>
                  {/if}
                </div>

                <div class="p-6 rounded-3xl bg-slate-50 dark:bg-slate-800/50 border border-slate-100 dark:border-slate-800 flex justify-between items-center">
                  <div>
                    <span class="text-[10px] font-bold text-slate-400 uppercase tracking-widest block mb-1">Contact</span>
                    <h4 class="font-bold text-slate-900 dark:text-white">{name}</h4>
                    <span class="text-xs text-slate-500">{phone}</span>
                  </div>
                  <button onclick={() => step = 3} class="text-xs font-bold text-blue-600 uppercase">Edit</button>
                </div>

                <div class="mt-4 p-6 rounded-3xl bg-blue-50 dark:bg-blue-900/10 border border-blue-100 dark:border-blue-900/30 flex gap-4">
                  <div class="text-blue-600 flex-shrink-0">
                    <AlertCircle size={24} />
                  </div>
                  <div>
                    <h5 class="text-sm font-bold text-blue-900 dark:text-blue-100 mb-1">Diagnosis Fee Notice</h5>
                    <p class="text-xs text-blue-700 dark:text-blue-300 leading-relaxed">A standard $29 diagnostic fee applies to all intake devices. This fee is **waived** if you proceed with the recommended repair.</p>
                  </div>
                </div>

                <!-- Terms Checkbox (PRD Requirement) -->
                <div class="mt-4 flex items-start gap-3 p-2">
                  <input 
                    type="checkbox" 
                    id="terms" 
                    bind:checked={termsAgreed}
                    class="mt-1 w-5 h-5 rounded border-slate-300 text-blue-600 focus:ring-blue-500 transition-all cursor-pointer"
                  />
                  <label for="terms" class="text-xs text-slate-600 dark:text-slate-400 leading-relaxed cursor-pointer select-none">
                    I acknowledge that a <span class="font-bold text-slate-900 dark:text-white">$29 mandatory diagnosis fee</span> applies. I agree to the <span class="text-blue-600 underline">Terms of Service</span> and authorize technical assessment of my device.
                  </label>
                </div>
              </div>
            </div>
          {/if}
        </div>

        <!-- Footer Actions -->
        <div class="px-8 py-6 bg-slate-50 dark:bg-slate-800/50 border-t border-slate-100 dark:border-slate-800 flex justify-between items-center">
          <button 
            onclick={prevStep}
            disabled={step === 1 || isSubmitting}
            class="px-6 py-3 rounded-xl text-sm font-bold text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white disabled:opacity-0 transition-all flex items-center gap-2"
          >
            <ChevronLeft size={18} />
            Back
          </button>

          {#if step < 4}
            <button 
              onclick={nextStep}
              disabled={!canContinue}
              class="px-8 py-3 rounded-xl bg-blue-600 text-white font-bold text-sm shadow-lg shadow-blue-600/20 hover:bg-blue-700 transition-all flex items-center gap-2 disabled:opacity-50 disabled:grayscale disabled:cursor-not-allowed"
            >
              Continue
              <ChevronRight size={18} />
            </button>
          {:else}
            <button 
              onclick={handleSubmit}
              disabled={isSubmitting || !canContinue}
              class="px-10 py-4 rounded-2xl bg-blue-600 text-white font-bold text-sm shadow-xl shadow-blue-600/30 hover:bg-blue-700 transition-all flex items-center gap-3 disabled:opacity-50 disabled:grayscale disabled:cursor-not-allowed"
            >
              {#if isSubmitting}
                <div class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                Creating Ticket...
              {:else}
                Confirm & Create Ticket
                <CheckCircle2 size={18} />
              {/if}
            </button>
          {/if}
        </div>
      </div>

      <!-- Help Section -->
      <div class="mt-12 grid grid-cols-1 sm:grid-cols-3 gap-6">
        <div class="flex items-center gap-4 p-4">
          <div class="w-10 h-10 rounded-full bg-blue-50 dark:bg-blue-900/20 flex items-center justify-center text-blue-600">
            <Clock size={20} />
          </div>
          <div>
            <h5 class="text-sm font-bold text-slate-900 dark:text-white">Quick Intake</h5>
            <p class="text-xs text-slate-500">Processing in 15 mins</p>
          </div>
        </div>
        <div class="flex items-center gap-4 p-4">
          <div class="w-10 h-10 rounded-full bg-blue-50 dark:bg-blue-900/20 flex items-center justify-center text-blue-600">
            <ShieldCheck size={20} />
          </div>
          <div>
            <h5 class="text-sm font-bold text-slate-900 dark:text-white">Secure Data</h5>
            <p class="text-xs text-slate-500">Privacy protected</p>
          </div>
        </div>
        <div class="flex items-center gap-4 p-4">
          <div class="w-10 h-10 rounded-full bg-blue-50 dark:bg-blue-900/20 flex items-center justify-center text-blue-600">
            <Wrench size={20} />
          </div>
          <div>
            <h5 class="text-sm font-bold text-slate-900 dark:text-white">Expert Techs</h5>
            <p class="text-xs text-slate-500">Certified precision</p>
          </div>
        </div>
      </div>

    {:else}
      <!-- Success State -->
      <div in:fly={{ y: 20 }} class="max-w-xl mx-auto text-center py-12">
        <div class="w-24 h-24 rounded-full bg-green-100 dark:bg-green-900/20 text-green-600 flex items-center justify-center mx-auto mb-8 shadow-[0_0_40px_rgba(34,197,94,0.2)]">
          <CheckCircle2 size={48} />
        </div>
        <h2 class="text-4xl font-black text-slate-900 dark:text-white mb-4">Ticket Created!</h2>
        <p class="text-slate-500 dark:text-slate-400 text-lg mb-10 leading-relaxed">
          Your repair ticket <span class="font-bold text-blue-600">#OB-8829</span> has been successfully logged. Please visit our service center with your device to proceed with the diagnostics.
        </p>
        
        <div class="bg-white dark:bg-slate-900 rounded-[2.5rem] p-8 border border-slate-100 dark:border-slate-800 shadow-premium mb-10">
          <div class="flex items-center justify-between mb-6 pb-6 border-b border-slate-50 dark:border-slate-800">
            <div class="text-left">
              <span class="text-[10px] font-bold text-slate-400 uppercase tracking-widest block mb-1">Queue Status</span>
              <span class="inline-flex items-center gap-2 px-3 py-1 bg-green-50 dark:bg-green-900/20 text-green-600 rounded-full text-xs font-bold uppercase tracking-wide">
                <span class="w-1.5 h-1.5 rounded-full bg-green-500 animate-pulse"></span>
                Fast Lane Available
              </span>
            </div>
            <div class="text-right">
              <span class="text-[10px] font-bold text-slate-400 uppercase tracking-widest block mb-1">Est. Wait Time</span>
              <span class="text-xl font-black text-slate-900 dark:text-white">~15 Mins</span>
            </div>
          </div>
          
          <div class="space-y-4">
            <button 
              onclick={() => goto('/track/OB-8829')}
              class="w-full py-4 rounded-2xl bg-slate-900 dark:bg-white text-white dark:text-slate-900 font-bold text-sm hover:opacity-90 transition-all flex items-center justify-center gap-2"
            >
              Track Progress
              <ChevronRight size={18} />
            </button>
            <button 
              onclick={() => goto('/')}
              class="w-full py-4 rounded-2xl border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-400 font-bold text-sm hover:bg-slate-50 dark:hover:bg-slate-800 transition-all"
            >
              Return Home
            </button>
          </div>
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .shadow-premium {
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.08);
  }

  .shadow-soft {
    box-shadow: 0 10px 30px -10px rgba(0, 0, 0, 0.05);
  }
  
  select {
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24' stroke='%2394a3b8' stroke-width='2'%3E%3Cpath stroke-linecap='round' stroke-linejoin='round' d='M19 9l-7 7-7-7'%3E%3C/path%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 1.5rem center;
    background-size: 1rem;
  }
</style>
