export const state = () => ({
  crawlURL: '',
})

export const mutations = {
  changecrawlURL(state, newstate) {
    state.crawlURL = newstate
  },
}
