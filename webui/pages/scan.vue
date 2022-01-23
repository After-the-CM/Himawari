<template>
  <v-app>
    <v-container>
      <v-card class="mx-auto mt-5">
        <v-card-title style="background: #272727; color: white" class="mb-3"
          >Scanner</v-card-title
        >
        <v-row>
          <v-col cols="5">
            <v-card flat>
              <v-card-title> Sitemap</v-card-title>
              <ul>
                <sitemap-node :item="sitemap" />
              </ul>
            </v-card>
          </v-col>
          <v-spacer></v-spacer>
          <v-col cols="7" style="border-left: solid 1px gray">
            <v-card flat>
              <v-card-title>ScanOption</v-card-title>
              <v-card-text
                ><div class="font-weight-bold">
                  予想診断時間は {{ approximateTime }}分です。
                </div>
              </v-card-text>
              <radio-btn
                v-model="scanOption"
                :btnSetting="btnSetting"
                class="ml-5"
              />
              <v-col cols="2">
                <input-number
                  v-model="delay"
                  labelText="delay(ms)"
                  :inputRule="delayRule"
                  textId="delay"
                />
              </v-col>
              <login-option-switch v-model="loginflag" />
              <div v-if="loginflag">
                <input-text
                  v-model.trim="loginReferer"
                  labelText="LoginフォームがあるURL(Referer)"
                  textId="loginref"
                  textClass="mx-5"
                />
                <input-text
                  v-model.trim="loginURL"
                  labelText="Loginリクエストの送信先"
                  textId="loginurl"
                  textClass="mx-5"
                />
                <login-option
                  v-model="loginOptions"
                  :methods="loginmethods"
                  :methodcols="4"
                />
                <add-form-btn
                  :addform="loginOptions"
                  btnText="ログインパラメータ追加"
                  :adddata="{ key: '', value: '', method: 'POST' }"
                />
              </div>
              <div v-if="landmarkFlag">
                <v-text-field
                  v-model="landmarkNumber"
                  label="LandMark(default : 0)"
                  type="number"
                />
              </div>
            </v-card>
          </v-col>
        </v-row>
        <v-row>
          <v-btn
            rounded
            x-large
            justify="start"
            class="my-auto ma-auto mb-2 text-capitalize"
            href="/sitemap/download"
            >Download Sitemap</v-btn
          >
          <v-btn
            rounded
            x-large
            justify="end"
            class="my-auto ma-auto mb-2 text-capitalize"
            @click="doScan(), transitionsitemap()"
            :disabled="delay === '' || delay < 0"
            >Start Scan</v-btn
          >
        </v-row>
      </v-card>
    </v-container>
    <v-container>
      <out-of-scope />
    </v-container>
  </v-app>
</template>

<script>
import { cloneDeep } from 'lodash'

export default {
  layout: 'original',
  middleware({ $cookies, redirect }) {
    if ($cookies.get('agree') !== 'Agree') {
      return redirect('/')
    }
  },
  data() {
    return {
      delay: null,
      delayRule: [(value) => Number(value) > 0 || '0以上を入力してください'],

      loginflag: false,
      loginReferer: null,
      loginURL: null,
      loginOptions: null,
      loginmethods: ['GET', 'POST'],
      btnSetting: [
        { text: 'Full Scan', color: 'red' },
        { text: 'Quick Scan', color: 'primary' },
      ],
      scanOption: 'Full Scan',
      debugParam: null,
      landmarkFlag: false,
      landmarkNumber: 0,

      sitemap: {},
      approximateTime: 0,
    }
  },
  watch: {
    scanOption() {
      this.calcApproximateTime()
    },
    delay() {
      this.calcApproximateTime()
    },
    loginflag() {
      this.calcApproximateTime()
    },
  },
  created() {
    this.delay = this.$store.state.delay.delay

    this.loginReferer = this.$store.state.loginPath.loginRef
    this.loginURL = this.$store.state.loginPath.loginURL
    this.loginOptions = cloneDeep(this.$store.state.loginParams.loginParams)
    if (this.loginURL !== '') {
      this.loginflag = true
    }
  },
  mounted() {
    this.debugParam = new URLSearchParams(
      window.location.search.substring(1)
    ).get('debug')
    if (this.debugParam !== null) {
      this.landmarkFlag = true
    }
  },
  methods: {
    doScan() {
      this.$store.commit('delay/changeDelay', this.delay)

      const forms = new FormData()

      forms.append('delay', this.delay)
      forms.append('scanOption', this.scanOption)
      forms.append('LandmarkNumber', this.landmarkNumber)

      if (this.loginflag) {
        forms.append('loginReferer', this.loginReferer)
        this.$store.commit('loginPath/changeloginRef', this.loginReferer)

        forms.append('loginURL', this.loginURL)
        this.$store.commit('loginPath/changeloginURL', this.loginURL)

        this.$store.commit('loginParams/changeloginParams', this.loginOptions)
        for (const i in this.loginOptions) {
          forms.append('loginKey[]', this.loginOptions[i].key)
          forms.append('loginValue[]', this.loginOptions[i].value)
          forms.append('loginMethod[]', this.loginOptions[i].method)
        }
      }

      this.$axios
        .$post('/api/scan', forms)
        .then((response) => {
          this.$store.commit('scanstate/changestate', 'finish', { root: true })
          // this.transitionsitemap()
        })
        .catch((err) => {
          console.log('err:', err)
        })
    },
    transitionsitemap() {
      this.$router.push('/report')
    },
    async calcApproximateTime() {
      if (!this.sitemap.path) {
        await this.$axios
          .$get('/api/sitemap')
          .then((response) => {
            this.sitemap = response
          })
          .catch((err) => {
            console.log('err:', err)
          })
      }
      const msgNum = this.countMsg(this.sitemap)
      const cookieNum = this.countCookie(this.sitemap)
      const paramNum = this.countParam(this.sitemap)
      let accessNum = (msgNum * 5 + paramNum + cookieNum) * 315
      let accessTime = this.retrieveAccessTime(this.sitemap)
      accessTime += Number(this.delay)

      if (this.scanOption === 'Full Scan') {
        accessNum += (msgNum * 2 + paramNum + cookieNum) * msgNum
      }
      if (this.loginflag) {
        accessNum *= 2
      }
      accessTime += accessNum * 0.0001

      this.approximateTime = Math.round((accessTime * accessNum) / 60000)
    },
    retrieveAccessTime(node) {
      let accessTime
      if (node.messages[0]) {
        accessTime = node.messages[0].time
      }
      for (const i in node.children) {
        if (accessTime) {
          break
        }
        accessTime = this.retrieveAccessTime(node.children[i])
      }
      return accessTime
    },
    countMsg(node) {
      let msgNum = node.messages.length
      for (const i in node.children) {
        msgNum += this.countMsg(node.children[i])
      }
      return msgNum
    },
    countCookie(node) {
      let cookieNum = node.cookies.length
      for (const i in node.children) {
        cookieNum += this.countParam(node.children[i])
      }
      return cookieNum
    },
    countParam(node) {
      let paramNum = 0
      for (const i in node.messages) {
        paramNum += Object.keys(node.messages[i].getParams).length
        paramNum += Object.keys(node.messages[i].postParams).length
      }
      for (const i in node.children) {
        paramNum += this.countParam(node.children[i])
      }
      return paramNum
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
