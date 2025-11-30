<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { appState, DeviceStatus } from "$lib/state.svelte";
  import DeviceCard from "$lib/components/DeviceCard.svelte";
  import ThroughputChart from "$lib/components/ThroughputChart.svelte";
  import AddDeviceModal from "$lib/components/AddDeviceModal.svelte";
  import type { DiagnosticsResponse } from "$lib/state.svelte";

  let showAddDeviceModal = $state(false);
  let streamCancelers: (() => void)[] = [];

  let lastUpdate = $state<Date | null>(null);
  let timeSinceUpdate = $state("--");
  let updateInterval: any;

  function formatTimeSince(date: Date | null): string {
    if (!date) return "--";

    const seconds = Math.floor((Date.now() - date.getTime()) / 1000);

    if (seconds < 5) return "0s";
    if (seconds < 60) return `${seconds}s`;
    if (seconds < 3600) return `${Math.floor(seconds / 60)}m`;
    if (seconds < 86400) return `${Math.floor(seconds / 3600)}h`;
    return `${Math.floor(seconds / 86400)}d`;
  }

  function updateTimeSince() {
    timeSinceUpdate = formatTimeSince(lastUpdate);
  }

  function onStreamUpdate() {
    lastUpdate = new Date();
    updateTimeSince();
  }

  async function loadDevices() {
    if (!appState.client) return;
    try {
      const deviceList = await appState.client.listDevices();
      appState.setDevices(deviceList);
      startStreaming(deviceList);
    } catch (e) {
      console.error("Failed to load devices:", e);
    }
  }

  function startStreaming(deviceList: any[]) {
    streamCancelers.forEach((cancel) => cancel());
    streamCancelers = [];

    deviceList.forEach(async (device) => {
      if (!appState.client) return;
      const cancel = await appState.client.streamDiagnostics(
        device.device_id,
        (data: DiagnosticsResponse) => {
          if (data.diagnostics) {
            appState.addDiagnosticData(device.device_id, data.diagnostics);
            const hasOutage = appState.devices.some((d) => {
              const diag = appState.diagnosticsMap[d.device_id];
              return (
                diag?.device_status === DeviceStatus.DEVICE_STATUS_DEGRADED ||
                diag?.device_status === DeviceStatus.DEVICE_STATUS_ERROR
              );
            });
            appState.isOutageActive = hasOutage;

            onStreamUpdate();
          }
        },
        (err) => {
          console.error(`Stream error for ${device.device_id}:`, err);
        }
      );
      streamCancelers.push(cancel);
    });
  }

  $effect(() => {
    if (appState.client) {
      loadDevices();
    }
  });

  onMount(() => {
    window.addEventListener("openAddDeviceModal", () => {
      showAddDeviceModal = true;
    });
    updateInterval = setInterval(updateTimeSince, 1000);
  });

  onDestroy(() => {
    streamCancelers.forEach((cancel) => cancel());
    if (updateInterval) clearInterval(updateInterval);
  });

  function handleDeviceAdded() {
    loadDevices();
    showAddDeviceModal = false;
  }
</script>

<!-- Header -->
<header
  class="mb-10 flex flex-col gap-6 opacity-0 animate-fade-in"
  style="animation-fill-mode: forwards;"
>
  <div
    class="flex flex-col md:flex-row justify-between items-start md:items-end gap-6"
  >
    <div class="space-y-2">
      <div class="flex items-center gap-3">
        <div
          id="systemStatusBadge"
          class="inline-flex items-center gap-2 px-3 py-1 rounded-full border border-blue-200 dark:border-blue-500/20 bg-blue-50 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400 text-[11px] font-medium tracking-wide uppercase transition-all duration-500 {appState.isOutageActive
            ? 'border-red-200 dark:border-red-500/20 bg-red-50 dark:bg-red-500/10 text-red-600 dark:text-red-400'
            : 'border-blue-200 dark:border-blue-500/20 bg-blue-50 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400'} text-[11px] font-medium tracking-wide uppercase"
        >
          <span class="relative flex h-1.5 w-1.5">
            <span
              id="statusPing"
              class="animate-ping absolute inline-flex h-full w-full rounded-full {appState.isOutageActive
                ? 'bg-red-400'
                : 'bg-blue-400'} opacity-75"
            ></span>
            <span
              id="statusDot"
              class="relative inline-flex rounded-full h-1.5 w-1.5 transition-colors duration-500 {appState.isOutageActive
                ? 'bg-red-500'
                : 'bg-blue-500'}"
            ></span>
          </span>
          <span id="statusText"
            >{appState.isOutageActive
              ? "System Degraded"
              : "System Operational"}</span
          >
        </div>
      </div>
      <h1
        class="text-3xl font-medium text-zinc-900 dark:text-white tracking-tight"
      >
        Device Overview
      </h1>
    </div>

    <!-- Aggregates -->
    <div class="flex items-center gap-10">
      <div class="text-right group cursor-default">
        <div
          class="text-[11px] font-medium uppercase tracking-wider text-zinc-400 dark:text-zinc-500 mb-1"
        >
          Last Update
        </div>
        <div
          class="text-xl font-medium text-zinc-900 dark:text-zinc-200 tracking-tight font-mono"
        >
          {timeSinceUpdate}
        </div>
      </div>
      <div class="w-px h-10 bg-zinc-200 dark:bg-zinc-800"></div>
      <div class="text-right group cursor-default">
        <div
          class="text-[11px] font-medium uppercase tracking-wider text-zinc-400 dark:text-zinc-500 mb-1"
        >
          Active Devices
        </div>
        <div class="flex items-center justify-end gap-1">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="12"
            height="12"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="text-blue-500"
            ><path d="m5 12 7-7 7 7" /><path d="M12 19V5" /></svg
          >
          <div
            class="text-xl font-medium text-zinc-900 dark:text-zinc-200 tracking-tight font-mono"
          >
            {appState.devices.length || 0}
          </div>
        </div>
      </div>
    </div>
  </div>
</header>

<!-- Chart Section -->
<section
  class="mb-8 opacity-0 animate-slide-up"
  style="animation-delay: 0.1s; animation-fill-mode: forwards;"
>
  <ThroughputChart />
</section>

<!-- Devices Grid -->
<section
  class="opacity-0 animate-slide-up"
  style="animation-delay: 0.2s; animation-fill-mode: forwards;"
>
  <div class="flex justify-between items-center mb-6 px-1">
    <h3
      class="text-xs font-semibold text-zinc-500 dark:text-zinc-500 uppercase tracking-widest"
    >
      Infrastructure
    </h3>
    <div class="flex items-center gap-2">
      <span class="w-1 h-1 rounded-full bg-zinc-300 dark:bg-zinc-700"></span>
      <span class="text-xs text-zinc-600 dark:text-zinc-400"
        >{appState.devices.length} Devices</span
      >
    </div>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-2 gap-4">
    {#each appState.devices as device (device.device_id)}
      <DeviceCard
        {device}
        diagnostics={appState.diagnosticsMap[device.device_id]}
        history={appState.diagnosticsHistory[device.device_id] || []}
      />
    {/each}
  </div>
</section>

{#if showAddDeviceModal}
  <AddDeviceModal
    on:close={() => (showAddDeviceModal = false)}
    on:success={handleDeviceAdded}
  />
{/if}
