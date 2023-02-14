<script lang="ts">
  import { onMount } from "svelte";

  let active_connections: number = 0,
    cpu_used: number = 0,
    memory_used: number = 0;

  const getStats = async () => {
    const res = await fetch("/api/v1/stats");
    const data = await res.json();

    active_connections = data.active_connections;
    memory_used = data.memory_used;
    cpu_used = data.cpu_used[0];
  };

  onMount(() => {
    getStats();
    setInterval(() => {
      getStats();
    }, 5000);
  });
</script>

<div>
  <dl class="mt-5 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
    <div
      class="relative bg-white pt-5 px-4 sm:px-6 border shadow rounded-lg overflow-hidden"
    >
      <dt>
        <div class="absolute border rounded-md p-3">
          <!-- Heroicon name: outline/users -->
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
        <p class="ml-16 text-sm font-medium text-gray-500 truncate">
          Active Connections
        </p>
      </dt>
      <dd class="ml-16 flex items-baseline">
        <p class="text-2xl font-semibold text-gray-900">{active_connections}</p>
      </dd>
    </div>

    <div
      class="relative bg-white pt-5 px-4 sm:px-6 border shadow rounded-lg overflow-hidden"
    >
      <dt>
        <div class="absolute border rounded-md p-3">
          <!-- Heroicon name: outline/mail-open -->
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
        <p class="ml-16 text-sm font-medium text-gray-500 truncate">CPU Used</p>
      </dt>
      <dd class="ml-16 pb-6 flex items-baseline sm:pb-7">
        <p class="text-2xl font-semibold text-gray-900">
          {cpu_used.toFixed(2)}%
        </p>
      </dd>
    </div>

    <div
      class="relative bg-white pt-5 px-4 sm:px-6 border shadow rounded-lg overflow-hidden"
    >
      <dt>
        <div class="absolute border rounded-md p-3">
          <!-- Heroicon name: outline/cursor-click -->
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
        <p class="ml-16 text-sm font-medium text-gray-500 truncate">
          Memory Used
        </p>
      </dt>
      <dd class="ml-16 pb-6 flex items-baseline sm:pb-7">
        <p class="text-2xl font-semibold text-gray-900">
          {memory_used.toFixed(2)}%
        </p>
      </dd>
    </div>
  </dl>
</div>
