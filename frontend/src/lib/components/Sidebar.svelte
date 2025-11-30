<script lang="ts">
  import { page } from "$app/stores";
  import { appState, DeviceStatus } from "../state.svelte";
  import { onMount } from "svelte";
  import { fade } from "svelte/transition";

  let isOutageActive = $state(false);

  async function toggleOutage() {
    if (!appState.client) return;

    try {
      if (!appState.isOutageActive) {
        // Simulate outage - set all devices to degraded
        const updatePromises = appState.devices.map((device) =>
          appState.client!.updateDevice(
            device.device_id,
            DeviceStatus.DEVICE_STATUS_DEGRADED
          )
        );
        await Promise.all(updatePromises);
      } else {
        // Restore to healthy - set all devices back to healthy
        const updatePromises = appState.devices.map((device) =>
          appState.client!.updateDevice(
            device.device_id,
            DeviceStatus.DEVICE_STATUS_HEALTHY
          )
        );
        await Promise.all(updatePromises);
      }
    } catch (error) {
      console.error("Failed to toggle outage:", error);
    }
  }
</script>

<aside
  class="fixed left-4 md:left-6 top-1/2 -translate-y-1/2 z-50 flex flex-col gap-4 pointer-events-none h-[85vh] justify-center"
>
  <nav
    class="glass-nav rounded-2xl p-2 flex flex-col items-center gap-6 pointer-events-auto animate-fade-in"
  >
    <!-- Logo -->
    <div
      class="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center shadow-lg shadow-blue-500/20 text-white mb-2 cursor-pointer hover:scale-95 transition-transform duration-300"
    >
      <svg
        class="pt-0.5"
        xmlns="http://www.w3.org/2000/svg"
        width="20"
        height="20"
        fill="none"
        viewBox="0 0 128 129"
        ><path
          fill="currentColor"
          d="M123.998 0h-7.999v7.999h7.999V0ZM96.01 56.02V39.994l.004.005h15.993v15.997H128v5.07c0 5.862-.498 12.807-1.644 18.26-.642 3.049-1.615 6.078-2.756 8.987-1.169 2.976-2.516 5.827-3.982 8.382a62.776 62.776 0 0 1-6.695 9.547l-.136.158-.224.262c-.617.724-1.227 1.44-1.898 2.139a63.215 63.215 0 0 1-2.415 2.415c-10.237 9.859-23.575 16.017-37.521 17.431-1.678.172-5.047.35-6.729.35-1.687-.005-5.051-.178-6.729-.35-13.946-1.414-27.284-7.577-37.52-17.431a63.155 63.155 0 0 1-2.416-2.415c-.704-.729-1.339-1.477-1.98-2.233l-.003-.002-.275-.324a62.658 62.658 0 0 1-6.695-9.547c-1.466-2.56-2.813-5.406-3.982-8.382-1.141-2.91-2.114-5.938-2.756-8.986C.498 73.867 0 66.928 0 61.067V1.002h31.99V56.02s0 4.218.053 5.598l.012.322v.002c.068 1.785.133 3.534.319 5.274.527 4.94 1.62 9.628 3.872 13.592.652 1.145 1.313 2.257 2.104 3.311 4.812 6.417 12.135 11.234 21.27 12.576 1.087.158 3.283.297 4.38.297s3.292-.139 4.38-.297c9.135-1.342 16.458-6.159 21.27-12.576.795-1.054 1.452-2.166 2.104-3.311 2.252-3.964 3.345-8.651 3.872-13.592.186-1.743.252-3.495.319-5.284l.012-.314c.053-1.38.053-5.598.053-5.598ZM100.002 12h12v11.996H128v15.998h-15.998V24.001h-12v-12Z"
        /></svg
      >
    </div>

    <!-- Main Links -->
    <div class="flex flex-col gap-2 w-full items-center">
      <a
        href="/"
        class="relative group p-2.5 rounded-xl transition-all {$page.url
          .pathname === '/'
          ? 'text-zinc-900 dark:text-white bg-black/5 dark:bg-white/10'
          : 'text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-200 hover:bg-black/5 dark:hover:bg-white/5'}"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          ><rect width="7" height="7" x="3" y="3" rx="1" /><rect
            width="7"
            height="7"
            x="14"
            y="3"
            rx="1"
          /><rect width="7" height="7" x="14" y="14" rx="1" /><rect
            width="7"
            height="7"
            x="3"
            y="14"
            rx="1"
          /></svg
        >
        <div
          class="tooltip absolute left-full top-1/2 -translate-y-1/2 ml-4 px-2.5 py-1.5 bg-zinc-900 text-white text-[11px] font-medium rounded-md shadow-xl group-hover:opacity-100 group-hover:translate-x-0 opacity-0 -translate-x-2 transition-all duration-200 pointer-events-none whitespace-nowrap z-50 border border-white/10"
        >
          Device Overview
        </div>
      </a>
    </div>

    <div class="w-6 h-px bg-zinc-200 dark:bg-white/10 my-1"></div>

    <!-- Actions -->
    <div class="flex flex-col gap-2 w-full items-center">
      <button
        id="outageToggleBtn"
        class="relative group p-2.5 rounded-xl text-zinc-400 transition-all cursor-pointer {appState.isOutageActive
          ? 'text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-500/10 hover:text-red-300 hover:bg-red-500/20'
          : 'hover:text-blue-600 dark:hover:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-500/10'}"
        onclick={toggleOutage}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="lucide lucide-zap-off-icon lucide-zap-off"
          ><path
            d="M10.513 4.856 13.12 2.17a.5.5 0 0 1 .86.46l-1.377 4.317"
          /><path d="M15.656 10H20a1 1 0 0 1 .78 1.63l-1.72 1.773" /><path
            d="M16.273 16.273 10.88 21.83a.5.5 0 0 1-.86-.46l1.92-6.02A1 1 0 0 0 11 14H4a1 1 0 0 1-.78-1.63l4.507-4.643"
          /><path d="m2 2 20 20" /></svg
        >
        <div
          class="tooltip absolute left-full top-1/2 -translate-y-1/2 ml-4 px-2.5 py-1.5 bg-zinc-900 text-white text-[11px] font-medium rounded-md shadow-xl opacity-0 group-hover:opacity-100 group-hover:translate-x-0 -translate-x-2 transition-all duration-200 pointer-events-none whitespace-nowrap z-50 border border-white/10"
        >
          Simulate Outage
        </div>
      </button>

      <button
        id="addDeviceBtn"
        class="relative group p-2.5 rounded-xl text-zinc-400 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-500/10 transition-all cursor-pointer"
        onclick={() =>
          window.dispatchEvent(new CustomEvent("openAddDeviceModal"))}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          ><path d="M5 12h14" /><path d="M12 5v14" /></svg
        >
        <div
          class="tooltip absolute left-full top-1/2 -translate-y-1/2 ml-4 px-2.5 py-1.5 bg-zinc-900 text-white text-[11px] group-hover:opacity-100 group-hover:translate-x-0 font-medium rounded-md shadow-xl opacity-0 -translate-x-2 transition-all duration-200 pointer-events-none whitespace-nowrap z-50 border border-white/10"
        >
          Add Device
        </div>
      </button>
    </div>

    <div class="flex-grow"></div>

    <!-- Profile -->
    <div class="relative group cursor-pointer mb-2">
      <div
        class="w-9 h-9 rounded-full bg-zinc-100 dark:bg-zinc-800 border border-zinc-200 dark:border-white/10 flex items-center justify-center overflow-hidden hover:border-zinc-300 dark:hover:border-white/30 transition-colors"
      >
        <span class="text-xs font-semibold text-zinc-600 dark:text-zinc-300"
          >JD</span
        >
      </div>
    </div>
  </nav>
</aside>
