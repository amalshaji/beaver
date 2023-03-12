<script lang="ts">
  import { onDestroy, onMount } from "svelte";
  import toast from "svelte-french-toast";
  import moment from "moment";

  import AddTunnelUser from "../lib/modals/AddTunnelUser.svelte";
  import ShowSecretKey from "./modals/ShowSecretKey.svelte";
  import DeleteUser from "./modals/DeleteUser.svelte";

  import { tunnelUserConnectionStatus } from "./store";

  let addTunnelUserModalOpen = false;

  let tunnelUsers: ITunnelUser[] = [];

  let secretKey = undefined;
  let showSecretKeyModalOpen = false;

  let tunnelUserToDelete;
  let deleteUserModalOpen = false;

  const deleteUser = (id: number, email: string) => {
    tunnelUserToDelete = {
      ID: id,
      Email: email,
    };
    deleteUserModalOpen = true;
  };

  const getTunnelUsers = async () => {
    const res = await fetch("/api/v1/tunnel-users");
    tunnelUsers = await res.json();
  };

  const rotateTunnelUserSecretKey = async (email: string) => {
    try {
      const res = await fetch("/api/v1/tunnel-users", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email }),
      });
      if (res.status == 200) {
        const data = await res.json();
        secretKey = data.SecretKey;
        showSecretKeyModalOpen = true;
      } else {
        const data: IError = await res.json();
        toast.error(data.error);
      }
    } catch (err) {
      console.error(err);
    } finally {
    }
  };

  const unsubscribe = tunnelUserConnectionStatus.subscribe((n) => {
    if (tunnelUsers.length === 0) {
      return;
    }
    const obj = {};

    for (const item of n) {
      obj[item.ID] = {
        Active: item.Active,
        LastActiveAt: item.LastActiveAt,
      };
    }
    const tempTunnelUsers: ITunnelUser[] = [];
    for (const tunnelUser of tunnelUsers) {
      tunnelUser.Active = obj[tunnelUser.ID].Active;
      tunnelUser.LastActiveAt = obj[tunnelUser.ID].LastActiveAt;
      tempTunnelUsers.push(tunnelUser);
    }

    tunnelUsers = [...tempTunnelUsers];
  });

  onMount(() => {
    getTunnelUsers();
  });

  onDestroy(() => {
    unsubscribe();
  });
</script>

<AddTunnelUser
  isOpen={addTunnelUserModalOpen}
  onClose={() => {
    addTunnelUserModalOpen = false;
  }}
  on:success={(e) => {
    secretKey = e.detail;
    showSecretKeyModalOpen = true;
    getTunnelUsers();
  }}
/>

<ShowSecretKey
  isOpen={showSecretKeyModalOpen}
  {secretKey}
  onClose={() => {
    showSecretKeyModalOpen = false;
    secretKey = undefined;
  }}
/>

<DeleteUser
  userToDelete={tunnelUserToDelete}
  isOpen={deleteUserModalOpen}
  onClose={() => {
    deleteUserModalOpen = false;
    tunnelUserToDelete = undefined;
    getTunnelUsers();
  }}
/>

<!-- Tunnel Users -->
<div class="mt-10 sm:hidden">
  <div class="px-4 sm:px-6">
    <h2 class="text-gray-500 text-xs font-medium uppercase tracking-wide">
      Tunnel Users
    </h2>
  </div>
  <ul
    role="list"
    class="mt-3 border-t border-gray-200 divide-y divide-gray-100"
  >
    <li>
      <a
        href="#"
        class="group flex items-center justify-between px-4 py-4 hover:bg-slate-700 text-white sm:px-6"
      >
        <span class="flex items-center truncate space-x-3">
          <span
            class="w-2.5 h-2.5 flex-shrink-0 rounded-full bg-rose-600"
            aria-hidden="true"
          />
          <span class="font-medium truncate text-sm leading-6">
            GraphQL API
            <span class="truncate font-normal text-gray-500"
              >in Engineering</span
            >
          </span>
        </span>
        <!-- Heroicon name: solid/chevron-right -->
        <svg
          class="ml-4 h-5 w-5 text-gray-400 group-hover:text-gray-500"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 20 20"
          fill="currentColor"
          aria-hidden="true"
        >
          <path
            fill-rule="evenodd"
            d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
            clip-rule="evenodd"
          />
        </svg>
      </a>
    </li>

    <!-- More projects... -->
  </ul>
