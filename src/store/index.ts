import { createStore } from 'vuex'
import * as t from 'io-ts/lib/Decoder'
import { isRight } from 'fp-ts/Either'

const state = t.type({
  name: t.string,
  service: t.string,
  status: t.boolean
})

const stateArr = t.array(state)

export default createStore({
  state: {
    counter: 0,
    apiDomain: "http://localhost:8000",
    status: [],
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
        const data = JSON.parse(text.replace('while(1);', ''))
        const result = stateArr.decode(data)
        if(isRight(result)) {
          // eslint-disable-next-line
          this.state.status = result.right as any
        } else {
          console.log(`Error! ${data}`)
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
