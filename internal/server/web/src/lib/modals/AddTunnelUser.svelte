<script lang="ts">
  import { Transition, TransitionChild } from "@rgossiaux/svelte-headlessui";
  import toast from "svelte-french-toast";
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  export let isOpen = false;
  export let onClose;
  export let loading = false;

  const toggle = () => {
    isOpen = !isOpen;
  };

  const closeModal = () => {
    toggle();
    onClose();
  };

  const createTunnelUser = async () => {
    if (email.length == 0) {
      inputRef.focus();
      return;
    }
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
        closeModal();
        dispatch("success", data.SecretKey);
      } else {
        toast.error(data.error);
      }
    } catch (err) {
      console.error(err);
    } finally {
      loading = false;
    }
  };

  let email = "";
  let inputRef: HTMLElement | null = null;
</script>

<div>
  <Transition
    show={isOpen}
    class="fixed inset-0 z-10 overflow-y-auto"
    on:afterEnter={() => {
      inputRef?.focus();
    }}
    on:afterLeave={(event) => {
      email = "";
    }}
  >
    <div
      class="flex items-end justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0"
    >
      <TransitionChild
        enter="ease-out duration-300"
        enterFrom="opacity-0"
        enterTo="opacity-100"
        leave="ease-in duration-200"
        leaveFrom="opacity-100"
        leaveTo="opacity-0"
      >
        <div class="fixed inset-0 transition-opacity">
          <div class="absolute inset-0 bg-gray-500 opacity-75" />
        </div>
      </TransitionChild>
      <!-- This element is to trick the browser into centering the modal contents. -->
      <span class="hidden sm:inline-block sm:align-middle sm:h-screen" />&#8203;
      <TransitionChild
        class="inline-block overflow-hidden text-left align-bottom transition-all transform bg-white rounded-lg shadow-xl sm:my-8 sm:align-middle sm:max-w-lg sm:w-full"
        role="dialog"
        aria-modal="true"
        aria-labelledby="modal-headline"
        enter="ease-out duration-300"
        enterFrom="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
        enterTo="opacity-100 translate-y-0 sm:scale-100"
        leave="ease-in duration-200"
        leaveFrom="opacity-100 translate-y-0 sm:scale-100"
        leaveTo="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
      >
        <div class="px-4 pt-5 pb-4 bg-white sm:p-6 sm:pb-4">
          <div class="sm:flex sm:items-start">
            <div
              class="flex items-center justify-center flex-shrink-0 w-12 h-12 mx-auto bg-green-100 rounded-full sm:mx-0 sm:h-10 sm:w-10"
            >
              <!-- Heroicon name: exclamation -->
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="w-6 h-6 text-green-900"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M12 9v6m3-3H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
            <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
              <h3
                class="text-lg font-medium leading-6 text-gray-900"
                id="modal-headline"
              >
                New tunnel user
              </h3>
              <div class="mt-2">
                <p class="text-sm leading-5 text-gray-500">
                  Once you create a user, you can create tunnel connections
                  using the generated secret key
                </p>
              </div>
              <div class="mt-2">
                <div>
                  <label
                    for="email"
                    class="block text-sm font-medium leading-5 text-gray-700"
                  >
                    Email address
                  </label>
                  <div class="relative mt-1 rounded-md shadow-sm">
                    <input
                      bind:this={inputRef}
                      bind:value={email}
                      id="email"
                      class="block rounded-lg w-full px-3 form-input sm:text-sm sm:leading-5"
                      placeholder="name@example.com"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="px-4 py-3 bg-gray-50 sm:px-6 sm:flex sm:flex-row-reverse">
          <span class="flex w-full rounded-md shadow-sm sm:ml-3 sm:w-auto">
            <button
              type="button"
              on:click={createTunnelUser}
              class="inline-flex justify-center w-full min-w-full px-4 py-2 text-base font-medium leading-6 text-white transition duration-150 ease-in-out bg-gray-600 border border-transparent rounded-md shadow-sm hover:bg-gray-500 focus:outline-none focus:border-gray-700 focus:shadow-outline-gray sm:text-sm sm:leading-5"
            >
              Create
            </button>
          </span>
          <span class="flex w-full mt-3 rounded-md shadow-sm sm:mt-0 sm:w-auto">
            <button
              on:click={closeModal}
              type="button"
              class="inline-flex justify-center w-full px-4 py-2 text-base font-medium leading-6 text-gray-700 transition duration-150 ease-in-out bg-white border border-gray-300 rounded-md shadow-sm hover:text-gray-500 focus:outline-none focus:border-blue-300 focus:shadow-outline-blue sm:text-sm sm:leading-5"
            >
              Cancel
            </button>
          </span>
        </div>
      </TransitionChild>
    </div>
  </Transition>
</div>
