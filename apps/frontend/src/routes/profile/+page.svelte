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
</script>

<svelte:head>
  <title>Profile | OpenBench</title>
</svelte:head>

<div class="pt-24 pb-12 min-h-screen bg-slate-50 dark:bg-slate-950">
  <div class="container mx-auto px-4 max-w-6xl">
    
    <!-- Profile Header -->
    <div class="card-premium p-8 mb-8 flex flex-col md:flex-row items-center gap-8">
      <div class="relative">
        <img src={user.avatar} alt={user.name} class="w-32 h-32 rounded-2xl object-cover shadow-premium border-4 border-white dark:border-slate-800" />
        <div class="absolute -bottom-2 -right-2 w-8 h-8 bg-emerald-500 border-4 border-white dark:border-slate-900 rounded-full"></div>
      </div>
      
      <div class="flex-1 text-center md:text-left">
        <div class="flex flex-col md:flex-row md:items-center gap-2 md:gap-4 mb-2">
          <h1 class="text-3xl font-bold text-slate-900 dark:text-white">{user.name}</h1>
          <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-bold bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400 border border-blue-200 dark:border-blue-800 uppercase tracking-wider">
            Premium Customer
          </span>
        </div>
        
        <div class="flex flex-wrap justify-center md:justify-start gap-4 text-slate-500 dark:text-slate-400 text-sm">
          <div class="flex items-center gap-1.5">
            <Mail size={16} />
            <span>{user.email}</span>
          </div>
          <div class="flex items-center gap-1.5">
            <Phone size={16} />
            <span>{user.phone}</span>
          </div>
          <div class="flex items-center gap-1.5">
            <Calendar size={16} />
            <span>Member since {user.memberSince}</span>
          </div>
        </div>
      </div>

      <div class="flex gap-3">
        <button class="px-5 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 font-semibold text-sm hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors">
          Edit Profile
        </button>
        <button class="w-10 h-10 flex items-center justify-center rounded-xl border border-slate-200 dark:border-slate-700 text-slate-500 hover:text-red-500 transition-colors">
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
            class="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-semibold transition-all {activeTab === 'repairs' ? 'bg-blue-600 text-white shadow-lg shadow-blue-600/20' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800'}"
          >
            <Package size={18} />
            <span>My Repairs</span>
            {#if tickets.filter(t => t.status !== 'Completed').length > 0}
              <span class="ml-auto bg-white/20 text-white text-[10px] px-2 py-0.5 rounded-full">
                {tickets.filter(t => t.status !== 'Completed').length}
              </span>
            {/if}
          </button>
          
          <button 
            onclick={() => activeTab = 'credentials'}
            class="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-semibold transition-all {activeTab === 'credentials' ? 'bg-blue-600 text-white shadow-lg shadow-blue-600/20' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800'}"
          >
            <User size={18} />
            <span>Personal Info</span>
          </button>
          
          <button 
            onclick={() => activeTab = 'settings'}
            class="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-semibold transition-all {activeTab === 'settings' ? 'bg-blue-600 text-white shadow-lg shadow-blue-600/20' : 'text-slate-600 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800'}"
          >
            <Settings size={18} />
            <span>Security Settings</span>
          </button>
        </nav>
      </aside>

      <!-- Main Content -->
      <main class="lg:col-span-9">
        
        {#if activeTab === 'repairs'}
          <div class="space-y-6">
            <div class="flex items-center justify-between">
              <h2 class="text-xl font-bold text-slate-900 dark:text-white">Active Repairs</h2>
              <a href="/book" class="text-blue-600 text-sm font-bold hover:underline">New Repair Ticket</a>
            </div>

            {#each tickets as ticket}
              <div class="card-premium overflow-hidden group !rounded-2xl hover:border-blue-500">
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
                      <span class="text-blue-600">{ticket.progress}%</span>
                    </div>
                    <div class="w-full h-2.5 bg-slate-100 dark:bg-slate-800 rounded-full overflow-hidden">
                      <div class="h-full bg-blue-600 rounded-full transition-all duration-1000 ease-out" style="width: {ticket.progress}%"></div>
                    </div>
                    
                    <!-- Steps indicator -->
                    <div class="hidden md:flex justify-between pt-2">
                      {#each ticket.steps as step, i}
                        {@const stepIndex = ticket.steps.indexOf(ticket.status)}
                        <div class="flex flex-col items-center gap-2 group/step">
                          <div class="w-2 h-2 rounded-full {i <= stepIndex ? 'bg-blue-600' : 'bg-slate-200 dark:bg-slate-700'}"></div>
                          <span class="text-[10px] font-bold uppercase {i <= stepIndex ? 'text-slate-700 dark:text-slate-300' : 'text-slate-400'}">{step}</span>
                        </div>
                      {/each}
                    </div>
                  </div>
                </div>
                
                <div class="bg-slate-50/50 dark:bg-slate-900/50 px-6 py-3 border-t border-slate-100 dark:border-slate-800 flex justify-between items-center">
                  <div class="flex items-center gap-2 text-xs text-slate-500 font-medium">
                    <AlertCircle size={14} />
                    <span>Est. completion: 2-3 business days</span>
                  </div>
                  <button class="flex items-center gap-1 text-xs font-bold text-blue-600 hover:text-blue-700">
                    View Details
                    <ChevronRight size={14} />
                  </button>
                </div>
              </div>
            {/each}
          </div>

        {:else}
          <div class="card-premium p-8 !rounded-2xl">
            <div class="flex items-center gap-4 mb-8">
              <div class="w-12 h-12 rounded-xl bg-blue-50 dark:bg-blue-900/20 text-blue-600 flex items-center justify-center">
                {#if activeTab === 'credentials'}
                  <User size={24} />
                {:else}
                  <Settings size={24} />
                {/if}
              </div>
              <div>
                <h2 class="text-xl font-bold text-slate-900 dark:text-white">
                  {activeTab === 'credentials' ? 'Personal Information' : 'Security Settings'}
                </h2>
                <p class="text-slate-500 text-sm">Update your account details and preferences.</p>
              </div>
            </div>

            {#if activeTab === 'credentials'}
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="space-y-2">
                  <label class="text-sm font-bold text-slate-700 dark:text-slate-300 ml-1">Full Name</label>
                  <div class="relative">
                    <User class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
                    <input type="text" value={user.name} class="w-full pl-10 pr-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-950 focus:outline-none focus:border-blue-600 transition-colors" />
                  </div>
                </div>
                
                <div class="space-y-2">
                  <label class="text-sm font-bold text-slate-700 dark:text-slate-300 ml-1">Email Address</label>
                  <div class="relative">
                    <Mail class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
                    <input type="email" value={user.email} class="w-full pl-10 pr-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-950 focus:outline-none focus:border-blue-600 transition-colors" />
                  </div>
                </div>

                <div class="space-y-2">
                  <label class="text-sm font-bold text-slate-700 dark:text-slate-300 ml-1">Phone Number</label>
                  <div class="relative">
                    <Phone class="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
                    <input type="tel" value={user.phone} class="w-full pl-10 pr-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-950 focus:outline-none focus:border-blue-600 transition-colors" />
                  </div>
                </div>

                <div class="md:col-span-2 space-y-2">
                  <label class="text-sm font-bold text-slate-700 dark:text-slate-300 ml-1">Delivery Address</label>
                  <div class="relative">
                    <MapPin class="absolute left-3 top-3 text-slate-400" size={18} />
                    <textarea class="w-full pl-10 pr-4 py-2.5 rounded-xl border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-950 focus:outline-none focus:border-blue-600 transition-colors h-24">{user.address}</textarea>
                  </div>
                </div>
              </div>
              
              <div class="mt-8 pt-6 border-t border-slate-100 dark:border-slate-800 flex justify-end">
                <button class="px-8 py-3 bg-blue-600 text-white font-bold rounded-xl hover:bg-blue-700 shadow-lg shadow-blue-600/20 active:scale-95 transition-all">
                  Save Changes
                </button>
              </div>
            {:else}
              <div class="flex flex-col items-center justify-center py-12 text-center">
                <div class="w-20 h-20 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center text-slate-400 mb-4">
                  <Settings size={40} />
                </div>
                <h3 class="text-lg font-bold text-slate-900 dark:text-white mb-2">Advanced Settings</h3>
                <p class="text-slate-500 max-w-xs mx-auto mb-6">Security preferences and notification settings are coming soon.</p>
                <button class="text-blue-600 font-bold hover:underline">Go back to repairs</button>
              </div>
            {/if}
          </div>
        {/if}

      </main>
    </div>
  </div>
</div>
