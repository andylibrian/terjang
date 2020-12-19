<template>
  <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/bulma@0.9.1/css/bulma.min.css">
  <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.1/css/all.min.css">
  <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/@creativebulma/bulma-tooltip@1.2.0/dist/bulma-tooltip.min.css">
  <link rel="preconnect" href="https://fonts.gstatic.com">
  <link href="https://fonts.googleapis.com/css2?family=Lato&family=Open+Sans&family=Roboto&display=swap" rel="stylesheet"> 
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
  import CreateWebsocket from './lib/websocket.js'

  let serverBaseUrl = process.env.VUE_APP_SERVER_BASE_URL;
  if (!serverBaseUrl) {
    serverBaseUrl = location.protocol+'//'+location.hostname+(location.port ? ':'+location.port: '');
  }

  export default {
    name: 'App',
    data: function() {
      return {
        serverInfo: {
          num_of_workers: 0,
          state: "",
        },
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
        console.log("error", evt.data);
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

