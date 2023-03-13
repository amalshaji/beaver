<script lang="ts">
  import { Transition, TransitionChild } from "@rgossiaux/svelte-headlessui";
  import toast from "svelte-french-toast";

  export let isOpen = true;
  export let onClose;
  let loading = false;
  export let userToDelete;

  const toggle = () => {
    isOpen = !isOpen;
  };

  const closeModal = () => {
    toggle();
    onClose();
  };

  const deleteTunnelUser = async () => {
    loading = true;
    try {
      const res = await fetch(`/api/v1/tunnel-users/${userToDelete.ID}`, {
        method: "DELETE",
      });
      const data = await res.json();
      if (res.status == 200) {
        closeModal();
      } else {
        toast.error(data.error);
      }
    } catch (err) {
      console.error(err);
    } finally {
      loading = false;
    }
  };
</script>

<div>
  <Transition
    show={isOpen}
    class="fixed inset-0 z-10 overflow-y-auto"
    on:afterLeave={(event) => {
      userToDelete = undefined;
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
              class="flex items-center justify-center flex-shrink-0 w-12 h-12 mx-auto bg-red-100 rounded-full sm:mx-0 sm:h-10 sm:w-10"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="w-6 h-6 text-red-900"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0"
                />
              </svg>
            </div>
            <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
              <h3
                class="text-lg font-medium leading-6 text-gray-900"
                id="modal-headline"
              >
                Delete tunnel user
              </h3>
              <div class="mt-2">
                <p class="text-sm leading-5 text-gray-500">
                  You are about to delete tunnel user <span
                    class="font-semibold"
                    >{userToDelete === undefined
                      ? ""
                      : userToDelete.Email}</span
                  >. This action is permanent and cannot be reversed.
                </p>
              </div>
            </div>
          </div>
        </div>
        <div class="px-4 py-3 bg-gray-50 sm:px-6 sm:flex sm:flex-row-reverse">
          <span class="flex w-full rounded-md shadow-sm sm:ml-3 sm:w-auto">
            <button
              type="button"
              on:click={deleteTunnelUser}
              class="inline-flex justify-center w-full min-w-full px-4 py-2 text-base font-medium leading-6 text-white transition duration-150 ease-in-out bg-red-600 border border-transparent rounded-md shadow-sm hover:bg-red-500 focus:outline-none focus:border-red-700 focus:shadow-outline-red sm:text-sm sm:leading-5"
            >
              Delete
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
