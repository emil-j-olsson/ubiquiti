<script lang="ts">
  import { onMount } from "svelte";
  import gsap from "gsap";
  import type { Device, Diagnostics } from "../state.svelte";
  import { DeviceStatus } from "../state.svelte";

  let { device, diagnostics, history } = $props<{
    device: Device;
    diagnostics: Diagnostics | undefined;
    history: { timestamp: number; cpu: number; memory: number }[];
  }>();

  let card: HTMLElement;
  let cpuPath = $state("");

  let cpuValue = $derived(diagnostics ? Math.round(diagnostics.cpu_usage) : 0);
  let memValue = $derived(
    diagnostics ? Math.round(diagnostics.memory_usage) : 0
  );
  let statusText = $derived(getStatusText(diagnostics?.device_status));
  let isHealthy = $derived(
    diagnostics?.device_status === DeviceStatus.DEVICE_STATUS_HEALTHY
  );

  $effect(() => {
    if (history.length > 0) {
      const width = 100;
      const height = 32;
      const data = history.map((h: { cpu: number }) => h.cpu / 100);

      if (data.length === 0) {
        cpuPath = "";
      } else {
        const step = width / Math.max(data.length - 1, 1);
        let d = `M 0 ${height - data[0] * height}`;

        for (let i = 0; i < data.length - 1; i++) {
          const x0 = i * step;
          const y0 = height - data[i] * height;
          const x1 = (i + 1) * step;
          const y1 = height - data[i + 1] * height;
          const cpX = (x0 + x1) / 2;
          d += ` C ${cpX} ${y0}, ${cpX} ${y1}, ${x1} ${y1}`;
        }
        cpuPath = d;
      }
    }
  });

  function getStatusText(status: DeviceStatus | undefined) {
    if (status === undefined) return "Unknown";
    switch (status) {
      case DeviceStatus.DEVICE_STATUS_HEALTHY:
        return "Operational";
      case DeviceStatus.DEVICE_STATUS_DEGRADED:
        return "Degraded";
      case DeviceStatus.DEVICE_STATUS_ERROR:
        return "Error";
      case DeviceStatus.DEVICE_STATUS_MAINTENANCE:
        return "Maintenance";
      case DeviceStatus.DEVICE_STATUS_BOOTING:
        return "Booting";
      case DeviceStatus.DEVICE_STATUS_OFFLINE:
        return "Offline";
      default:
        return "Unknown";
    }
  }

  onMount(() => {
    // Entrance animation
    gsap.fromTo(
      card,
      { opacity: 0, y: 10 },
      { opacity: 1, y: 0, duration: 0.6, ease: "power2.out" }
    );
  });
</script>

<div
  bind:this={card}
  class="device-card glass-panel rounded-xl p-6 relative group transition-all duration-500 {isHealthy
    ? 'hover:border-blue-500/30 dark:hover:border-blue-500/30'
    : 'hover:border-red-500/20'}"
