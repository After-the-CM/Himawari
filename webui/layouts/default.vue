<template>
  <v-app>
    <v-navigation-drawer
      v-model="drawer"
      :clipped="clipped"
      fixed
      app
      temporary
    >
      <v-list-item-group>
        <v-list-item
          v-for="(item, i) in items"
          :key="i"
          :to="item.to"
          router
          exact
        >
          <v-list-item-content>
            <v-list-item-title v-text="item.title" />
          </v-list-item-content>
        </v-list-item>
        <a
          target="_blank"
          rel="noopener noreferrer"
          href="/Operation_manual.pdf"
          style="text-decoration: none; color: #000000"
        >
          <v-list-item>
            <v-list-item-content>
              <v-list-item-title>操作手順書</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </a>
      </v-list-item-group>

      <v-switch
        v-model="theme"
        class="ml-5"
        style="position: absolute; bottom: 0px"
      ></v-switch>
    </v-navigation-drawer>
    <v-app-bar :clipped-left="clipped" fixed app>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer" />
      <v-toolbar-title>
        <v-btn class="text-capitalize" text plain to="/"> Himawari </v-btn>
      </v-toolbar-title>
      <v-spacer />
    </v-app-bar>
    <v-main>
      <Nuxt />
    </v-main>
    <v-footer>
      <span
        >&copy; {{ new Date().getFullYear() }} After_the_CM
        <img src="favicon.ico" width="24px" height="24px"
      /></span>
    </v-footer>
  </v-app>
</template>

<script>
export default {
  data() {
    return {
      theme: this.$store.state.theme,
      clipped: false,
      drawer: false,
      fixed: true,
      items: [
        {
          title: 'Top Page',
          to: '/',
        },
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
