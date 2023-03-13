<script lang="ts">
  import { Transition, TransitionChild } from "@rgossiaux/svelte-headlessui";
  import CopyToClipboard from "../CopyToClipboard.svelte";

  export let isOpen = false;
  export let onClose,
    secretKey = "";

  const toggle = () => {
    isOpen = !isOpen;
  };

  const closeModal = () => {
    toggle();
    onClose();
  };
</script>

<div>
  <Transition show={isOpen} class="fixed inset-0 z-10 overflow-y-auto">
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
              class="flex items-center justify-center flex-shrink-0 w-12 h-12 mx-auto bg-yellow-100 rounded-full sm:mx-0 sm:h-10 sm:w-10"
            >
              <!-- Heroicon name: exclamation -->
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
                  d="M12 9v3.75m0-10.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.75c0 5.592 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.57-.598-3.75h-.152c-3.196 0-6.1-1.249-8.25-3.286zm0 13.036h.008v.008H12v-.008z"
                />
              </svg>
            </div>
            <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
              <h3
                class="text-lg font-medium leading-6 text-gray-900"
                id="modal-headline"
              >
                Secret key
              </h3>
              <div class="mt-2">
                <p class="text-sm leading-5 text-gray-500">
                  The secret key will only be displayed once. Please copy and
                  save it locally.
                </p>
              </div>
              <div class="mt-2">
                <div>
                  <label
                    for="email"
                    class="block text-sm font-medium leading-5 text-gray-700"
                  >
                    Secret Key
                  </label>
                  <div
                    class="relative mt-1 rounded-md shadow-sm flex items-center"
                  >
                    <input
                      readonly
                      class="block rounded-lg w-full px-3 form-input border-none bg-gray-50 sm:text-sm sm:leading-5"
                      value={secretKey !== undefined ? secretKey : ""}
                    />
                    <CopyToClipboard text={secretKey} />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="px-4 py-3 bg-gray-50 sm:px-6 sm:flex sm:flex-row-reverse">
          <span class="flex w-full mt-3 rounded-md shadow-sm sm:mt-0 sm:w-auto">
            <button
              on:click={closeModal}
              type="button"
              class="inline-flex justify-center w-full px-4 py-2 text-base font-medium leading-6 text-gray-700 transition duration-150 ease-in-out bg-white border border-gray-300 rounded-md shadow-sm hover:text-gray-500 focus:outline-none focus:border-blue-300 focus:shadow-outline-blue sm:text-sm sm:leading-5"
            >
              I have saved the secret key
            </button>
          </span>
        </div>
      </TransitionChild>
    </div>
  </Transition>
</div>
