import client from '@/api/client';

export const done = {
  namespaced: true,
  state: {
    doneList: {},
  },
  getters: {
    doneList: state => {
      return state.doneList;
    },
  },
  mutations: {
    doneList: (state, data) => {
      state.doneList = data;
    },
  },
  actions: {
    async getDone({ commit }) {
      const res = await client.list();
      commit('doneList', res.list);
    },
  },
};
