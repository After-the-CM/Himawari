export const state = () => ({
  // carawlURL: 'crawlurl',
  crawlURL: '',
})

export const mutations = {
  changecrawlURL(state, newstate) {
    state.crawlURL = newstate
  },
}
