<template>
  <ul>
    <sitemap-node :item="tree" />
  </ul>
</template>

<script>
export default {
  data() {
    return {
      tree: {},
    }
  },
  created() {
    this.doFetchSitemap()
  },
  methods: {
    doFetchSitemap() {
      this.$axios
        .$get('/api/sitemap')
        .then((response) => {
          const resultSitemap = response
          this.tree = resultSitemap
        })
        .catch((err) => {
          // eslint-disable-next-line no-console
          console.log('err:', err)
        })
    },
  },
}
</script>

<style>
.item {
  cursor: pointer;
}
.bold {
  font-weight: bold;
}
ul {
  padding-left: 1em;
  line-height: 1.5em;
  list-style-type: dot;
}
</style>
