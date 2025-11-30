<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { appState, Protocol } from "../state.svelte";
  import { fade, scale } from "svelte/transition";

  const dispatch = createEventDispatcher();

  let deviceId = $state("ubiquiti-device-access-point-05da");
  let alias = $state("U7 Pro Max Ultimate");
  let host = $state("ubiquiti-device-access-point");
  let port = $state(8080);
  let portGateway = $state(8081);
  let protocol = $state(Protocol.PROTOCOL_HTTP);
  let loading = $state(false);
  let error = $state("");

  async function handleSubmit() {
    if (!appState.client) {
      error = "Client not initialized";
      return;
    }

    loading = true;
    error = "";

    const device = {
      device_id: deviceId,
      alias,
      host,
      port,
      port_gateway: portGateway,
      protocol,
    };

    try {
      await appState.client.registerDevice(device);
      dispatch("close");
      dispatch("success");
    } catch (e: any) {
      console.error(e);
      error = e.message || "Failed to register device";
    } finally {
      loading = false;
    }
  }
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
  class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/60 backdrop-blur-md"
  transition:fade
  onclick={() => dispatch("close")}
>
  <div
    class="w-full max-w-md glass-panel rounded-3xl shadow-2xl overflow-hidden border border-white/10"
    onclick={(e) => e.stopPropagation()}
    transition:scale={{ start: 0.95, duration: 200 }}
  >
    <div
      class="p-6 border-b border-white/5 flex justify-between items-center bg-white/5"
    >
      <h3 class="text-lg font-bold text-white tracking-tight">
        Add New Device
      </h3>
      <button
        onclick={() => dispatch("close")}
        class="text-zinc-500 hover:text-white transition-colors"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="w-5 h-5"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>

    <div class="p-6 space-y-5">
      {#if error}
        <div
          class="p-3 bg-red-500/10 border border-red-500/20 rounded-xl text-red-400 text-sm flex items-center gap-2"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="w-4 h-4"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            ><circle cx="12" cy="12" r="10"></circle><line
              x1="12"
              y1="8"
              x2="12"
              y2="12"
            ></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg
          >
          {error}
        </div>
      {/if}

      <div class="space-y-2">
        <label
          class="text-xs font-bold text-zinc-500 uppercase tracking-wider ml-1"
          >Device ID</label
        >
        <input
          bind:value={deviceId}
          type="text"
          placeholder="e.g. device-01"
          class="w-full bg-black/20 border border-white/10 rounded-xl px-4 py-3 text-white placeholder-zinc-600 focus:outline-none focus:border-indigo-500/50 focus:ring-1 focus:ring-indigo-500/50 transition-all font-mono text-sm"
        />
      </div>

      <div class="space-y-2">
        <label
          class="text-xs font-bold text-zinc-500 uppercase tracking-wider ml-1"
          >Alias</label
        >
        <input
          bind:value={alias}
          type="text"
          placeholder="e.g. Production Server"
          class="w-full bg-black/20 border border-white/10 rounded-xl px-4 py-3 text-white placeholder-zinc-600 focus:outline-none focus:border-indigo-500/50 focus:ring-1 focus:ring-indigo-500/50 transition-all text-sm"
        />
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div class="space-y-2">
          <label
            class="text-xs font-bold text-zinc-500 uppercase tracking-wider ml-1"
            >Host</label
          >
          <input
            bind:value={host}
            type="text"
            class="w-full bg-black/20 border border-white/10 rounded-xl px-4 py-3 text-white placeholder-zinc-600 focus:outline-none focus:border-indigo-500/50 focus:ring-1 focus:ring-indigo-500/50 transition-all font-mono text-sm"
          />
        </div>
        <div class="space-y-2">
          <label
            class="text-xs font-bold text-zinc-500 uppercase tracking-wider ml-1"
            >Port</label
          >
          <input
            bind:value={port}
            type="number"
            class="w-full bg-black/20 border border-white/10 rounded-xl px-4 py-3 text-white placeholder-zinc-600 focus:outline-none focus:border-indigo-500/50 focus:ring-1 focus:ring-indigo-500/50 transition-all font-mono text-sm"
          />
        </div>
      </div>

      <div class="space-y-2">
        <label
          class="text-xs font-bold text-zinc-500 uppercase tracking-wider ml-1"
          >Protocol</label
        >
        <div class="relative">
          <select
            bind:value={protocol}
            class="w-full bg-black/20 border border-white/10 rounded-xl px-4 py-3 text-white focus:outline-none focus:border-indigo-500/50 focus:ring-1 focus:ring-indigo-500/50 transition-all appearance-none text-sm"
          >
            <option value={Protocol.PROTOCOL_GRPC}>gRPC</option>
            <option value={Protocol.PROTOCOL_GRPC}>gRPC Stream</option>
            <option value={Protocol.PROTOCOL_HTTP}>HTTP</option>
            <option value={Protocol.PROTOCOL_GRPC}>HTTP Stream</option>
          </select>
          <div
            class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-zinc-500"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="w-4 h-4"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"><path d="M6 9l6 6 6-6" /></svg
            >
          </div>
        </div>
      </div>
    </div>

    <div class="p-6 border-t border-white/5 flex justify-end gap-3 bg-black/20">
      <button
        onclick={() => dispatch("close")}
        class="px-5 py-2.5 text-sm font-medium text-zinc-400 hover:text-white transition-colors"
      >
        Cancel
      </button>
      <button
        onclick={handleSubmit}
        disabled={loading}
        class="px-6 py-2.5 text-sm font-semibold bg-blue-500 hover:bg-blue-400 text-white rounded-xl shadow-[0_0_20px_rgba(99,102,241,0.3)] hover:shadow-[0_0_30px_rgba(99,102,241,0.5)] transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
      >
        {#if loading}
          <svg
            class="animate-spin h-4 w-4 text-white"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            ></circle>
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
        {/if}
        Register Device
      </button>
    </div>
  </div>
</div>
