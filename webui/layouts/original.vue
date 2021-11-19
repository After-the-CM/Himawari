<template>
  <v-app>
    <v-navigation-drawer v-model="drawer" app fixed temporary>
      <!--  -->
      <v-list-item-group>
        <v-list-item v-for="(item, i) in items" :key="i" :to="item.link">
          <v-list-item-content>
            <v-list-item-title v-text="item.text"></v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list-item-group>
      <v-switch
        v-model="theme"
        class="ml-5"
        style="position: absolute; bottom: 0px"
      ></v-switch>
    </v-navigation-drawer>

    <v-app-bar app>
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>
        <v-btn class="text-capitalize" text plain to="/"> Himawari </v-btn>
      </v-toolbar-title>
    </v-app-bar>

    <v-main>
      <!--  -->
      <nuxt />
    </v-main>

    <span
      >&copy; {{ new Date().getFullYear() }} After_the_CM
      <img src="favicon.ico" width="24px" height="24px"
    /></span>
  </v-app>
</template>

<script>
export default {
  data() {
    return {
      theme: this.$store.state.theme,
      drawer: false,
      items: [
        { text: 'TopPage', link: '/' },
        { text: 'Crawler', link: '/crawler' },
        { text: 'Scanner', link: '/scanner' },
      ],
    }
  },
  watch: {
    theme() {
      this.$store.dispatch('theme', this.theme)
      this.$vuetify.theme.dark = this.theme
    },
  },
}
</script>
