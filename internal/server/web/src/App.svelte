<script>
  import { Router, Link, Route } from "svelte-routing";
  import Index from "./pages/Index.svelte";
  import Dashboard from "./pages/Dashboard.svelte";
  import { Toaster } from "svelte-french-toast";
  import CreateSuperUser from "./pages/CreateSuperUser.svelte";

  export let url = "";
</script>

<Router {url}>
  <div>
    <Route path="dashboard" component={Dashboard} />
    <Route path="createsuperuser" component={CreateSuperUser} />
    <Route path="/"><Index /></Route>
  </div>
</Router>

<Toaster position="top-center" />
