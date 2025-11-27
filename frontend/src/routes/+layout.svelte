<script lang="ts">
  import "./layout.css";
  import favicon from "$lib/assets/favicon-320x320.png";

  let { children } = $props();
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
  <title>Ubiquiti Device Monitoring</title>
</svelte:head>

{@render children()}
