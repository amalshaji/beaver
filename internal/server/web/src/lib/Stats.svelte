<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import { tunnelUserConnectionStatus } from "./store";

  let active_connections: number = 0,
    cpu_used: number = 0,
    memory_used: number = 0;

  let interval;

  const getStats = async () => {
    const res = await fetch("/api/v1/stats");
    const data: IStats = await res.json();

    active_connections = data.active_connections;
    memory_used = data.memory_used;
    cpu_used = data.cpu_used[0];

    if (
      JSON.stringify($tunnelUserConnectionStatus) !=
      JSON.stringify(data.connection_status)
    ) {
      $tunnelUserConnectionStatus = data.connection_status;
    }
  };

  onMount(() => {
    getStats();
    interval = setInterval(() => {
      getStats();
    }, 5000);
  });

  onDestroy(() => {
    clearInterval(interval);
  });
</script>

<div class="px-4 mt-6 sm:px-6 lg:px-8">
  <h2 class="text-gray-500 text-xs font-medium uppercase tracking-wide">
    Server Stats
  </h2>
  <ul
    role="list"
    class="grid grid-cols-1 gap-4 sm:gap-6 sm:grid-cols-2 xl:grid-cols-4 mt-3"
  >
    <li class="relative col-span-1 flex shadow-sm rounded-md">
      <div
        class="flex-shrink-0 flex items-center justify-center w-16 bg-pink-600 text-white text-sm font-medium rounded-l-md"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
          class="w-6 h-6"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M8.288 15.038a5.25 5.25 0 017.424 0M5.106 11.856c3.807-3.808 9.98-3.808 13.788 0M1.924 8.674c5.565-5.565 14.587-5.565 20.152 0M12.53 18.22l-.53.53-.53-.53a.75.75 0 011.06 0z"
          />
        </svg>
      </div>
      <div
        class="flex-1 flex items-center justify-between border-t border-r border-b border-gray-200 bg-white rounded-r-md truncate"
      >
        <div class="flex-1 px-4 py-2 text-sm">
          <a href="#" class="text-gray-900 font-medium hover:text-gray-600">
            Active Connections
          </a>
          <p class="text-gray-500">{active_connections}</p>
        </div>
      </div>
    </li>

    <li class="relative col-span-1 flex shadow-sm rounded-md">
      <div
        class="flex-shrink-0 flex items-center justify-center w-16 bg-emerald-600 text-white text-sm font-medium rounded-l-md"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
          class="w-6 h-6"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M8.25 3v1.5M4.5 8.25H3m18 0h-1.5M4.5 12H3m18 0h-1.5m-15 3.75H3m18 0h-1.5M8.25 19.5V21M12 3v1.5m0 15V21m3.75-18v1.5m0 15V21m-9-1.5h10.5a2.25 2.25 0 002.25-2.25V6.75a2.25 2.25 0 00-2.25-2.25H6.75A2.25 2.25 0 004.5 6.75v10.5a2.25 2.25 0 002.25 2.25zm.75-12h9v9h-9v-9z"
          />
        </svg>
      </div>
      <div
        class="flex-1 flex items-center justify-between border-t border-r border-b border-gray-200 bg-white rounded-r-md truncate"
      >
        <div class="flex-1 px-4 py-2 text-sm">
          <a href="#" class="text-gray-900 font-medium hover:text-gray-600">
            CPU
          </a>
          <p class="text-gray-500">{cpu_used.toFixed(2)}%</p>
        </div>
      </div>
    </li>

    <li class="relative col-span-1 flex shadow-sm rounded-md">
      <div
        class="flex-shrink-0 flex items-center justify-center w-16 bg-sky-600 text-white text-sm font-medium rounded-l-md"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          stroke-width="1.5"
          stroke="currentColor"
          class="w-6 h-6"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M20.25 6.375c0 2.278-3.694 4.125-8.25 4.125S3.75 8.653 3.75 6.375m16.5 0c0-2.278-3.694-4.125-8.25-4.125S3.75 4.097 3.75 6.375m16.5 0v11.25c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125V6.375m16.5 0v3.75m-16.5-3.75v3.75m16.5 0v3.75C20.25 16.153 16.556 18 12 18s-8.25-1.847-8.25-4.125v-3.75m16.5 0c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125"
          />
        </svg>
      </div>
      <div
        class="flex-1 flex items-center justify-between border-t border-r border-b border-gray-200 bg-white rounded-r-md truncate"
      >
        <div class="flex-1 px-4 py-2 text-sm">
          <a href="#" class="text-gray-900 font-medium hover:text-gray-600">
            Memory
          </a>
          <p class="text-gray-500">{memory_used.toFixed(2)}%</p>
        </div>
      </div>
    </li>
  </ul>
</div>
