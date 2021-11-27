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
              <sitemap />
            </v-card>
          </v-col>
          <v-spacer></v-spacer>
          <v-col cols="7" style="border-left: solid 1px gray">
            <v-card flat>
              <v-card-title>ScanOption</v-card-title>
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
                  labelText="LoginフォームがあるURL"
                  textId="url"
                  textClass="mx-5"
                />
                <input-text
                  v-model.trim="loginURL"
                  labelText="Loginリクエストの送信先"
                  textId="url"
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
              <div v-if="randmarkFlag">
                <v-text-field
                  v-model="randmarkNumber"
                  label="RandMark(default : 0)"
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
import InputText from '~/components/InputText.vue'
import LoginOption from '~/components/LoginOption.vue'
import LoginOptionSwitch from '~/components/LoginOptionSwitch.vue'
import OutOfScope from '~/components/OutOfScope.vue'
import InputNumber from '~/components/InputNumber.vue'

export default {
  components: {
    LoginOptionSwitch,
    LoginOption,
    OutOfScope,
    InputText,
    InputNumber,
  },
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
      randmarkFlag: false,
      randmarkNumber: 0,
    }
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
      this.randmarkFlag = true
    }
  },
  methods: {
    doScan() {
      this.$store.commit('delay/changeDelay', this.delay)

      const forms = new FormData()

      forms.append('delay', this.delay)
      forms.append('scanOption', this.scanOption)
      forms.append('RandmarkNumber', this.randmarkNumber)

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
  },
}
</script>

<style></style>
