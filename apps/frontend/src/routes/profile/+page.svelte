<script lang="ts">
  import { 
    User, 
    MapPin, 
    Smartphone, 
    Clock, 
    CheckCircle2, 
    Package, 
    Settings, 
    LogOut,
    Mail,
    Phone,
    Calendar,
    ChevronRight,
    AlertCircle
  } from 'lucide-svelte';

  // Dummy User Data
  const user = {
    name: "Denden",
    email: "denden@example.com",
    phone: "+62 812-3456-7890",
    address: "Jl. Sudirman No. 123, Jakarta Pusat, DKI Jakarta 10220",
    memberSince: "May 2024",
    avatar: "https://ui-avatars.com/api/?name=Denden&background=0369a1&color=fff"
  };

  // Dummy Repair Tickets
  const tickets = [
    {
      id: "TK-7842",
      device: "iPhone 15 Pro",
      issue: "Broken Screen Replacement",
      status: "In Progress",
      statusColor: "text-blue-600 bg-blue-50 border-blue-100",
      progress: 65,
      date: "Oct 24, 2024",
      steps: ["Received", "Diagnosing", "Repairing", "Testing", "Ready"]
    },
    {
      id: "TK-7521",
      device: "MacBook Air M2",
      issue: "Battery Health Service",
      status: "Diagnosing",
      statusColor: "text-amber-600 bg-amber-50 border-amber-100",
      progress: 20,
      date: "Oct 28, 2024",
      steps: ["Received", "Diagnosing", "Repairing", "Testing", "Ready"]
    },
    {
      id: "TK-6910",
      device: "Google Pixel 8",
      issue: "Charging Port Repair",
      status: "Completed",
      statusColor: "text-emerald-600 bg-emerald-50 border-emerald-100",
      progress: 100,
      date: "Oct 15, 2024",
      steps: ["Received", "Diagnosing", "Repairing", "Testing", "Ready"]
    }
  ];

  let activeTab = $state('repairs');
  let repairFilter = $state('active');
</script>

<svelte:head>
  <title>Profile | OpenBench</title>
  <meta name="description" content="View your repair history, update personal details, and track device repair status on OpenBench." />
</svelte:head>

