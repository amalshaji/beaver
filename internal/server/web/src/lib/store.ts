import { writable } from "svelte/store";

export const tunnelUserConnectionStatus = writable<IActiveConnectionStatus[]>(
  []
);
