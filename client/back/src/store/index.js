import Vue from 'vue';
import Vuex from 'vuex';
import { requestResults } from './modules/requestResults.js';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {},
  mutations: {},
  actions: {},
  modules: {
    requestResults,
  },
});