>
  <!-- Header -->
  <div class="flex justify-between items-start mb-6 relative z-10">
    <div class="flex gap-4 items-center">
      <div
        class="w-12 h-12 rounded-lg bg-zinc-50 dark:bg-zinc-900 border border-zinc-200 dark:border-white/10 flex items-center justify-center text-zinc-400 dark:text-zinc-500 group-hover:text-zinc-300 dark:group-hover:text-zinc-300 transition-colors duration-300"
      >
        {#if device.host.toLowerCase().includes("switch") || device.host
            .toLowerCase()
            .includes("router")}
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="35"
            height="28"
            viewBox="0 0 50 40"
            class="device-icon transition-transform duration-500 group-hover:scale-110"
            ><path
              fill="currentColor"
              fill-rule="evenodd"
              d="M7.33 12.746c-.322 0-.641.064-.939.188l-1.516.632c-.205.085-.144.39.078.39h40.094c.222 0 .283-.305.078-.39l-1.517-.632a2.4 2.4 0 0 0-.938-.188zM4.343 14.2a.407.407 0 0 0-.406.406v3.303c0 .224.182.406.406.406h41.313a.407.407 0 0 0 .407-.406v-3.303a.407.407 0 0 0-.407-.406zm.727.968a.407.407 0 0 0-.407.407v1.366c0 .224.182.406.407.406h1.365a.407.407 0 0 0 .407-.406v-1.366a.407.407 0 0 0-.407-.407zM2.407 14.2a.407.407 0 0 0-.407.406v3.303c0 .224.182.406.407.406h.881a.407.407 0 0 0 .407-.406v-3.303a.407.407 0 0 0-.407-.406zm.44.484a.363.363 0 1 0 0 .726.363.363 0 0 0 0-.726m-.363 2.784a.363.363 0 1 1 .727 0 .363.363 0 0 1-.727 0M46.712 14.2a.407.407 0 0 0-.407.406v3.303c0 .224.182.406.407.406h.881a.407.407 0 0 0 .407-.406v-3.303a.407.407 0 0 0-.407-.406zm.44.484a.363.363 0 1 0 0 .726.363.363 0 0 0 0-.726m-.362 2.784a.363.363 0 1 1 .726 0 .363.363 0 0 1-.726 0m-6.701-2.542H36.96v2.663h7.487v-2.663zm1.452.407h1.211v.804h-1.21zm-.242.804v-.804h-1.21v.804zm-1.21.242h1.21v.804h-1.21zm-.242-.242v-.804h-1.14v.804zm-1.14.242h1.14v.804h-1.14zm-.242-.242v-.804h-1.098v.804zm-1.098.242h1.098v.804h-1.098zm4.174 0h1.211v.804h-1.21zm1.453 0v.804h1.046v-.804zm1.046-.242v-.804h-1.046v.804zm-39.697 2.53a.407.407 0 0 0-.406.406v3.303c0 .225.182.407.406.407h41.313a.407.407 0 0 0 .407-.407v-3.303a.407.407 0 0 0-.407-.406zm.727.968a.407.407 0 0 0-.407.407v1.366c0 .224.182.406.407.406h1.365a.407.407 0 0 0 .407-.406V20.04a.407.407 0 0 0-.407-.407zm-2.663-.968a.407.407 0 0 0-.407.406v3.303c0 .225.182.407.407.407h.881a.407.407 0 0 0 .407-.407v-3.303a.407.407 0 0 0-.407-.406zm.44.484a.363.363 0 1 0 0 .726.363.363 0 0 0 0-.726m-.363 2.784a.363.363 0 1 1 .727 0 .363.363 0 0 1-.727 0m44.228-3.268a.407.407 0 0 0-.407.406v3.303c0 .225.182.407.407.407h.881a.407.407 0 0 0 .407-.407v-3.303a.407.407 0 0 0-.407-.406zm.44.484a.363.363 0 1 0 0 .727.363.363 0 0 0 0-.727m-.362 2.784a.363.363 0 1 1 .726 0 .363.363 0 0 1-.726 0M4.343 23.138a.407.407 0 0 0-.406.406v3.303c0 .224.182.406.406.406h41.313a.407.407 0 0 0 .407-.406v-3.303a.407.407 0 0 0-.407-.406zm.727.968a.407.407 0 0 0-.407.407v1.365c0 .225.182.407.407.407h1.365a.407.407 0 0 0 .407-.407v-1.365a.407.407 0 0 0-.407-.407zm-2.663-.968a.407.407 0 0 0-.407.406v3.303c0 .224.182.406.407.406h.881a.407.407 0 0 0 .407-.406v-3.303a.407.407 0 0 0-.407-.406zm.44.484a.363.363 0 1 0 0 .726.363.363 0 0 0 0-.726m-.363 2.784a.363.363 0 1 1 .727 0 .363.363 0 0 1-.727 0m44.228-3.268a.407.407 0 0 0-.407.406v3.303c0 .224.182.406.407.406h.881a.407.407 0 0 0 .407-.406v-3.303a.407.407 0 0 0-.407-.406zm.44.484a.363.363 0 1 0 0 .726.363.363 0 0 0 0-.726m-.362 2.784a.363.363 0 1 1 .726 0 .363.363 0 0 1-.726 0m-11.457-2.542h-3.128v2.663h7.486v-2.663zm1.452.406h1.211v.805h-1.21zm-.242.805v-.805h-1.21v.805zm-1.21.242h1.21v.804h-1.21zm-.242-.242v-.805h-1.14v.805zm-1.14.242h1.14v.804h-1.14zm-.242-.242v-.805H32.61v.805zm-1.098.242h1.098v.804H32.61zm4.174 0h1.211v.804h-1.21zm1.453 0v.804h1.046v-.804zm1.046-.242v-.805h-1.046v.805zm-12.323-.805h-1.14v.805h1.14zm.242.805v-.805h1.21v.805zm0 .242v.804h1.21v-.804zm-.242 0v.804h-1.14v-.804zm-1.383 0v.804h-1.097v-.804zm0-.242v-.805h-1.097v.805zm3.077.242v.804h1.21v-.804zm0-.242v-.805h1.21v.805zm1.453 0v-.805h1.046v.805zm0 .242v.804h1.046v-.804zm-6.034-1.453h7.486v2.663h-7.486zm-5.244.406h-1.14v.805h1.14zm.242.805v-.805h1.21v.805zm0 .242v.804h1.21v-.804zm-.242 0v.804h-1.14v-.804zm-1.382 0v.804H16.35v-.804zm0-.242v-.805H16.35v.805zm3.076.242v.804h1.211v-.804zm0-.242v-.805h1.211v.805zm1.453 0v-.805h1.046v.805zm0 .242v.804h1.046v-.804zm-6.033-1.453h7.486v2.663h-7.486zm-5.244.406H9.56v.805h1.14zm.242.805v-.805h1.21v.805zm0 .242v.804h1.21v-.804zm-.242 0v.804H9.56v-.804zm-1.383 0v.804H8.22v-.804zm0-.242v-.805H8.22v.805zm3.077.242v.804h1.21v-.804zm0-.242v-.805h1.21v.805zm1.453 0v-.805h1.046v.805zm0 .242v.804h1.046v-.804zm-6.034-1.453H15.3v2.663H7.813zm35.907.406h-.882v.805h.882zm.406.805v1.452h-1.694v-2.663h1.694zm-1.288.242v.803h.882v-.803zm-1.15-1.047h-.883v.805h.882zm.406.805v1.452h-1.695v-2.663h1.695zm-1.289.242v.803h.882v-.803zM9.336 20.465H8.225v.407h1.111zm.242 0v.407h1.156v-.407zm2.591.407h-1.193v-.407h1.193zm.242 0h1.207v-.407H12.41zm2.496 0H13.86v-.407h1.047zm-7.089-.813v1.22h7.496v-1.22zm8.525.406h1.11v.407h-1.11zm1.353.407v-.407h1.155v.407zm1.397 0h1.194v-.407h-1.194zm2.642 0H20.53v-.407h1.206zm.242 0h1.048v-.407h-1.047zm-6.04.407v-1.22h7.494v1.22zm9.635-.814H24.46v.407h1.112zm.242 0v.407h1.155v-.407zm2.59.407h-1.193v-.407h1.194zm.243 0h1.206v-.407h-1.206zm2.496 0h-1.048v-.407h1.047zm-7.09-.813v1.22h7.496v-1.22zm8.55.406h1.112v.407h-1.112zm1.354.407v-.407h1.155v.407zm1.397 0h1.194v-.407h-1.194zm2.642 0H36.79v-.407h1.206zm.242 0h1.048v-.407h-1.048zm-6.041.407v-1.22h7.495v1.22zm9.423-.814h-.813v.407h.813zm-1.22-.406v1.22h1.626v-1.22zm3.252.406h-.813v.407h.813zm-1.22-.406v1.22h1.627v-1.22z"
              clip-rule="evenodd"
            /></svg
          >
        {:else}
          <svg
            width="22"
            height="22"
            viewBox="0 0 138 138"
            fill="currentColor"
            color="currentColor"
            xmlns="http://www.w3.org/2000/svg"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="device-icon transition-transform duration-500 group-hover:scale-110"
          >
            <path
              d="M69 0C107.108 0 138 30.8924 138 69C138 107.108 107.108 138 69 138C30.8924 138 0 107.108 0 69C0 30.8924 30.8924 0 69 0ZM69 38.8125C52.3279 38.8125 38.8125 52.3279 38.8125 69C38.8125 85.6721 52.3279 99.1875 69 99.1875C85.6721 99.1875 99.1875 85.6721 99.1875 69C99.1875 52.3279 85.6721 38.8125 69 38.8125ZM69 45C82.2548 45 93 55.7452 93 69C93 82.2548 82.2548 93 69 93C55.7452 93 45 82.2548 45 69C45 55.7452 55.7452 45 69 45Z"
            />
          </svg>
        {/if}
      </div>
      <div>
        <h3
          class="device-name text-base font-medium text-zinc-900 dark:text-zinc-200 tracking-tight"
        >
          {device.alias || device.device_id}
        </h3>
        <div class="flex items-center gap-2 mt-0.5">
          <span
            class="device-ip text-[10px] text-zinc-400 dark:text-zinc-500 font-mono tracking-wide"
          >
            {device.host}
          </span>
        </div>
      </div>
    </div>

    <!-- Status with Tooltip -->
    <div class="relative flex items-center group/status cursor-help">
      <span class="relative flex h-2 w-2">
        <span
          class="animate-ping absolute inline-flex h-full w-full rounded-full opacity-75 {isHealthy
            ? 'bg-blue-500'
            : 'bg-red-500'}"
        ></span>
        <span
          class="status-indicator relative inline-flex rounded-full h-2 w-2 ring-2 ring-white dark:ring-zinc-900 transition-all duration-300 {isHealthy
            ? 'bg-blue-500 shadow-[0_0_8px_rgba(59,130,246,0.5)]'
            : 'bg-red-500 shadow-[0_0_8px_rgba(239,68,68,0.6)]'}"
        ></span>
      </span>
      <div
        class="absolute right-full top-1/2 -translate-y-1/2 mr-3 px-2 py-1 bg-zinc-900 dark:bg-white text-white dark:text-black text-[10px] font-medium rounded shadow-lg opacity-0 -translate-x-1 group-hover/status:opacity-100 group-hover/status:translate-x-0 transition-all duration-200 whitespace-nowrap z-20 pointer-events-none"
      >
        {statusText}
      </div>
    </div>
  </div>

  <!-- Mini Sparklines -->
  <div class="grid grid-cols-2 gap-8 relative z-10">
    <!-- CPU -->
    <div>
      <div class="flex justify-between items-center mb-2">
        <span
          class="text-[10px] text-zinc-400 dark:text-zinc-500 uppercase tracking-wider font-medium"
          >CPU</span
        >
        <div class="flex items-center gap-1">
          <span
            class="cpu-value text-xs font-medium text-zinc-700 dark:text-zinc-300 font-mono"
            >{cpuValue}%</span
          >
        </div>
      </div>
      <div class="h-8 w-full">
        <svg
          class="w-full h-full overflow-visible"
          preserveAspectRatio="none"
          viewBox="0 0 100 32"
        >
          <path
            class="cpu-chart-line"
            fill="none"
            stroke={isHealthy
              ? document.documentElement.classList.contains("dark")
                ? "#52525b"
                : "#a1a1aa"
              : "#ef4444"}
            stroke-width="1.5"
            d={cpuPath}
          ></path>
        </svg>
      </div>
    </div>

    <!-- MEM -->
    <div>
      <div class="flex justify-between items-center mb-2">
        <span
          class="text-[10px] text-zinc-400 dark:text-zinc-500 uppercase tracking-wider font-medium"
          >Memory</span
        >
        <div class="flex items-center gap-1">
          <span
            class="mem-value text-xs font-medium text-zinc-700 dark:text-zinc-300 font-mono"
            >{memValue}%</span
          >
        </div>
      </div>
      <div class="h-8 w-full flex items-end gap-[3px] mem-bar-container">
        {#each Array(15) as _, i}
          {@const val =
            history[Math.max(0, history.length - 15 + i)]?.memory || 0}
          {@const h = Math.min(100, Math.max(5, val))}
          <div
            class="flex-1 bg-zinc-100 dark:bg-zinc-800 rounded-[1px] h-full relative overflow-hidden"
          >
            <div
              class="absolute bottom-0 left-0 w-full transition-all duration-700 {isHealthy
                ? 'bg-blue-500/80 dark:bg-blue-600'
                : 'bg-red-500'}"
              style="height: {h}%"
            ></div>
          </div>
        {/each}
      </div>
    </div>
  </div>

  <!-- Hidden Details Reveal -->
  <div class="reveal-details border-t border-zinc-100 dark:border-white/5 pt-4">
    <div class="grid grid-cols-3 gap-2">
      <div>
        <div class="text-[9px] text-zinc-400 uppercase tracking-wider mb-0.5">
          Software
        </div>
        <div class="text-[10px] font-mono text-zinc-600 dark:text-zinc-300">
          {diagnostics?.software_version || "N/A"}
        </div>
      </div>
      <div>
        <div class="text-[9px] text-zinc-400 uppercase tracking-wider mb-0.5">
          Firmware
        </div>
        <div class="text-[10px] font-mono text-zinc-600 dark:text-zinc-300">
          {diagnostics?.firmware_version || "N/A"}
        </div>
      </div>
      <div>
        <div class="text-[9px] text-zinc-400 uppercase tracking-wider mb-0.5">
          Hardware
        </div>
        <div class="text-[10px] font-mono text-zinc-600 dark:text-zinc-300">
          {diagnostics?.hardware_version || "N/A"}
        </div>
      </div>
    </div>
  </div>
</div>
