<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { appState } from "../state.svelte";
  import gsap from "gsap";

  let mainChartData = $state(Array(40).fill(0.4));
  let throughputVal = $state(2.8);
  let interval: any;

  $effect(() => {
    const mainLine = document.getElementById("mainChartLine");
    const gradStop1 = document.getElementById("gradStop1");
    const gradStop2 = document.getElementById("gradStop2");

    if (mainLine && gradStop1 && gradStop2) {
      const targetColor = appState.isOutageActive ? "#ef4444" : "#3b82f6";
      const targetOpacity = appState.isOutageActive ? "0.3" : "0.2";

      gsap.to(mainLine, {
        stroke: targetColor,
        duration: 0.6,
        ease: "power2.inOut",
      });

      gsap.to(gradStop1, {
        stopColor: targetColor,
        stopOpacity: targetOpacity,
        duration: 0.6,
        ease: "power2.inOut",
      });

      gsap.to(gradStop2, {
        stopColor: targetColor,
        duration: 0.6,
        ease: "power2.inOut",
      });
    }
  });

  function getPath(data: number[], width: number, height: number) {
    if (data.length === 0) return "";
    const step = width / (data.length - 1);
    let d = `M 0 ${height - data[0] * height}`;

    for (let i = 0; i < data.length - 1; i++) {
      const x0 = i * step;
      const y0 = height - data[i] * height;
      const x1 = (i + 1) * step;
      const y1 = height - data[i + 1] * height;
      const cpX = (x0 + x1) / 2;
      d += ` C ${cpX} ${y0}, ${cpX} ${y1}, ${x1} ${y1}`;
    }
    return d;
  }

  function updateSimulation() {
    const last = mainChartData[mainChartData.length - 1];
    let target = appState.isOutageActive ? 0.05 : 0.6;
    let next = last + (Math.random() - 0.5) * 0.15;

    next = next + (target - next) * 0.1;
    next = Math.max(0.01, Math.min(0.95, next));

    mainChartData.shift();
    mainChartData.push(next);
    mainChartData = [...mainChartData];

    const width = 1000;
    const height = 300;
    const linePath = getPath(mainChartData, width, height);
    const areaPath = `${linePath} L ${width} ${height} L 0 ${height} Z`;

    const mainLine = document.getElementById("mainChartLine");
    const mainArea = document.getElementById("mainChartArea");

    if (mainLine && mainArea) {
      gsap.to(mainLine, { attr: { d: linePath }, duration: 0.8, ease: "none" });
      gsap.to(mainArea, { attr: { d: areaPath }, duration: 0.8, ease: "none" });
    }

    throughputVal = parseFloat((next * 5).toFixed(1));
  }

  onMount(() => {
    interval = setInterval(updateSimulation, 1000);
    updateSimulation();
  });

  onDestroy(() => {
    if (interval) clearInterval(interval);
  });
</script>

<div class="glass-panel rounded-3xl p-1 relative group">
  <!-- Controlled Glow -->
  <div
    class="absolute -bottom-40 -right-20 w-[600px] h-[400px] blur-[80px] rounded-full pointer-events-none transition-opacity duration-1000 opacity-0 group-hover:opacity-100 {appState.isOutageActive
      ? 'bg-red-500/5 dark:bg-red-500/10'
      : 'bg-blue-500/5 dark:bg-blue-500/10'}"
  ></div>

  <div class="relative z-10 p-8">
    <!-- Chart Header -->
    <div class="flex justify-between items-start mb-2">
      <div>
        <div class="flex items-center gap-2 mb-8">
          <h2 class="text-sm font-medium text-zinc-900 dark:text-zinc-200">
            Throughput
          </h2>
        </div>
        <div>
          <div class="flex items-end gap-3">
            <span
              class="text-5xl font-medium text-zinc-900 dark:text-white tracking-tighter-custom"
              >{throughputVal}</span
            >
            <span
              class="text-lg text-zinc-500 dark:text-zinc-500 font-normal mb-1.5"
              >Gbps</span
            >
            <div
              class="flex items-center gap-1 mb-2 px-2 py-0.5 rounded-full text-xs font-medium border {appState.isOutageActive
                ? 'bg-red-50 dark:bg-red-500/10 border-red-100 dark:border-red-500/20 text-red-600 dark:text-red-400'
                : 'bg-blue-50 dark:bg-blue-500/10 border-blue-100 dark:border-blue-500/20 text-blue-600 dark:text-blue-400'}"
            >
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
                class="transition-transform duration-500 {appState.isOutageActive
                  ? 'rotate-90'
                  : ''}"><path d="M7 7h10v10" /><path d="M7 17 17 7" /></svg
              >
              <span>12.4%</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Filter Pills -->
      <div
        class="flex bg-zinc-100 dark:bg-zinc-900/50 rounded-lg p-1 border border-zinc-200 dark:border-white/5"
      >
        <button
          class="px-4 py-1 text-[11px] font-medium rounded-md text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-300 transition-colors"
          >24h</button
        >
        <button
          class="px-4 py-1 text-[11px] font-medium rounded-md bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white shadow-sm border border-zinc-200 dark:border-white/10"
          >Live</button
        >
        <button
          class="px-4 py-1 text-[11px] font-medium rounded-md text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-300 transition-colors"
          >Max</button
        >
      </div>
    </div>

    <!-- Chart Area -->
    <div class="h-[280px] w-full relative mt-4 overflow-hidden rounded-xl">
      <svg
        class="w-full h-full overflow-visible"
        preserveAspectRatio="none"
        viewBox="0 0 1000 300"
      >
        <defs>
          <linearGradient id="chartGradient" x1="0" y1="0" x2="0" y2="1">
            <stop
              id="gradStop1"
              offset="0%"
              stop-color={appState.isOutageActive ? "#ef4444" : "#3b82f6"}
              stop-opacity={appState.isOutageActive ? "0.3" : "0.2"}
              class="dark:stop-opacity-0.3"
            ></stop>
            <stop
              id="gradStop2"
              offset="100%"
              stop-color={appState.isOutageActive ? "#ef4444" : "#3b82f6"}
              stop-opacity="0"
            ></stop>
          </linearGradient>
          <clipPath id="chartClip">
            <rect width="1000" height="300" rx="0"></rect>
          </clipPath>
        </defs>

        <!-- Clean Grid -->
        <g
          class="opacity-10 dark:opacity-20"
          stroke="currentColor"
          stroke-dasharray="0"
          stroke-width="1"
        >
          <line x1="0" y1="299" x2="1000" y2="299"></line>
          <line x1="0" y1="225" x2="1000" y2="225" stroke-dasharray="4,4"
          ></line>
          <line x1="0" y1="150" x2="1000" y2="150" stroke-dasharray="4,4"
          ></line>
          <line x1="0" y1="75" x2="1000" y2="75" stroke-dasharray="4,4"></line>
        </g>

        <path
          id="mainChartArea"
          d=""
          fill="url(#chartGradient)"
          clip-path="url(#chartClip)"
        ></path>
        <path
          id="mainChartLine"
          d=""
          fill="none"
          stroke={appState.isOutageActive ? "#ef4444" : "#3b82f6"}
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          vector-effect="non-scaling-stroke"
        ></path>
      </svg>
    </div>
  </div>
</div>
