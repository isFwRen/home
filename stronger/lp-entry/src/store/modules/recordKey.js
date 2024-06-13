const state = {
  recordKey: [],
  page: 1
};

const mutations = {
  UPDATE_KEY(state, str) {
    state.recordKey.push(str)
  },
  UPDATE_PAGE(state, number) {
    state.page = number + 1
  },
  async RESET_LOG(state) {
    state.recordKey = []
    state.page = 1
  }
};

export default {
  state,
  mutations,
};