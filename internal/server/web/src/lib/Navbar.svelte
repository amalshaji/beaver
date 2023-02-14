<script lang="ts">
  import beaverPng from "../assets/beaver.png";
  import { navigate } from "svelte-routing";

  const logoutUser = async () => {
    const res = await fetch("/api/v1/logout", {
      method: "POST",
      credentials: "same-origin",
    });
    if (res.status == 200) {
      navigate("/", { replace: true });
    }
  };
</script>

<nav class="bg-white shadow">
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
    <div class="flex justify-between h-16">
      <div class="flex">
        <div class="flex-shrink-0 flex items-center">
          <img
            class="block lg:hidden h-8 w-auto"
            src={beaverPng}
            alt="Workflow"
          />
          <img
            class="hidden lg:block h-8 w-auto"
            src={beaverPng}
            alt="Workflow"
          />
        </div>
      </div>
      <div class="hidden sm:ml-6 sm:flex sm:items-center">
        <button
          on:click={logoutUser}
          type="button"
          class="bg-white p-1 rounded-md border px-4 py-2  text-gray-800 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-900"
        >
          Logout
        </button>
      </div>
      <div class="-mr-2 flex items-center sm:hidden">
        <!-- Mobile menu button -->
        <button
          type="button"
          class="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500"
          aria-controls="mobile-menu"
          aria-expanded="false"
        >
          <span class="sr-only">Open main menu</span>
          <!--
              Icon when menu is closed.
  
              Heroicon name: outline/menu
  
              Menu open: "hidden", Menu closed: "block"
            -->
          <svg
            class="block h-6 w-6"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            aria-hidden="true"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 6h16M4 12h16M4 18h16"
            />
          </svg>
          <!--
              Icon when menu is open.
  
              Heroicon name: outline/x
  
              Menu open: "block", Menu closed: "hidden"
            -->
          <svg
            class="hidden h-6 w-6"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            aria-hidden="true"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>
      </div>
    </div>
  </div>

  <!-- Mobile menu, show/hide based on menu state. -->
</nav>
