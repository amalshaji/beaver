<script lang="ts">
  import { onMount } from "svelte";
  import toast from "svelte-french-toast";
  import Loader from "./Loader.svelte";
  import { Tooltip } from "@svelte-plugins/tooltips";

  let tunnelUsers = [];
  let email;
  let loading = false;

  const generateSecretKeyMask = (input: string): string => {
    let maskedString = "";
    for (let i = 0; i < input.length; i++) {
      if (input[i] == "-") {
        maskedString += "-";
      } else {
        maskedString += "x";
      }
    }
    // pad strings to get rid of cls
    return maskedString;
  };

  const copyToClipboard = async (text) => {
    await navigator.clipboard.writeText(text);
  };

  const getTunnelUsers = async () => {
    const res = await fetch("/api/v1/tunnel-users");
    tunnelUsers = await res.json();
  };

  const createTunnelUser = async () => {
    loading = true;
    try {
      const res = await fetch("/api/v1/tunnel-users", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email }),
      });
      const data = await res.json();
      if (res.status == 200) {
        email = "";
        await getTunnelUsers();
        toast.success(`New user added: ${data.email}`);
      } else {
        toast.error(data.error);
      }
    } catch (err) {
      console.error(err);
    } finally {
      loading = false;
    }
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
      const data = await res.json();
      if (res.status == 200) {
        await getTunnelUsers();
        toast.success(`New tunnel SecretKey generated for: ${data.email}`);
      } else {
        toast.error(data.error);
      }
    } catch (err) {
      console.error(err);
    } finally {
    }
  };

  onMount(() => {
    getTunnelUsers();
  });
</script>

<div class="flex flex-col">
  <div class="grid grid-cols-1 md:grid-cols-2 my-4 space-y-2">
    <div>
      <h1 class="font-semibold text-xl my-3">Tunnel Users</h1>
      <h3 class="text-gray-700 text-sm">
        A list of all the registered tunnel users
      </h3>
    </div>
    <div>
      <form on:submit|preventDefault={createTunnelUser} class="mt-6 flex">
        <label for="email" class="sr-only">Email address</label>
        <input
          type="email"
          name="email"
          id="email"
          bind:value={email}
          class="shadow-sm focus:ring-gray-500 focus:border-gray-500 block w-full sm:text-sm border-gray-300 rounded-md"
          placeholder="Enter an email"
        />
        <button
          type="submit"
          class="ml-4 w-24 h-12 flex-shrink-0 px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-gray-800 hover:bg-gray-900 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
        >
          {#if loading}
            <Loader />
          {:else}
            Add user
          {/if}
        </button>
      </form>
    </div>
  </div>
  <div class="-my-2 overflow-x-auto w-full lg:w-5/6 mx-auto mt-6">
    <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
      <div
        class="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg"
      >
        {#if tunnelUsers.length > 0}
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th
                  scope="col"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                  >Email</th
                >
                <th
                  scope="col"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                  >Secret Key</th
                >
                <th scope="col" class="relative px-6 py-3">
                  <span class="sr-only">Rotate Secret Key</span>
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              {#each tunnelUsers as tunnelUser}
                <tr>
                  <td
                    class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900"
                    >{tunnelUser.email}</td
                  >
                  <!-- svelte-ignore a11y-click-events-have-key-events -->
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    <Tooltip content="click to copy">
                      <input
                        class="hover:cursor-pointer border-0 w-full text-sm focus:ring-0 px-0"
                        readonly
                        type="text"
                        value={generateSecretKeyMask(tunnelUser.secret_key)}
                        on:mouseenter={(e) =>
                          (e.target.value = tunnelUser.secret_key)}
                        on:mouseleave={(e) =>
                          (e.target.value = generateSecretKeyMask(
                            tunnelUser.secret_key
                          ))}
                        on:click={async () =>
                          await copyToClipboard(tunnelUser.secret_key)}
                      />
                    </Tooltip>
                  </td>
                  <td
                    class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium"
                  >
                    <button
                      on:click={() =>
                        rotateTunnelUserSecretKey(tunnelUser.email)}
                      class="text-gray-600 hover:text-gray-900 border rounded-lg px-2 py-1"
                      >Rotate Secret Key</button
                    >
                  </td>
                </tr>
              {/each}

              <!-- More people... -->
            </tbody>
          </table>
        {:else}
          <div
            class="relative block w-full border-2 border-gray-300 border-dashed rounded-lg p-12 text-center hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              class="mx-auto w-8 h-8 text-gray-400"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z"
              />
            </svg>

            <span class="mt-2 block text-sm font-semibold text-gray-900">
              No tunnel users
            </span>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