</div>

<!-- Tunnel users table (small breakpoint and up) -->
<div class="hidden mt-8 sm:block mb-16 sm:px-6 lg:px-8">
  <div class="sm:flex sm:items-center w-full py-4">
    <h2 class="text-gray-500 text-xs font-medium uppercase tracking-wide py-3">
      Tunnel Users
    </h2>
    <div class="mt-4 sm:mt-0 sm:ml-4 sm:flex-none">
      <button
        type="button"
        on:click={() => (addTunnelUserModalOpen = true)}
        class="my-auto float-right inline-flex items-center justify-center rounded-sm border border-transparent bg-gray-600 px-2 py-1 text-sm font-medium text-white shadow-sm hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 sm:w-auto"
        >Add user</button
      >
    </div>
  </div>
  {#if tunnelUsers.length > 0}
    <div class="align-middle inline-block w-full border">
      <table class="w-full table-fixed rounded-lg">
        <thead class="rounded-lg">
          <tr class="border-t border-gray-200">
            <th
              class="px-6 py-3 border-b border-gray-200 bg-zinc-500 text-white text-left text-xs font-medium uppercase tracking-wider"
            >
              <span class="lg:pl-8">Email</span>
            </th>
            <th
              class="px-6 py-3 border-b border-gray-200 bg-zinc-500 text-white text-left text-xs font-medium uppercase tracking-wider"
              >Last Active</th
            >
            <th
              class="pr-6 py-3 border-b border-gray-200 bg-zinc-500 text-white text-right text-xs font-medium uppercase tracking-wider"
            />
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-100">
          {#each tunnelUsers as tunnelUser}
            <tr>
              <td
                class="px-6 py-3 max-w-0 w-full whitespace-nowrap text-sm font-medium text-gray-900"
              >
                <div class="flex items-center space-x-3 lg:pl-2">
                  <div
                    class="flex-shrink-0 w-2.5 h-2.5 rounded-full {tunnelUser.Active
                      ? 'bg-green-600 backdrop-blur-lg'
                      : 'bg-gray-100 border border-gray-300'}"
                    aria-hidden="true"
                  >
                    {#if tunnelUser.Active}
                      <div class="w-2 h-2 bg-green-600 blur-sm" />
                    {/if}
                  </div>
                  <p class="truncate hover:text-gray-600">
                    <span>
                      {tunnelUser.Email}
                    </span>
                  </p>
                </div>
              </td>
              <td
                title={tunnelUser.LastActiveAt === null
                  ? "Not available"
                  : tunnelUser.Active
                  ? "Online"
                  : moment(tunnelUser.LastActiveAt).format(
                      "MMMM Do YYYY, h:mm:ss a"
                    )}
                class="px-6 py-3 whitespace-nowrap text-sm text-gray-500 text-left"
              >
                {#if tunnelUser.Active}
                  Online
                {:else}
                  {tunnelUser.LastActiveAt === null
                    ? "Not available"
                    : moment(tunnelUser.LastActiveAt).from(new Date())}
                {/if}
              </td>
              <td
                class="px-6 py-3 whitespace-nowrap text-right text-sm font-medium sm:space-x-4"
              >
                <button
                  class="text-indigo-600 hover:text-indigo-900"
                  on:click={() => rotateTunnelUserSecretKey(tunnelUser.Email)}
                  >Rotate Key</button
                >
                <button
                  class="text-red-600 hover:text-red-900"
                  on:click={() => deleteUser(tunnelUser.ID, tunnelUser.Email)}
                  >Delete</button
                >
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {:else}
    <button
      on:click={() => (addTunnelUserModalOpen = true)}
      type="button"
      class="relative block w-full border-2 border-gray-300 border-dashed rounded-lg p-12 text-center hover:border-gray-400 focus:outline-none"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke-width="1.5"
        stroke="currentColor"
        class="mx-auto h-12 w-12 text-gray-400"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z"
        />
      </svg>

      <span class="mt-2 block text-sm font-medium text-gray-900">
        Create a new tunnel user
      </span>
    </button>
  {/if}
</div>
