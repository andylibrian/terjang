<template>
  <nav id="main-navbar" class="navbar is-fixed-top">
    <div class="navbar-brand">
      <h1 style="padding: 14px 14px 11px 14px" class="title is-5">Terjang</h1>
    </div>
    <div class="navbar-end">
      <div class="navbar-item">
          <strong>{{ serverInfo.state }}</strong>
      </div>
      <div class="navbar-item">
          <strong>{{ serverInfo.num_of_workers }}</strong>&nbsp;Worker Nodes connected
      </div>
    </div>
  </nav>
</template>

<script>
export default {
  name: 'Navbar',
  props: {
    serverInfo: Object,
  },
}
</script>

<style scoped>
  #main-navbar {
    border-bottom: 1px solid rgba(14, 16, 18, 0.06);
    background-color: #16304e;
    background-image: linear-gradient(to bottom, #192d3d, #162a3a);
    color: #dedede !important;
  }

  #main-navbar * {
    color: #dedede !important;
  }
</style>
