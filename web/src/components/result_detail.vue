<template>
  <div class="content">
    <div class="card">
      <div class="card-header">
        <strong class="card-header-title"><span class="icon"><i class="fas fa-table" aria-hidden="true"></i></span> Result</strong>
      </div>
      <div class="card-content">
        <table class="table is-bordered is-hoverable is-fullwidth">
          <thead>
            <tr>
              <th>Worker</th>
              <th>Requests (sum)</th>
              <th>Rate (per second)</th>
              <th>Throughput (per second)</th>
              <th>Success (%)</th>
              <th>Avg Resp Time (ms)</th>
              <th>P95 Resp Time (ms)</th>
              <th><span class="has-tooltip-arrow" data-tooltip="The total number of bytes received with the response bodies.">Total Bytes In</span></th>
              <th><span class="has-tooltip-arrow" data-tooltip="The total number of bytes sent with the request bodies.">Total Bytes Out</span></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(worker, name) in workers" :key="name">
              <td>{{ name }}</td>
              <td>{{ worker.metrics.requests }}</td>
              <td>{{ Math.floor(worker.metrics.rate) }}</td>
              <td>{{ Math.floor(worker.metrics.throughput) }}</td>
              <td>{{ Math.floor(worker.metrics.success * 100) }}</td>
              <td>{{ worker.metrics.latencies.mean / 1000000 }}</td>
              <td>{{ worker.metrics.latencies['95th'] / 1000000 }}</td>
              <td>{{ worker.metrics.bytes_in.total }}</td>
              <td>{{ worker.metrics.bytes_out.total }}</td>
            </tr>
          </tbody>
          <tfoot>
            <tr>
              <th>Summary</th>
              <th>{{ summary.requests }}</th>
              <th>{{ Math.floor(summary.rate) }}</th>
              <th>{{ Math.floor(summary.throughput) }}</th>
              <th>{{ Math.floor(summary.success * 100) }}</th>
              <th>{{ summary.meanLatencies / 1000000 }}</th>
              <th>{{ summary.meanP95Latencies / 1000000 }}</th>
              <th>{{ summary.totalBytesIn }}</th>
              <th>{{ summary.totalBytesOut }}</th>
            </tr>
          </tfoot>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    name: 'ResultDetail',
    props: {
      workers: Object,
      summary: Object,
    },
    methods: {
    }
  }
</script>

<style scoped>
  .table {
    color: #464646;
  }
</style>
