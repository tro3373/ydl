import client from '@/api/client';

export const requestResults = {
  namespaced: true,
  state: {
    requestResults: [],
  },
  getters: {
    requestResults: (state) => {
      return state.requestResults;
    },
  },
  mutations: {
    requestResults: (state, data) => {
      state.requestResults = data;
    },
  },
  actions: {
    async getRequestResults({ commit }) {
      const res = await client.list();
      commit('requestResults', res.list || []);
    },
  },
};
