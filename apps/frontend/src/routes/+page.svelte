<script lang="ts">
  import { 
    Wrench, 
    Bell, 
    Shield, 
    ArrowRight, 
    Search, 
    Smartphone, 
    Battery, 
    Cpu, 
    Image as ImageIcon, 
    ClipboardCheck, 
    History,
    CheckCircle2,
    Calendar,
    Package,
    ShieldCheck,
    Zap,
    QrCode
  } from 'lucide-svelte';
</script>

<div class="flex flex-col">
  <!-- Hero Section -->
  <section class="relative pt-32 pb-20 lg:pt-48 lg:pb-32 overflow-hidden bg-white dark:bg-slate-950">
    <!-- Background Decor -->
    <div class="absolute top-0 right-0 w-1/3 h-full opacity-5 pointer-events-none hidden lg:block">
      <div class="w-full h-full bg-gradient-to-l from-blue-600 to-transparent"></div>
    </div>
    
    <div class="container mx-auto px-4 relative z-10">
      <div class="grid lg:grid-cols-2 gap-16 items-center">
        <div>
          <div class="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-blue-50 dark:bg-blue-900/30 text-sm font-semibold text-blue-600 dark:text-blue-400 mb-8 animate-fade-in">
            <span class="flex h-2 w-2 rounded-full bg-blue-600"></span>
            Certified Technician Service
          </div>
          
          <h1 class="text-5xl lg:text-7xl font-extrabold tracking-tight text-slate-900 dark:text-white mb-8 leading-[1.1]">
            Precision Repair. <br/>
            <span class="text-blue-600">Total Peace of Mind.</span>
          </h1>
          
          <p class="text-lg lg:text-xl text-slate-600 dark:text-slate-400 max-w-xl mb-10 leading-relaxed">
            Your device is essential. Our expert technicians use genuine parts and clinical precision to restore your iPhone or Android to factory standards. Fast, transparent, and guaranteed.
          </p>
          
          <div class="flex flex-col sm:flex-row items-center gap-4">
            <a href="/book" class="w-full sm:w-auto px-8 py-4 bg-blue-600 text-white font-bold rounded-xl shadow-lg shadow-blue-600/20 hover:bg-blue-700 transition-all active:scale-95 flex items-center justify-center gap-2">
              Book a Repair
              <ArrowRight size={18} />
            </a>
            <a href="/track" class="w-full sm:w-auto px-8 py-4 bg-white dark:bg-slate-900 text-slate-900 dark:text-white font-bold rounded-xl border border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800 transition-all active:scale-95 flex items-center justify-center">
              Track My Device
            </a>
          </div>

          <div class="mt-12 flex items-center gap-6">
            <div class="flex items-center gap-2 text-sm font-medium text-slate-500">
              <ShieldCheck class="text-green-500" size={20} />
              <span>90-Day Warranty on All Repairs</span>
            </div>
            <div class="h-4 w-px bg-slate-200 dark:bg-slate-800"></div>
            <div class="flex items-center gap-2 text-sm font-medium text-slate-500">
              <Zap class="text-yellow-500" size={20} />
              <span>Same-Day Service Available</span>
            </div>
          </div>
        </div>

        <div class="relative">
          <div class="rounded-3xl overflow-hidden shadow-2xl border-8 border-white dark:border-slate-900 bg-slate-100 dark:bg-slate-800 aspect-square group">
            <img 
              src="/images/hero.png" 
              alt="Precision Phone Repair" 
              class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-700" 
            />
          </div>
          
          <!-- Floating Badge -->
          <div class="absolute -bottom-6 -left-6 bg-white dark:bg-slate-900 p-6 rounded-2xl shadow-xl border border-slate-100 dark:border-slate-800 flex items-center gap-4 animate-bounce-subtle">
            <div class="bg-blue-50 dark:bg-blue-900/30 p-3 rounded-xl">
              <CheckCircle2 class="text-blue-600 dark:text-blue-400" />
            </div>
            <div>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Quality Guarantee</p>
              <p class="text-2xl font-bold text-slate-900 dark:text-white">Genuine Parts</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Quick Tracker Section -->
  <section class="py-12 bg-slate-50 dark:bg-slate-900/50">
    <div class="container mx-auto px-4">
      <div class="bg-white dark:bg-slate-900 p-8 lg:p-12 rounded-3xl shadow-soft border border-slate-200 dark:border-slate-800 -mt-24 relative z-20">
        <div class="grid lg:grid-cols-4 gap-8 items-end">
          <div class="lg:col-span-1">
            <h3 class="text-2xl font-bold text-slate-900 dark:text-white mb-2">Track Your Repair</h3>
            <p class="text-sm text-slate-500 dark:text-slate-400">Enter your Ticket ID for real-time status updates.</p>
          </div>
          <div class="space-y-2">
            <label class="text-xs font-bold text-slate-400 uppercase tracking-wider block">Ticket ID</label>
            <div class="relative">
              <input 
                type="text" 
                placeholder="e.g. FIX-98231" 
                class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl py-3 px-4 text-sm focus:outline-none focus:border-blue-600 dark:focus:border-blue-500 transition-colors"
              />
            </div>
          </div>
          <div class="space-y-2">
            <label class="text-xs font-bold text-slate-400 uppercase tracking-wider block">Phone Number</label>
            <input 
              type="tel" 
              placeholder="Last 4 digits" 
              class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl py-3 px-4 text-sm focus:outline-none focus:border-blue-600 dark:focus:border-blue-500 transition-colors"
            />
          </div>
          <a href="/track" class="bg-slate-900 dark:bg-blue-600 text-white w-full py-4 rounded-xl font-bold flex items-center justify-center gap-2 hover:bg-slate-800 dark:hover:bg-blue-700 transition-all active:scale-95">
            <Search size={18} />
            Track Status
          </a>
        </div>
      </div>
    </div>
  </section>

  <!-- Features Grid -->
  <section id="features" class="py-32 bg-white dark:bg-slate-950">
    <div class="container mx-auto px-4">
      <div class="text-center mb-20">
        <h2 class="text-sm font-bold text-blue-600 uppercase tracking-widest mb-4">Why Choose OpenBench</h2>
        <p class="text-4xl lg:text-5xl font-extrabold text-slate-900 dark:text-white">Expertise You Can Trust.</p>
      </div>

      <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
        <!-- Feature 1 -->
        <div class="group p-8 rounded-3xl border border-slate-100 dark:border-slate-900 hover:border-blue-100 dark:hover:border-blue-900/30 hover:bg-blue-50/30 dark:hover:bg-blue-900/5 transition-all duration-300">
          <div class="w-14 h-14 rounded-2xl bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 flex items-center justify-center mb-8 group-hover:scale-110 transition-transform">
            <Wrench size={28} />
          </div>
          <h3 class="text-xl font-bold text-slate-900 dark:text-white mb-4">Expert Diagnosis</h3>
          <p class="text-slate-600 dark:text-slate-400 leading-relaxed">
            Professional circuit-level assessment. We identify the root cause before you commit to a repair.
          </p>
        </div>

        <!-- Feature 2 -->
        <div class="group p-8 rounded-3xl border border-slate-100 dark:border-slate-900 hover:border-blue-100 dark:hover:border-blue-900/30 hover:bg-blue-50/30 dark:hover:bg-blue-900/5 transition-all duration-300">
          <div class="w-14 h-14 rounded-2xl bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 flex items-center justify-center mb-8 group-hover:scale-110 transition-transform">
            <QrCode size={28} />
          </div>
          <h3 class="text-xl font-bold text-slate-900 dark:text-white mb-4">Live Updates</h3>
          <p class="text-slate-600 dark:text-slate-400 leading-relaxed">
            Scan your QR receipt to see live progress. No phone calls needed—your device status is always online.
          </p>
        </div>

        <!-- Feature 3 -->
        <div class="group p-8 rounded-3xl border border-slate-100 dark:border-slate-900 hover:border-blue-100 dark:hover:border-blue-900/30 hover:bg-blue-50/30 dark:hover:bg-blue-900/5 transition-all duration-300">
          <div class="w-14 h-14 rounded-2xl bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 flex items-center justify-center mb-8 group-hover:scale-110 transition-transform">
            <ShieldCheck size={28} />
          </div>
          <h3 class="text-xl font-bold text-slate-900 dark:text-white mb-4">Guaranteed Parts</h3>
          <p class="text-slate-600 dark:text-slate-400 leading-relaxed">
            We only use high-grade components. Choose between Original or premium ODM parts with full warranty.
          </p>
        </div>

        <!-- Feature 4 -->
        <div class="group p-8 rounded-3xl border border-slate-100 dark:border-slate-900 hover:border-blue-100 dark:hover:border-blue-900/30 hover:bg-blue-50/30 dark:hover:bg-blue-900/5 transition-all duration-300">
          <div class="w-14 h-14 rounded-2xl bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 flex items-center justify-center mb-8 group-hover:scale-110 transition-transform">
            <ImageIcon size={28} />
          </div>
          <h3 class="text-xl font-bold text-slate-900 dark:text-white mb-4">Visual Evidence</h3>
          <p class="text-slate-600 dark:text-slate-400 leading-relaxed">
            We document every repair with before-and-after photos, giving you complete visibility of the work performed.
          </p>
        </div>
      </div>
    </div>
  </section>

  <!-- Engineered Section -->
  <section id="solutions" class="py-32 bg-slate-950 text-white overflow-hidden relative">
    <div class="absolute inset-0 opacity-10 pointer-events-none">
      <div class="w-full h-full bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] from-blue-600 via-transparent to-transparent scale-150"></div>
    </div>
    
    <div class="container mx-auto px-4 relative z-10">
      <div class="grid lg:grid-cols-2 gap-20 items-center">
        <div class="relative rounded-3xl overflow-hidden border border-slate-800 shadow-2xl">
          <img 
            src="/images/equipment.png" 
            alt="Professional Repair Equipment" 
            class="w-full h-[500px] object-cover" 
          />
        </div>
        
        <div>
          <h2 class="text-4xl lg:text-5xl font-extrabold mb-12 leading-tight">Clinical Standards. <br/><span class="text-blue-500">Expert Results.</span></h2>
          
          <div class="space-y-10">
            <div class="flex gap-6">
              <div class="w-12 h-12 rounded-xl bg-blue-600/20 text-blue-500 flex items-center justify-center flex-shrink-0">
                <ClipboardCheck size={24} />
              </div>
              <div>
                <h4 class="text-xl font-bold mb-2">Transparent Pricing</h4>
                <p class="text-slate-400 leading-relaxed">Know exactly what you're paying for. Detailed breakdown of parts, labor, and diagnosis fees—no surprises.</p>
              </div>
            </div>

            <div class="flex gap-6">
              <div class="w-12 h-12 rounded-xl bg-blue-600/20 text-blue-500 flex items-center justify-center flex-shrink-0">
                <History size={24} />
              </div>
              <div>
                <h4 class="text-xl font-bold mb-2">Extended Warranty</h4>
                <p class="text-slate-400 leading-relaxed">We stand by our work. Repairs are covered for up to 30 days, ensuring your device stays fixed.</p>
              </div>
            </div>

            <div class="flex gap-6">
              <div class="w-12 h-12 rounded-xl bg-blue-600/20 text-blue-500 flex items-center justify-center flex-shrink-0">
                <Battery size={24} />
              </div>
              <div>
                <h4 class="text-xl font-bold mb-2">Battery & Screen Health</h4>
                <p class="text-slate-400 leading-relaxed">High-precision calibration for touch sensitivity and battery health reporting that matches original specs.</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Journey Section -->
  <section class="py-32 bg-white dark:bg-slate-950">
    <div class="container mx-auto px-4 text-center mb-20">
      <h2 class="text-4xl font-extrabold text-slate-900 dark:text-white mb-4">How Your Repair Works</h2>
      <p class="text-slate-600 dark:text-slate-400">Simple, transparent, and professional from start to finish.</p>
    </div>

    <div class="container mx-auto px-4">
      <div class="relative flex flex-col md:flex-row justify-between items-center md:items-start gap-12">
        <!-- Progress Line -->
        <div class="hidden md:block absolute top-16 left-0 right-0 h-0.5 bg-slate-100 dark:bg-slate-900 -z-10"></div>
        
        <!-- Step 1 -->
        <div class="flex-1 flex flex-col items-center text-center group">
          <div class="w-32 h-32 rounded-full bg-slate-50 dark:bg-slate-900 flex items-center justify-center mb-8 border-2 border-transparent group-hover:border-blue-600 transition-all relative">
            <span class="absolute -top-2 -left-2 bg-blue-600 text-white w-10 h-10 rounded-full flex items-center justify-center font-bold shadow-lg">1</span>
            <Calendar size={40} class="text-slate-400 group-hover:text-blue-600 transition-colors" />
          </div>
          <h4 class="text-xl font-bold text-slate-900 dark:text-white mb-3">Book Your Slot</h4>
          <p class="text-sm text-slate-500 dark:text-slate-400 max-w-[200px]">Select your device and issue online for an instant estimate and slot.</p>
        </div>

        <!-- Step 2 -->
        <div class="flex-1 flex flex-col items-center text-center group">
          <div class="w-32 h-32 rounded-full bg-slate-50 dark:bg-slate-900 flex items-center justify-center mb-8 border-2 border-transparent group-hover:border-blue-600 transition-all relative">
            <span class="absolute -top-2 -left-2 bg-blue-600 text-white w-10 h-10 rounded-full flex items-center justify-center font-bold shadow-lg">2</span>
            <Package size={40} class="text-slate-400 group-hover:text-blue-600 transition-colors" />
          </div>
          <h4 class="text-xl font-bold text-slate-900 dark:text-white mb-3">Safe Drop-off</h4>
          <p class="text-sm text-slate-500 dark:text-slate-400 max-w-[200px]">Hand over your device and receive your secure digital intake receipt.</p>
        </div>

        <!-- Step 3 -->
        <div class="flex-1 flex flex-col items-center text-center group">
          <div class="w-32 h-32 rounded-full bg-slate-50 dark:bg-slate-900 flex items-center justify-center mb-8 border-2 border-transparent group-hover:border-blue-600 transition-all relative">
            <span class="absolute -top-2 -left-2 bg-blue-600 text-white w-10 h-10 rounded-full flex items-center justify-center font-bold shadow-lg">3</span>
            <CheckCircle2 size={40} class="text-slate-400 group-hover:text-blue-600 transition-colors" />
          </div>
          <h4 class="text-xl font-bold text-slate-900 dark:text-white mb-3">Track & Collect</h4>
          <p class="text-sm text-slate-500 dark:text-slate-400 max-w-[200px]">Watch your repair live. Collect your device the moment it's ready.</p>
        </div>
      </div>
    </div>
  </section>

  <!-- CTA Section -->
  <section class="py-32">
    <div class="container mx-auto px-4">
      <div class="bg-blue-600 rounded-[3rem] p-12 lg:p-24 text-center text-white relative overflow-hidden shadow-2xl shadow-blue-600/20">
        <!-- Decor -->
        <div class="absolute top-0 left-0 w-64 h-64 bg-white/10 rounded-full -translate-x-1/2 -translate-y-1/2 blur-3xl"></div>
        <div class="absolute bottom-0 right-0 w-64 h-64 bg-black/20 rounded-full translate-x-1/2 translate-y-1/2 blur-3xl"></div>
        
        <div class="relative z-10 max-w-2xl mx-auto">
          <h2 class="text-4xl lg:text-6xl font-extrabold mb-8">Broken device? Let's fix it today.</h2>
          <p class="text-lg lg:text-xl text-blue-100 mb-12 leading-relaxed">
            Get your instant repair estimate and join thousands of satisfied customers who trust OpenBench for their essential tech.
          </p>
          <div class="flex flex-col sm:flex-row items-center justify-center gap-4">
            <button class="w-full sm:w-auto px-10 py-5 bg-white text-blue-600 font-bold rounded-2xl hover:bg-blue-50 transition-all active:scale-95 text-lg">
              Start Your Booking
            </button>
            <button class="w-full sm:w-auto px-10 py-5 bg-blue-700/50 text-white font-bold rounded-2xl border border-white/20 hover:bg-blue-700 transition-all active:scale-95 text-lg">
              View Price List
            </button>
          </div>
        </div>
      </div>
    </div>
  </section>
</div>

<style>
  :global(html) {
    scroll-behavior: smooth;
  }

  .shadow-soft {
    box-shadow: 0 10px 30px -10px rgba(0, 0, 0, 0.05);
  }

  @keyframes bounce-subtle {
    0%, 100% { transform: translateY(0); }
    50% { transform: translateY(-10px); }
  }

  .animate-bounce-subtle {
    animation: bounce-subtle 4s ease-in-out infinite;
  }

  @keyframes fade-in {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .animate-fade-in {
    animation: fade-in 0.8s ease-out forwards;
  }
</style>

