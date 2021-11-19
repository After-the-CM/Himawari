export const state = () => ({
  // carawlURL: 'crawlurl',
  loginRef: '',
  loginURL: '',
})

export const mutations = {
  changeloginRef(state, newstate) {
    state.loginRef = newstate
  },
  changeloginURL(state, newstate) {
    state.loginURL = newstate
  },
}
