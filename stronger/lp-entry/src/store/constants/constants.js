const state = {
  constants: {
    total: 0,
    count: 0,
    tally: 0
  }
}

const getters = {
  constants: () => state.constants
}

const mutations = {
  UPDATE_CONSTANTS(state, info) {
    state.constants = Object.assign(state.constants, info)
  }
}

export default {
  state,
  getters,
  mutations
}