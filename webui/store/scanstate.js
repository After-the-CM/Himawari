export const state = () => ({
  scanstate: null,
})

export const mutations = {
  changestate(state, { newstate }) {
    state.scanstate = newstate
  },
}
