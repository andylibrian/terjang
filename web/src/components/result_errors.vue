<template>
  <div class="content">
    <div class="card">
      <div class="card-header">
        <strong class="card-header-title"><span class="icon"><i class="fas fa-exclamation-triangle" aria-hidden="true"></i></span> Errors</strong>
      </div>
      <div class="card-content">
        <div v-for="(worker, name) in workers" :key="name">
          <div v-for="(err, id) in worker.metrics.errors" :key="name + '-' + id">
            <strong>{{ name }}</strong>&nbsp; {{ err }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    name: 'ResultErrors',
    props: {
      workers: Object,
    },
    methods: {
    }
  }
</script>

<style scoped>
</style>
