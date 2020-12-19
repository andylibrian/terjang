<template>
  <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/bulma@0.9.1/css/bulma.min.css">
  <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.1/css/all.min.css">
  <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/@creativebulma/bulma-tooltip@1.2.0/dist/bulma-tooltip.min.css">
  <link rel="preconnect" href="https://fonts.gstatic.com">
  <link href="https://fonts.googleapis.com/css2?family=Lato&family=Open+Sans&family=Roboto&display=swap" rel="stylesheet"> 

  <Navbar :serverInfo="serverInfo" />
  <div id="section-launch-test" class="section">
    <LaunchTest :serverInfo="serverInfo" :serverBaseUrl="serverBaseUrl" />
  </div>

</template>

<script>
  import CreateWebsocket from './lib/websocket.js'
  import Navbar from './components/navbar.vue'
  import LaunchTest from './components/launch_test.vue'

  let serverBaseUrl = process.env.VUE_APP_SERVER_BASE_URL;
  if (!serverBaseUrl) {
    serverBaseUrl = location.protocol+'//'+location.hostname+(location.port ? ':'+location.port: '');
  }

  export default {
    name: 'App',
    components: {
      Navbar,
      LaunchTest,
    },
    data: function() {
      return {
        serverInfo: {
          num_of_workers: 0,
          state: "",
        },
        serverBaseUrl: serverBaseUrl,
      }
    },
    created: function() {
      const splitted = serverBaseUrl.split("/");
      const wsBaseUrl = 'ws://' + splitted[2];
      const _this = this;
      const ws = CreateWebsocket(wsBaseUrl);

      ws.onopen = function() {
        console.log("connected to the server");
      }

      ws.onclose = function() {
        console.log("connection to the server closed");
      }

      ws.onmessage = function(evt) {
        console.log("received from the server", evt.data);

        try {
            const msg = JSON.parse(evt.data);
            console.log(msg);

            if (!('kind' in msg)) {
              return;
            }

            const obj = JSON.parse(msg.data);
            if (!obj) {
              return;
            }

            if (msg.kind === "ServerInfo") {
              Object.assign(_this.serverInfo, obj);
            }
        } catch(e) {
            console.error("Can not parse incoming notification message as JSON", evt.data, e);
        }
      }

      ws.onerror = function(evt) {
        console.log("error", evt);
      }
    },
  }
</script>

<style>
  body {
    padding-top: 4rem !important;
    padding-bottom: 4rem !important;
    background-color: #ecedef;
    background-image: linear-gradient(315deg, #dcdddf 0%, #ecedef 100%);

    font-family: Lato, Roboto, 'Open sans', sans-serif;
  }

  .section {
    padding: 1.5rem 0 1.5rem 0 !important;
  }

  .section .card {
    box-shadow: none;
    border: 1px solid rgba(14, 18, 23, 0.08);
    background: rgba(255, 255, 255, 0.9);
  }

  .section .card .card-header {
    background: rgba(100, 110, 128, 0.2) !important;
  }
  .section .card .card-header .card-header-title {
    border-bottom: 1px solid rgba(100, 110, 128, 0.1) !important;
  }

  .section .card .card-footer strong {
    color: #525252;
  }

  .button.is-primary {
    background-color: #0477c1 !important;
  }

  .button.is-primary.is-hovered, .button.is-primary:hover {
    background-color: #0467b1 !important;
  }

  .button.is-danger {
    background-color: #c12424 !important;
  }

  .button.is-danger.is-hovered, .button.is-danger:hover {
    background-color: #d11414 !important;
  }

  .button.is-danger[disabled], fieldset[disabled] .button.is-danger {
    background-color: #e18888 !important;
  }
</style>

