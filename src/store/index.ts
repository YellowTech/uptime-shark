import { createStore } from 'vuex'
import * as iots from 'io-ts'
import { withFallback } from "io-ts-types/lib/withFallback";
import { isRight } from 'fp-ts/Either'

const logEntry = iots.type({
    Failed: iots.boolean,
    Message: iots.string,
    Time: iots.number,
});

const monitor = iots.type({
  id: iots.string,
  name: iots.string,
  interval: iots.number,
  status: withFallback(iots.boolean, false),
  statusMessage: iots.string,
  inverted: withFallback(iots.boolean, false),
  mode: iots.string,
  url: iots.string,
  logs: iots.array(logEntry),
})


const monitors = iots.array(monitor)

export default createStore({
  state: {
    counter: 0,
    apiDomain: "http://localhost:8000",
    error: false,
    loaded: false,
    monitors: monitors._A,
  },
  mutations: {
    increment(state) {
      state.counter++;
      // this.$store.commit('increase');
    }
  },
  actions: {
    increment(context) {
      context.commit('increment');
      // this.$store.dispatch('increase');
    },

    fetchData() {
      fetch(this.state.apiDomain + '/api/status')
      .then(response => response.text())
      .then(text => {
        const data = JSON.parse(text.replace("while(1);", ""));
        const result = monitors.decode(data)
        if(isRight(result)) {
          // eslint-disable-next-line
          this.state.monitors = result.right;
          this.state.loaded = true
          this.state.error = false
        } else {
          console.log(`Error when decoding!`)
          console.log(data)
          this.state.loaded = true
          this.state.error = true
        }
      });
    }
  },
  modules: {
  },
  getters: {
    absCounter(state) {
      return Math.abs(state.counter);
    }
  }
})
