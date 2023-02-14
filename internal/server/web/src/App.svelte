<script>
  import { Router, Link, Route } from "svelte-routing";
  import Index from "./pages/Index.svelte";
  import Dashboard from "./pages/Dashboard.svelte";
  import { Toaster } from "svelte-french-toast";

  export let url = "";
</script>

<Router {url}>
  <div>
    <Route path="dashboard" component={Dashboard} />
    <Route path="/"><Index /></Route>
  </div>
</Router>

<Toaster position="top-center" />
