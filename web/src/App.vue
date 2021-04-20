<template>
  <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/bulma@0.9.1/css/bulma.min.css">
  <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.1/css/all.min.css">
  <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/@creativebulma/bulma-tooltip@1.2.0/dist/bulma-tooltip.min.css">
  <link rel="preconnect" href="https://fonts.gstatic.com">
  <link href="https://fonts.googleapis.com/css2?family=Lato&family=Open+Sans&family=Roboto&display=swap" rel="stylesheet"> 

  <Navbar :serverInfo="serverInfo" />

  <div class="section">
    <LaunchTest :serverInfo="serverInfo" :serverBaseUrl="serverBaseUrl" />
  </div>

  <div class="section" :class="{'is-hidden': !isResultVisible}">
    <ResultSummary :summary="metricsSummary"/>
  </div>

  <div class="section" :class="{'is-hidden': !isResultVisible}">
    <ResultDetail :workers="workers" :summary="metricsSummary" />
  </div>

  <div class="section" :class="{'is-hidden': !isErrorsVisible}">
    <ResultErrors :workers="workers" />
  </div>
</template>

<script>
  import CreateWebsocket from './lib/websocket.js'
  import Navbar from './components/navbar.vue'
  import LaunchTest from './components/launch_test.vue'
  import ResultSummary from './components/result_summary.vue'
  import ResultDetail from './components/result_detail.vue'
  import ResultErrors from './components/result_errors.vue'

  let serverBaseUrl = process.env.VUE_APP_SERVER_BASE_URL;
  if (!serverBaseUrl) {
    serverBaseUrl = location.protocol+'//'+location.hostname+(location.port ? ':'+location.port: '');
  }

  export default {
    name: 'App',
    components: {
      Navbar,
      LaunchTest,
      ResultSummary,
      ResultDetail,
      ResultErrors,
    },
    data: function() {
      return {
        serverInfo: {
          num_of_workers: 0,
          state: "",
        },
        workers: {},
        serverBaseUrl: serverBaseUrl,
        isResultVisible: false,
        isErrorsVisible: false,
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

            if (!("kind" in msg)) {
              return;
            }

            const obj = JSON.parse(msg.data);
            if (!obj) {
              return;
            }

            if (msg.kind === "ServerInfo") {
              Object.assign(_this.serverInfo, obj);

              if ('state' in obj) {
                if (!_this.isResultVisible && _this.serverInfo.state.toLowerCase() != 'notstarted') {
                  _this.isResultVisible = true;
                }
              }
            } else if (msg.kind === "WorkersInfo") {

              //reset workers object and start reading from server again
              _this.workers = {};
              
              for (let key in obj) {
                const worker = obj[key];
                if ("name" in worker && "metrics" in worker) {
                  _this.workers[worker.name] = worker;

                  if (!_this.isErrorsVisible && 'errors' in worker.metrics && worker.metrics.errors && worker.metrics.errors.length) {
                    _this.isErrorsVisible = true;
                  }
                }
              }
            }
        } catch(e) {
            console.error("Can not parse incoming notification message as JSON", evt.data, e);
        }
      }

      ws.onerror = function(evt) {
        console.log("error", evt);
      }
    },
    computed: {
      metricsSummary() {
        let requests = 0;
        let rate = 0;
        let throughput = 0;
        let success = 0;
        let workerCount = 0;
        let sumMeanLatencies = 0;
        let sumP50Latencies = 0;
        let sumP95Latencies = 0;
        let sumP99Latencies = 0;
        let meanLatencies = 0;
        let meanP50Latencies = 0;
        let meanP95Latencies = 0;
        let meanP99Latencies = 0;
        let totalBytesIn = 0;
        let totalBytesOut = 0;
        let statusCodes = {'2xx': 0, '4xx': 0, '5xx': 0};

        for (let w in this.workers) {
          requests += this.workers[w].metrics.requests;
          rate += this.workers[w].metrics.rate;
          throughput += this.workers[w].metrics.throughput;
          success += this.workers[w].metrics.success;
          sumMeanLatencies += this.workers[w].metrics.latencies.mean;
          sumP50Latencies += this.workers[w].metrics.latencies["50th"];
          sumP95Latencies += this.workers[w].metrics.latencies["95th"];
          sumP99Latencies += this.workers[w].metrics.latencies["99th"];
          totalBytesIn += this.workers[w].metrics.bytes_in.total;
          totalBytesOut += this.workers[w].metrics.bytes_out.total;

          for (let s in this.workers[w].metrics.status_codes) {
            if (s >= 200 && s <= 299) {
              statusCodes['2xx'] += this.workers[w].metrics.status_codes[s]
            } else if (s >= 400 && s <= 499) {
              statusCodes['4xx'] += this.workers[w].metrics.status_codes[s]
            } else if (s >= 500 && s <= 599) {
              statusCodes['5xx'] += this.workers[w].metrics.status_codes[s]
            }
          }

          workerCount++;
        }

        if (workerCount > 0) {
          success = success / workerCount;
          meanLatencies = sumMeanLatencies / workerCount;
          meanP50Latencies = sumP50Latencies / workerCount;
          meanP95Latencies = sumP95Latencies / workerCount;
          meanP99Latencies = sumP99Latencies / workerCount;
        }

        return {
          requests, rate, throughput, success, meanLatencies, meanP50Latencies, meanP95Latencies, meanP99Latencies, totalBytesIn, totalBytesOut, statusCodes
        };
      },
    }
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
    border: 1px solid rgba(14, 18, 23, 0.09);
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

  .section > .content {
      padding: 0 1rem !important;
  }
</style>

