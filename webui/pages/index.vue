<template v-slot:activator="{ on, attrs }">
  <div id="app">
    <v-row class="pa-10" justify="center">
      <img src="Himawari.png" />
    </v-row>
    <v-row justify="center">
      <v-card elevation="24" width="700px" height="300px">
        <v-card-title class="text-h3 justify-center">
          使用上の注意
        </v-card-title>
        <v-card-text class="text-h5 justify-center">
          <div class="text-center">
            <div>許可なく他者のサーバーに使用しないでください。</div>
            <div>違法行為となり、罰せられる可能性があります。</div>
            <div>製作者一同は一切の責任を取りません。</div>
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-row class="mx-auto, my-5">
            <v-btn
              rounded
              @click="
                asyncData()
                crawler()
              "
            >
              同意して開始する
            </v-btn>
          </v-row>
        </v-card-actions>
      </v-card>
    </v-row>
    <v-row justify="center"> </v-row>
  </div>
</template>

<script>
export default {
  asyncData({ app }) {
    app.$cookies.remove('agree')
  },
  data() {
    return {
      theme: this.$store.state.theme,
      dialog: true,
      isflag: true,
    }
  },
  watch: {
    theme() {
      this.$store.dispatch('theme', this.theme)
      this.$vuetify.theme.dark = this.theme
    },
  },
  mounted() {
    this.$axios
      .$get('/api/reset')
      .then((response) => {
        this.vuln = response
      })
      .catch((err) => {
        console.log('err:', err)
      })
  },
  methods: {
    crawler() {
      this.$router.push('crawl')
    },
    asyncData() {
      this.$cookies.set('agree', 'Agree')
    },
    disableOn() {
      this.isflag = true
    },
    disableOff() {
      this.isflag = false
    },
  },
}
</script>
