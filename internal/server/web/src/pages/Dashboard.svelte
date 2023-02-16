<script lang="ts">
  import Navbar from "../lib/Navbar.svelte";
  import Stats from "../lib/Stats.svelte";
  import TunnelUsers from "../lib/TunnelUsers.svelte";
</script>

<Navbar />

<div class="mx-1 md:mx-auto w-full md:w-3/4 pt-16 space-y-8">
  <Stats />
  <TunnelUsers />
</div>
