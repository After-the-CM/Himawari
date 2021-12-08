export const state = () => ({
  delay: '0',
})

export const mutations = {
  changeDelay(state, newstate) {
    state.delay = newstate
  },
}
