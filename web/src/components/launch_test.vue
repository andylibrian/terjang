<template>
  <div class="tiles">
    <div class="tile is-parent">
      <div class="tile is-child is-6">
        <div class="card">
          <div class="card-header">
            <strong class="card-header-title"><span class="icon"><i class="fas fa-rocket" aria-hidden="true"></i></span> Launch a load test</strong>
          </div>
          <div class="card-content">
            <form @submit.prevent>
              <div class="tabs">
                <ul>
                  <li :class="{'is-active': launchFormActiveTab == 'basic'}"><a @click="launchFormSwitchTab('basic')">Basic Settings</a></li>
                  <li :class="{'is-active': launchFormActiveTab == 'headers'}"><a @click="launchFormSwitchTab('headers')">Headers</a></li>
                  <li :class="{'is-active': launchFormActiveTab == 'body'}"><a @click="launchFormSwitchTab('body')">Body</a></li>
                </ul>
              </div>
              <div id="launchForm-tab-content-basic" :class="{'is-hidden': launchFormActiveTab != 'basic'}">
                <div class="field is-horizontal">
                  <div class="field-label is-normal">
                    <label class="label" for="load-test-url">URL</label>
                  </div>
                  <div class="field-body">
                    <div class="field has-addons">
                      <p class="control">
                        <span class="select">
                          <select id="load-test-method" name="load-test-method">
                            <option value="GET">GET</option>
                            <option value="POST">POST</option>
                            <option value="PUT">PUT</option>
                            <option value="DELETE">DELETE</option>
                          </select>
                        </span>
                      </p>
                      <p class="control is-expanded has-icons-left">
                        <input class="input" id="load-test-url" type="text" placeholder="URL">
                        <span class="icon is-small is-left">
                          <span class="fas fa-globe-asia"></span>
                        </span>
                      </p>
                    </div>
                  </div>
                </div>
                <div class="field is-horizontal">
                  <div class="field-label is-normal">
                    <label class="label" for="load-test-duration">Duration</label>
                  </div>
                  <div class="field-body">
                    <div class="field has-addons">
                      <p class="control has-icons-left">
                        <input class="input" name="load-test-duration" id="load-test-duration" type="text" placeholder="Duration" value="30">
                        <span class="icon is-small is-left">
                          <span class="fas fa-stopwatch"></span>
                        </span>
                      </p>
                      <p class="control">
                        <span class="select">
                          <select id="load-test-duration-unit" name="load-test-duration-unit">
                            <option value="second">seconds</option>
                            <option value="minute">minutes</option>
                            <option value="hour">hours</option>
                          </select>
                        </span>
                      </p>
                    </div>
                  </div>
                </div>
                <div class="field is-horizontal">
                  <div class="field-label is-normal">
                    <label class="label" for="load-test-rate">Rate per worker</label>
                  </div>
                  <div class="field-body">
                    <div class="field has-addons">
                      <p class="control has-icons-left">
                        <input class="input" name="load-test-rate" id="load-test-rate" type="text" placeholder="Rate" value="100">
                        <span class="icon is-small is-left">
                          <span class="fas fa-bolt"></span>
                        </span>
                      </p>
                      <p class="control">
                        <a class="button is-static">
                          reqs per second
                        </a>
                      </p>
                    </div>
                  </div>
                </div>
              </div>
              <div id="launchForm-tab-content-headers" :class="{'is-hidden': launchFormActiveTab != 'headers'}">
                <div class="field is-horizontal">
                  <div class="field-label is-normal">
                    <label class="label" for="load-test-header">Header</label>
                  </div>
                  <div class="field-body">
                    <div class="field has-addons">
                      <div class="control is-expanded">
                        <textarea class="textarea" id="load-test-header" placeholder="Key: Value
Key: Value"></textarea>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div id="launchForm-tab-content-body" :class="{'is-hidden': launchFormActiveTab != 'body'}">
                <div class="field is-horizontal">
                  <div class="field-label is-normal">
                    <label class="label" for="load-test-body">Body</label>
                  </div>
                  <div class="field-body">
                    <textarea class="textarea" id="load-test-body"></textarea>
                  </div>
                </div>
              </div>
              <div class="field is-horizontal mt-4">
                <div class="field-label"></div>
                <div class="field-body">
                  <div class="field is-grouped">
                    <p class="control buttons">
                      <button v-bind:disabled="serverInfo.state.toLowerCase() == 'running'" class="button is-primary" @click="runLoadTest()">
                        Run
                      </button>
                      <button v-bind:disabled="serverInfo.state.toLowerCase() != 'running'" class="button is-danger" @click="stopLoadTest()">
                        Stop
                      </button>
                    </p>
                  </div>
                </div>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'LaunchTest',
  props: {
    serverInfo: Object,
    serverBaseUrl: String,
  },
  data: function() {
    return {
      launchFormActiveTab: "basic",
    }
  },
  methods: {
    launchFormSwitchTab(tabId) {
      this.launchFormActiveTab = tabId;
    },
    runLoadTest() {
      // TODO: validate form

      const methodEl = document.getElementById("load-test-method");
      if (!methodEl) {
        return false;
      }

      const urlEl = document.getElementById("load-test-url");
      if (!urlEl) {
        return false;
      }

      const durationUnitEl = document.getElementById("load-test-duration-unit");
      if (!durationUnitEl) {
        return false;
      }

      const durationEl = document.getElementById("load-test-duration");
      if (!durationEl) {
        return false;
      }

      const rateEl = document.getElementById("load-test-rate");
      if (!rateEl) {
        return false;
      }

      const headerEl = document.getElementById("load-test-header");
      if (!headerEl) {
        return false;
      }

      const bodyEl = document.getElementById("load-test-body");
      if (!bodyEl) {
        return false;
      }

      const method = methodEl.value;
      const url = urlEl.value;
      const durationUnit = durationUnitEl.value;
      let duration = durationEl.value;

      if (durationUnit == "minute") {
        duration *= 60;
      } else if (durationUnit == "hour") {
        duration *= 60 * 60;
      }
      const rate = rateEl.value;
      const body = bodyEl.value;
      const header = headerEl.value;

      var xhr = new XMLHttpRequest();
      xhr.open("POST", this.serverBaseUrl + '/api/v1/load_test', true)
      xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8");

      const postData = JSON.stringify({
        method: method,
        url: url,
        duration: duration,
        rate: rate,
        header: header,
        body: body,
      });

      xhr.send(postData);

      // TODO: handle response
    },
    stopLoadTest() {
      var xhr = new XMLHttpRequest();
      xhr.open("DELETE", this.serverBaseUrl + '/api/v1/load_test', true)
      xhr.send();

      // TODO: handle response
    },
  },
}
</script>

<style scoped>
</style>
