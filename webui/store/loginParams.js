export const state = () => ({
  // carawlURL: 'crawlurl',
  loginParams: [
    { key: '', value: '', method: 'POST' },
    { key: '', value: '', method: 'POST' },
  ],
})

export const mutations = {
  changeloginParams(state, newstate) {
    state.loginParams = newstate
  },
}