<div class="pt-24 pb-12 min-h-screen bg-slate-50 dark:bg-slate-950">
  <div class="container mx-auto px-4 max-w-6xl">
    
    <!-- Profile Header -->
    <div class="card-premium p-6 md:p-8 mb-8 flex flex-col md:flex-row items-center gap-6 md:gap-8">
      <div class="relative shrink-0">
        <img src={user.avatar} alt={user.name} class="w-24 h-24 md:w-32 md:h-32 rounded-2xl object-cover shadow-premium border-4 border-white dark:border-slate-800" />
        <div class="absolute -bottom-1 -right-1 md:-bottom-2 md:-right-2 w-6 h-6 md:w-8 md:h-8 bg-emerald-500 border-4 border-white dark:border-slate-900 rounded-full"></div>
      </div>
      
      <div class="flex-1 text-center md:text-left w-full">
        <div class="flex flex-col md:flex-row md:items-center gap-2 md:gap-4 mb-3 md:mb-2">
          <h1 class="text-2xl md:text-3xl font-bold text-slate-900 dark:text-white">{user.name}</h1>
          <span class="inline-flex items-center px-3 py-1 rounded-full text-[10px] md:text-xs font-bold bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400 border border-blue-200 dark:border-blue-800 uppercase tracking-wider mx-auto md:mx-0">
            Premium Customer
          </span>
        </div>
        
        <div class="flex flex-col sm:flex-row flex-wrap justify-center md:justify-start gap-2 sm:gap-4 text-slate-500 dark:text-slate-400 text-xs sm:text-sm">
          <div class="flex items-center justify-center gap-1.5">
            <Mail size={16} />
            <span>{user.email}</span>
          </div>
          <div class="flex items-center justify-center gap-1.5">
            <Phone size={16} />
            <span>{user.phone}</span>
          </div>
          <div class="flex items-center justify-center gap-1.5">
            <Calendar size={16} />
            <span>Member since {user.memberSince}</span>
          </div>
        </div>
      </div>

      <div class="flex gap-3 w-full sm:w-auto justify-center mt-2 md:mt-0">
        <button class="px-5 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 font-semibold text-sm text-slate-700 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-800 hover:text-slate-900 dark:hover:text-white transition-all cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-600/50">
          Edit Profile
        </button>
        <button class="w-10 h-10 flex items-center justify-center rounded-xl border border-slate-200 dark:border-slate-700 text-slate-500 hover:text-red-500 hover:border-red-200 dark:hover:border-red-900/50 transition-colors cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-red-500/50 shrink-0">
          <LogOut size={20} />
        </button>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
      
      <!-- Sidebar Navigation -->
      <aside class="lg:col-span-3">
        <nav class="card-premium p-2 !rounded-2xl space-y-1">
          <button 
            onclick={() => activeTab = 'repairs'}
            class="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-semibold transition-all cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-600/50 {activeTab === 'repairs' ? 'bg-blue-600 text-white shadow-lg shadow-blue-600/20' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800'}"
          >
            <Package size={18} />
            <span>My Repairs</span>
            {#if tickets.filter(t => t.status !== 'Completed').length > 0}
              <span class="ml-auto text-[10px] px-2 py-0.5 rounded-full {activeTab === 'repairs' ? 'bg-white/20 text-white' : 'bg-blue-100 text-blue-600 dark:bg-blue-900/40 dark:text-blue-400'}">
                {tickets.filter(t => t.status !== 'Completed').length}
              </span>
            {/if}
          </button>
          
          <button 
            onclick={() => activeTab = 'credentials'}
            class="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-semibold transition-all cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-600/50 {activeTab === 'credentials' ? 'bg-blue-600 text-white shadow-lg shadow-blue-600/20' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800'}"
          >
            <User size={18} />
            <span>Personal Info</span>
          </button>
          
          <button 
            disabled
            class="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-semibold text-slate-400 dark:text-slate-500 opacity-70 cursor-not-allowed"
          >
            <Settings size={18} />
            <span>Security Settings</span>
            <span class="ml-auto text-[9px] px-1.5 py-0.5 rounded bg-slate-100 dark:bg-slate-800 text-slate-400 font-bold tracking-widest">SOON</span>
          </button>
        </nav>
      </aside>

      <!-- Main Content -->
      <main class="lg:col-span-9">
        
        {#if activeTab === 'repairs'}
          <div class="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-300">
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
              <div class="flex bg-slate-100 dark:bg-slate-900 p-1 rounded-xl w-fit border border-slate-200 dark:border-slate-800">
                <button 
                  onclick={() => repairFilter = 'active'}
                  class="px-4 py-2 rounded-lg text-xs font-bold uppercase tracking-wider transition-all cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-600/50 {repairFilter === 'active' ? 'bg-white dark:bg-slate-800 text-blue-600 shadow-sm' : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300'}"
                >
                  Active Repairs
                </button>
                <button 
                  onclick={() => repairFilter = 'completed'}
                  class="px-4 py-2 rounded-lg text-xs font-bold uppercase tracking-wider transition-all cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-600/50 {repairFilter === 'completed' ? 'bg-white dark:bg-slate-800 text-blue-600 shadow-sm' : 'text-slate-500 hover:text-slate-700 dark:hover:text-slate-300'}"
                >
                  Completed
                </button>
              </div>
              <a href="/book" class="inline-flex justify-center items-center gap-2 px-4 py-2 bg-blue-600 text-white text-xs font-bold uppercase tracking-wider rounded-xl hover:bg-blue-700 transition-all shadow-lg shadow-blue-600/20 cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-600/50 w-full sm:w-auto">
                <Smartphone size={14} />
                New Repair Ticket
              </a>
            </div>

            {#each tickets.filter(t => repairFilter === 'active' ? t.status !== 'Completed' : t.status === 'Completed') as ticket}
              <a href="/track/{ticket.id}" class="card-premium overflow-hidden group !rounded-2xl hover:border-blue-500 hover:shadow-premium hover:-translate-y-1 block transition-all cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-600/50">
                <div class="p-6">
                  <div class="flex flex-col md:flex-row justify-between gap-4 mb-6">
                    <div class="flex items-start gap-4">
                      <div class="w-12 h-12 rounded-xl bg-slate-100 dark:bg-slate-800 flex items-center justify-center text-slate-500">
                        <Smartphone size={24} />
                      </div>
                      <div>
                        <div class="flex items-center gap-2 mb-1">
                          <h3 class="font-bold text-slate-900 dark:text-white text-lg">{ticket.device}</h3>
                          <span class="text-xs font-medium text-slate-400">{ticket.id}</span>
                        </div>
                        <p class="text-slate-500 dark:text-slate-400 text-sm">{ticket.issue}</p>
                      </div>
                    </div>
                    
                    <div class="flex flex-row md:flex-col items-center md:items-end justify-between md:justify-start gap-2">
                      <span class="px-3 py-1 rounded-full text-xs font-bold border dark:bg-slate-800 {ticket.statusColor}">
                        {ticket.status}
                      </span>
                      <span class="text-xs text-slate-400">{ticket.date}</span>
                    </div>
                  </div>

                  <!-- Progress Section -->
                  <div class="space-y-4">
                    <div class="flex justify-between text-xs font-bold mb-1 uppercase tracking-wider">
                      <span class="text-slate-400">Completion Progress</span>
                      <span class={ticket.status === 'Completed' ? 'text-emerald-600' : 'text-blue-600'}>{ticket.progress}%</span>
                    </div>
                    <div class="w-full h-2.5 bg-slate-100 dark:bg-slate-800 rounded-full overflow-hidden">
                      <div class="h-full {ticket.status === 'Completed' ? 'bg-emerald-600' : 'bg-blue-600'} rounded-full transition-all duration-1000 ease-out" style="width: {ticket.progress}%"></div>
                    </div>
                    
                    <!-- Steps indicator -->
                    <div class="hidden md:flex justify-between pt-2">
                      {#each ticket.steps as step, i}
                        {@const isCompleted = ticket.status === 'Completed'}
                        {@const normalizedStatus = ticket.status === 'In Progress' ? 'Repairing' : ticket.status}
                        {@const stepIndex = isCompleted ? ticket.steps.length - 1 : ticket.steps.indexOf(normalizedStatus)}
                        <div class="flex flex-col items-center gap-2 group/step">
                          <div class="relative">
                            <div class="w-2 h-2 rounded-full transition-all duration-500 {i <= stepIndex ? (isCompleted ? 'bg-emerald-600' : 'bg-blue-600') : 'bg-slate-200 dark:bg-slate-700'}"></div>
                            {#if i === stepIndex && !isCompleted}
                              <div class="absolute inset-0 w-2 h-2 rounded-full bg-blue-600 motion-safe:animate-ping opacity-40"></div>
                            {/if}
                          </div>
                          <span class="text-[10px] font-bold uppercase transition-colors duration-500 {i <= stepIndex ? (isCompleted ? 'text-emerald-600' : 'text-slate-700 dark:text-slate-300') : 'text-slate-400'}">{step}</span>
                        </div>
                      {/each}
                    </div>
                  </div>
                </div>
                
                <div class="bg-slate-50/50 dark:bg-slate-900/50 px-6 py-3 border-t border-slate-100 dark:border-slate-800 flex justify-between items-center transition-colors group-hover:bg-blue-50/50 dark:group-hover:bg-blue-900/10">
                  <div class="flex items-center gap-2 text-xs text-slate-500 font-medium">
                    <Clock size={14} />
                    <span>Last updated: {ticket.date}</span>
                  </div>
                  <span class="flex items-center gap-1 text-xs font-bold text-blue-600 group-hover:translate-x-1 transition-transform">
                    View Details
                    <ChevronRight size={14} />
                  </span>
                </div>
              </a>
            {:else}
              <div class="card-premium p-12 text-center flex flex-col items-center justify-center !rounded-2xl border-dashed border-2 bg-slate-50/50 dark:bg-slate-900/50 shadow-none">
                <div class="w-16 h-16 bg-white dark:bg-slate-800 rounded-full flex items-center justify-center text-slate-400 mb-4 shadow-sm">
                  {#if repairFilter === 'active'}
                    <Package size={32} />
                  {:else}
                    <CheckCircle2 size={32} class="text-emerald-500" />
                  {/if}
                </div>
                <h3 class="text-lg font-bold text-slate-900 dark:text-white mb-2">No {repairFilter} repairs</h3>
                <p class="text-slate-500 dark:text-slate-400 text-sm max-w-sm mb-6">
                  {#if repairFilter === 'active'}
                    You don't have any devices currently in the shop.
                  {:else}
                    You don't have any completed repairs yet.
                  {/if}
                </p>
                {#if repairFilter === 'active'}
                  <a href="/book" class="inline-flex items-center gap-2 text-sm font-bold text-blue-600 hover:text-blue-700 transition-colors cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-600/50 rounded-lg px-2 py-1">
                    Start a new repair <ChevronRight size={16} />
                  </a>
                {/if}
              </div>
            {/each}
          </div>

        {:else if activeTab === 'credentials'}
          <div class="card-premium p-6 md:p-8 !rounded-2xl animate-in fade-in slide-in-from-bottom-4 duration-300">
            <div class="flex items-center gap-4 mb-8">
              <div class="w-12 h-12 rounded-xl bg-blue-50 dark:bg-blue-900/20 text-blue-600 flex items-center justify-center">
                <User size={24} />
              </div>
              <div>
                <h2 class="text-xl font-bold text-slate-900 dark:text-white">
                  Personal Information
                </h2>
                <p class="text-slate-500 text-sm">Update your account details and preferences.</p>
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="space-y-2">
                <label for="fullName" class="text-sm font-bold text-slate-700 dark:text-slate-300 ml-1">Full Name</label>
                <div class="relative">
                  <User class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
                  <input id="fullName" type="text" value={user.name} class="w-full pl-10 pr-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-950 text-slate-900 dark:text-white focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-600/50 focus:border-blue-600 transition-colors" />
                </div>
              </div>
              
              <div class="space-y-2">
                <label for="emailAddr" class="text-sm font-bold text-slate-700 dark:text-slate-300 ml-1">Email Address</label>
                <div class="relative">
                  <Mail class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
                  <input id="emailAddr" type="email" value={user.email} class="w-full pl-10 pr-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-950 text-slate-900 dark:text-white focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-600/50 focus:border-blue-600 transition-colors" />
                </div>
              </div>

              <div class="space-y-2">
                <label for="phoneNum" class="text-sm font-bold text-slate-700 dark:text-slate-300 ml-1">Phone Number</label>
                <div class="relative">
                  <Phone class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
                  <input id="phoneNum" type="tel" value={user.phone} class="w-full pl-10 pr-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-950 text-slate-900 dark:text-white focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-600/50 focus:border-blue-600 transition-colors" />
                </div>
              </div>

              <div class="md:col-span-2 space-y-2">
                <label for="deliveryAddr" class="text-sm font-bold text-slate-700 dark:text-slate-300 ml-1">Delivery Address</label>
                <div class="relative">
                  <MapPin class="absolute left-3 top-3 text-slate-400" size={18} />
                  <textarea id="deliveryAddr" class="w-full pl-10 pr-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-950 text-slate-900 dark:text-white focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-600/50 focus:border-blue-600 transition-colors h-24">{user.address}</textarea>
                </div>
              </div>
            </div>
            
            <div class="mt-8 pt-6 border-t border-slate-100 dark:border-slate-800 flex justify-end">
              <button class="px-8 py-3 bg-blue-600 text-white font-bold rounded-xl hover:bg-blue-700 shadow-lg shadow-blue-600/20 active:scale-95 transition-all cursor-pointer focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-600/50">
                Save Changes
              </button>
            </div>
          </div>
        {/if}

      </main>
    </div>
  </div>
</div>
