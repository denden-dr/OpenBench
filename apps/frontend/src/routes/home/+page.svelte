<script lang="ts">
  import HomeHeader from "./components/HomeHeader.svelte";
  import HomeHero from "./components/HomeHero.svelte";
  import HomeBrands from "./components/HomeBrands.svelte";
  import HomeHowItWorks from "./components/HomeHowItWorks.svelte";
  import HomeFeatures from "./components/HomeFeatures.svelte";
  import HomeLocation from "./components/HomeLocation.svelte";
  import HomeFAQ from "./components/HomeFAQ.svelte";
  import HomeCTA from "./components/HomeCTA.svelte";
  import HomeFooter from "./components/HomeFooter.svelte";
</script>

<svelte:head>
  <title>OpenBench - Fast, Transparent Gadget Repairs</title>
  <meta
    name="description"
    content="Welcome to OpenBench, the premium repair shop orchestration engine. Track device status, submit repair requests, and view diagnostic updates."
  />
</svelte:head>

<div class="min-h-screen bg-neubrutalism-bg flex flex-col font-sans select-none text-neubrutalism-charcoal">
  <HomeHeader />
  <HomeHero />
  <HomeBrands />
  <HomeHowItWorks />
  <HomeFeatures />
  <HomeLocation />
  <HomeFAQ />
  <HomeCTA />
  <HomeFooter />
</div>
