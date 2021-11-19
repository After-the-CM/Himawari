export const state = () => ({
  // carawlURL: 'crawlurl',
  crawlParams: [
    { name: '*', value: 'Himawari' },
    { name: 'key', value: 'hello' },
    { name: 'email', value: 'Himawari@example.com' },
    { name: 'url', value: 'http://example.com' },
    { name: 'tel', value: '00012345678' },
    { name: 'date', value: '2020-12-16' },
    { name: 'text', value: 'Himawari' },
    { name: 'textarea', value: 'Himawari' },
    { name: 'input', value: 'I am Himawari' },
  ],
})

export const mutations = {
  changecrawlParams(state, newstate) {
    state.crawlParams = newstate
  },
}
