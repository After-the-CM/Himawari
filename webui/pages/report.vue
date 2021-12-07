<template>
  <div>
    <scaning-progress-bar :flag="scanningFlag" />
    <div v-if="alertFlag">
      <v-alert
        color="pink darken-1"
        border="bottom"
        dark
        dismissible
        style="position: fixed; width: 100%; z-index: 1000"
        ><div class="font-weight-bold">Scan completed</div></v-alert
      >
    </div>
    <v-app class="mx-auto" style="width: 90%">
      <v-container>
        <div class="text-center d-flex pb-4"></div>
        <v-expansion-panels
          multiple
          class="mx-auto"
          v-for="(item, i) in vuln"
          :key="i"
        >
          <v-expansion-panel v-if="item.Issues.length > 0">
            <v-expansion-panel-header>
              <v-row>
                <v-col cols="1">{{ item.CWE }}</v-col>
                <v-col> {{ item.Name }}</v-col>
                <v-spacer></v-spacer>
                <v-col xs="4" sm="4" md="4" right> {{ item.Severity }}</v-col>
                <v-col
                  xs="1"
                  sm="1"
                  md="1"
                  right
                  class="font-weight-black"
                  style="color: RED"
                >
                  {{ item.Issues.length }}
                </v-col>
              </v-row>
            </v-expansion-panel-header>
            <v-expansion-panel-content>
              <v-card flat>
                <v-card-title>脆弱性の説明</v-card-title>
                <v-card-text
                  style="white-space: pre-wrap; word-wrap: break-word"
                  >{{ item.Description }}</v-card-text
                >
              </v-card>
              <v-card flat>
                <v-card-title>必須対策</v-card-title>
                <v-card-text
                  style="white-space: pre-wrap; word-wrap: break-word"
                  >{{ item.Mandatory }}</v-card-text
                >
              </v-card>
              <v-card flat>
                <v-card-title
                  style="white-space: pre-wrap; word-wrap: break-word"
                  >保険的対策</v-card-title
                >
                <v-card-text>{{ item.Insurance }}</v-card-text>
              </v-card>

              <v-expansion-panels
                flat
                multiple
                class="mx-auto"
                v-for="(issue, j) in item.Issues"
                :key="j"
              >
                <v-expansion-panel>
                  <v-expansion-panel-header>
                    <v-row>
                      <v-col cols="1" style="color: gray">{{ j + 1 }}</v-col>
                      <v-col> {{ issue.URL }} </v-col>
                      <v-col xs="5" sm="5" md="5" right>
                        Parameter: {{ issue.Parameter }}
                      </v-col>
                    </v-row>
                  </v-expansion-panel-header>
                  <v-row>
                    <v-col cols="3" sm="12" md="12" lg="3">
                      <v-row>
                        <v-expansion-panel-content>
                          <v-card
                            outlined
                            class="overflow-y-auto overflow-x-auto"
                          >
                            <v-card-title>Parameter</v-card-title>
                            <v-card-text>
                              {{ issue.Parameter }}
                            </v-card-text>
                          </v-card>
                          <v-card
                            outlined
                            class="overflow-y-auto overflow-x-auto"
                          >
                            <v-card-title>Payload</v-card-title>
                            <v-card-text>
                              {{ issue.Payload }}
                            </v-card-text>
                          </v-card>
                          <v-card
                            outlined
                            class="overflow-y-auto overflow-x-auto"
                          >
                            <v-card-title>Evidence</v-card-title>
                            <v-card-text>
                              {{ issue.Evidence }}
                            </v-card-text>
                          </v-card>
                        </v-expansion-panel-content>
                      </v-row>
                    </v-col>
                    <v-col cols="4">
                      <v-expansion-panel-content>
                        <v-card
                          outlined
                          sm="12"
                          md="12"
                          lg="4"
                          class="overflow-y-auto overflow-x-auto"
                          max-height="500"
                          style="width: 500px"
                        >
                          <v-card-title>Request</v-card-title>
                          <v-card-text
                            style="white-space: pre-wrap; word-wrap: break-word"
                            >{{ issue.Request }}</v-card-text
                          >
                        </v-card>
                      </v-expansion-panel-content>
                    </v-col>
                    <v-col cols="4">
                      <v-expansion-panel-content>
                        <v-card
                          outlined
                          sm="12"
                          md="12"
                          lg="4"
                          color=""
                          class="overflow-y-auto overflow-x-auto"
                          max-height="500"
                          style="width: 500px"
                        >
                          <v-card-title>Response</v-card-title>
                          <v-card-text
                            style="white-space: pre-wrap; word-wrap: break-word"
                            >{{ issue.Response }}</v-card-text
                          >
                        </v-card>
                      </v-expansion-panel-content>
                    </v-col>
                  </v-row>
                </v-expansion-panel>
              </v-expansion-panels>
            </v-expansion-panel-content>
          </v-expansion-panel>
        </v-expansion-panels>
      </v-container>
      <v-btn
        class="text-capitalize"
        rounded
        v-show="!scanningFlag"
        href="/report/markdown"
        >Download Report(markdown)</v-btn
      >
    </v-app>
  </div>
</template>

<script>
export default {
  middleware({ $cookies, redirect }) {
    if ($cookies.get('agree') !== 'Agree') {
      return redirect('/')
    }
  },
  data() {
    return {
      panel: [],
      items: 5,
      vuln: [],
      finishGetVulnsKey: null,
      scanflag: 'scanning',
      scanningFlag: true,
      alertFlag: false,
    }
  },
  methods: {
    vulnSort() {
      // error回避のためslice()をいれる
      const sortmap = {
        High: 3,
        Medium: 2,
        Low: 1,
      }
      this.vuln.sort((a, b) => {
        if (sortmap[a.Severity] > sortmap[b.Severity]) {
          return -1
        }
        if (sortmap[a.Severity] < sortmap[b.Severity]) {
          return 1
        }
        return 0
      })
    },
    getVulns() {
      this.$axios
        .$get('/api/vuln')
        .then((response) => {
          this.vuln = []
          for (const i in response) {
            this.vuln.push(response[i])
          }
          this.vulnSort()
        })
        .catch((err) => {
          console.log('err:', err)
        })
    },
    getscanflag() {
      this.$axios
        .$get('/api/scanflag')
        .then((response) => {
          this.scanflag = response
          if (this.scanflag !== 'scanning') {
            this.scanningFlag = false
            this.setalertflag()
          }
        })
        .catch((err) => {
          // eslint-disable-next-liue no-console
          console.log('err:', err)
        })
    },
    setalertflag() {
      this.alertFlag = true
      setTimeout(() => {
        this.alertFlag = false
      }, 5000)
    },
  },
  mounted() {
    this.getVulns()
    this.getscanflag()

    this.finishGetVulnsKey = setInterval(() => {
      this.getVulns()
      this.getscanflag()
      if (this.scanflag !== 'scanning') {
        clearInterval(this.finishGetVulnsKey)
      }
    }, 5000)

    window.addEventListener('beforeunload', () => {
      clearInterval(this.finishGetVulnsKey)
    })
  },
  beforeDestroy() {
    window.removeEventListener('beforeunload', () => {
      clearInterval(this.finishGetVulnsKey)
    })

    clearInterval(this.finishGetVulnsKey)
  },
}
</script>
