<template>
  <div class="tile">
    <div class="tile is-parent">
      <div class="card tile is-child">
        <div class="card-content">
          <div class="columns">
            <div class="column">
              <span class="title">{{ formatNumber(Math.floor(summary.throughput)) }}</span>
              <p class="content">responses / second</p>
            </div>
            <div class="column has-text-right is-one-fifth m-1">
              <span class="icon is-large"><i class="fas fa-space-shuttle fa-3x"></i></span>
            </div>
          </div>
        </div>
        <div class="card-footer">
          <div class="card-footer-item has-text-left">
            <div class="content">
              Total requests: <strong>{{ formatNumber(summary.requests) }}</strong>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="tile is-parent">
      <div class="card tile is-child">
        <div class="card-content">
          <div class="columns">
            <div class="column">
              <span class="title">{{ (summary.meanLatencies / 1000000).toFixed(1) }}ms</span>
              <p class="content">avg response time</p>
            </div>
            <div class="column has-text-right is-one-fifth m-1">
              <span class="icon is-large"><i class="fas fa-hourglass-end fa-3x"></i></span>
            </div>
          </div>
        </div>
        <div class="card-footer">
          <div class="card-footer-item has-text-left">
            P50:&nbsp;<strong>{{ (summary.meanP50Latencies / 1000000).toFixed(1) }}ms</strong>
          </div>
          <div class="card-footer-item has-text-left">
            P95:&nbsp;<strong>{{ (summary.meanP95Latencies / 1000000).toFixed(1) }}ms</strong>
          </div>
          <div class="card-footer-item has-text-left">
            P99:&nbsp;<strong>{{ (summary.meanP99Latencies / 1000000).toFixed(1) }}ms</strong>
          </div>
        </div>
      </div>
    </div>
    <div class="tile is-parent">
      <div class="card tile is-child">
        <div class="card-content">
          <div class="columns">
            <div class="column">
              <span class="title">{{ Math.floor(summary.success * 100) }}%</span>
              <p class="content">success rate</p>
            </div>
            <div class="column has-text-right is-one-fifth m-1">
              <span class="icon is-large"><i class="fas fa-chart-pie fa-3x"></i></span>
            </div>
          </div>
        </div>
        <div class="card-footer">
          <div class="card-footer-item has-text-left">
            2xx:&nbsp;<strong>{{ formatNumber(summary.statusCodes['2xx']) }}</strong>
          </div>
          <div class="card-footer-item has-text-left">
            4xx:&nbsp;<strong>{{ formatNumber(summary.statusCodes['4xx']) }}</strong>
          </div>
          <div class="card-footer-item has-text-left">
            5xx:&nbsp;<strong>{{ formatNumber(summary.statusCodes['5xx']) }}</strong>
          </div>
        </div>
      </div>
    </div>
    <div class="tile is-parent">
      <div class="card tile is-child">
        <div class="card-content">
          <div class="columns">
            <div class="column">
              <span class="title">{{ formatBytes(summary.totalBytesIn, 2) }}</span>
              <p class="content">bytes in</p>
            </div>
            <div class="column has-text-right is-one-fifth m-1">
              <span class="icon is-large"><i class="fas fa-exchange-alt fa-3x"></i></span>
            </div>
          </div>
        </div>
        <div class="card-footer">
          <div class="card-footer-item has-text-left">
            Bytes out:&nbsp;<strong>{{ formatBytes(summary.totalBytesOut, 2) }}</strong>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ResultSummary',
  props: {
    summary: Object,
  },
  methods: {
    formatBytes(bytes, decimals = 2) {
      if (bytes === 0) return '0 Byte';

      const k = 1024;
      const dm = decimals < 0 ? 0 : decimals;
      const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

      const i = Math.floor(Math.log(bytes) / Math.log(k));

      return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
    },
    formatNumber(value){
      return value.toLocaleString()
    }
  }
}
</script>

<style scoped>
.card .title {
  color: #3064D0;
}

.card-footer-item {
  justify-content: left;
  font-size: 0.9rem;
  color: #6a6a6a;
}

.card-content .column .icon.is-large {
  color: rgba(28, 80, 188, 0.35);
}

.card-footer-item:first-child {
  margin-left: 0.8rem;
}
</style>
