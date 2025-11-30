<script lang="ts">
  import "../app.css";
  import Sidebar from "$lib/components/Sidebar.svelte";
  import { onMount } from "svelte";
  import { ApiClient, appState } from "$lib/state.svelte";

  onMount(() => {
    const client = new ApiClient("http://localhost:8081");
    appState.setClient(client);

    document.documentElement.classList.add("dark");
  });
</script>

<div
  class="bg-zinc-50 dark:bg-[#050505] text-zinc-600 dark:text-zinc-400 antialiased selection:bg-blue-500/20 min-h-screen relative overflow-x-hidden flex"
>
  <div
    class="fixed top-[-10%] left-[20%] w-[40vw] h-[40vw] rounded-full bg-blue-500/5 dark:bg-blue-900/10 blur-[120px] pointer-events-none transition-colors duration-1000"
  ></div>

  <Sidebar />

  <main
    class="relative z-10 w-full pl-28 md:pl-36 pr-8 md:pr-12 py-12 max-w-[1600px] mx-auto"
  >
    <slot />
  </main>
</div>
