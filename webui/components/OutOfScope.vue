<template>
  <v-expansion-panels multiple class="mx-auto">
    <v-expansion-panel>
      <v-expansion-panel-header class="text-h6"
        >Out-Of-Scope</v-expansion-panel-header
      >

      <v-expansion-panel-content>
        <v-card flat>
          <v-expansion-panels
            flat
            v-for="(outOfScpeURLs, whereURLFound) in outOfScopes"
            :key="whereURLFound"
          >
            <v-expansion-panel>
              <v-expansion-panel-header>
                <v-row>
                  <v-col>{{ 'Find in ' + whereURLFound }}</v-col>
                </v-row>
              </v-expansion-panel-header>
              <v-expansion-panel-content
                v-for="(outOfScopeURL, j) in outOfScpeURLs"
                :key="j"
              >
                <v-row>
                  <v-col cols="1" style="color: gray">
                    {{ j + 1 }}
                  </v-col>
                  <v-col cols="11">
                    {{ outOfScopeURL }}
                  </v-col>
                </v-row>
              </v-expansion-panel-content>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-card>
      </v-expansion-panel-content>
    </v-expansion-panel>
  </v-expansion-panels>
</template>

<script>
export default {
  data() {
    return {
      outOfScopes: null,
    }
  },
  created() {
    this.getOutOfScpoe()
  },
  methods: {
    getOutOfScpoe() {
      this.$axios
        .$get('/api/outoforigin')
        .then((response) => {
          this.outOfScopes = response
        })
        .catch((err) => {
          // eslint-disable-next-line no-console
          console.log('err:', err)
        })
    },
  },
}
</script>

<style></style>
