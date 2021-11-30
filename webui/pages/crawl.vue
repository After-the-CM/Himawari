<template>
  <v-app>
    <scaning-progress-bar :flag="crawlingFlag" />
    <div v-if="alertFlag">
      <v-alert
        color="pink darken-1"
        border="bottom"
        dark
        dismissible
        style="position: fixed; width: 100%; z-index: 1000"
        ><div class="font-weight-bold">Something error.</div></v-alert
      >
    </div>
    <v-card width="50%" class="mx-auto mt-5 mb-12" style="position: relative">
      <v-toolbar color="gray" dark flat>
        <v-toolbar-title>Crawler</v-toolbar-title>

        <template v-slot:extension>
          <v-tabs v-model="tab" centered>
            <v-tabs-slider color="yellow"></v-tabs-slider>

            <v-tab
              v-for="item in items"
              :key="item"
              @click="
                flagOn(item)
                printFlag()
              "
            >
              {{ item }}
            </v-tab>
          </v-tabs>
        </template>
      </v-toolbar>
      <v-tabs-items v-model="tab">
        <v-tab-item v-for="(item, i) in items" :key="i">
          <v-card flat v-if="isflag">
            <v-row>
              <v-col cols="9">
                <input-text
                  v-model="url"
                  labelText="URL"
                  :inputRule="crawlURLRule"
                  textId="url"
                  textClass="mx-5"
                />
              </v-col>
              <v-col cols="3">
                <input-number
                  v-model="delay"
                  labelText="delay(ms)"
                  :inputRule="delayRule"
                  textId="delay"
                  textClass="mx-5"
                />
              </v-col>
            </v-row>

            <login-option-switch v-model="loginflag" />
            <v-card v-if="loginflag">
              <v-card-title>ログイン情報入力</v-card-title>
              <input-text
                v-model="loginReferer"
                labelText="LoginフォームがあるURL(Referer)"
                textId="loginref"
                textClass="mx-5"
              />
              <input-text
                v-model="loginURL"
                labelText="Loginリクエストの送信先"
                textId="loginurl"
                textClass="mx-5"
              />
              <v-list>
                <v-row>
                  <v-col cols="12">
                    <login-option
                      v-model="loginOptions"
                      :methods="loginmethods"
                      :methodcols="3"
                    />
                  </v-col>
                </v-row>
              </v-list>
              <add-form-btn
                :addform="loginOptions"
                btnText="ログインパラメータ追加"
                :adddata="{ key: '', value: '', method: 'POST' }"
              />
            </v-card>
            <v-expansion-panels>
              <v-expansion-panel>
                <v-expansion-panel-header>
                  Crawl時に代入するパラメーター設定はこちら
                </v-expansion-panel-header>
                <v-expansion-panel-content>
                  <v-list>
                    <v-list-item
                      v-for="(data, index) in formdatas"
                      :key="index"
                      style="position: relative"
                    >
                      <v-text-field
                        v-model="data.name"
                        label="key"
                      ></v-text-field>
                      <v-text-field
                        v-model="data.value"
                        label="value"
                      ></v-text-field>

                      <delete-form-btn
                        :deleteform="formdatas"
                        btnText="削除"
                        :i="index"
                      />
                    </v-list-item>
                  </v-list>
                  <v-btn @click="addForm()">パラメーター追加</v-btn>
                </v-expansion-panel-content>
              </v-expansion-panel>
            </v-expansion-panels>
          </v-card>
          <v-card flat v-else>
            <v-file-input
              id="Json"
              class="mx-5"
              v-model="file"
              accept=".json"
              :loading="inputfileflag"
              @change="finishUpload"
              @focus="onLoding"
              @blur="offLoding"
            ></v-file-input>
          </v-card>
        </v-tab-item>
      </v-tabs-items>
      <v-btn
        v-if="isflag"
        :disabled="url === '' || delay === '' || delay < 0"
        rounded
        absolute
        right
        class="mt-3 text-capitalize"
        @click="doCrawl()"
      >
        Crawl
      </v-btn>

      <v-btn
        v-else
        :disabled="fileUploadFlag"
        rounded
        absolute
        right
        class="mt-3 text-capitalize"
        @click="doUpload()"
      >
        Upload JSON
      </v-btn>
    </v-card>
  </v-app>
</template>

<script>
import { cloneDeep } from 'lodash'
import InputText from '~/components/InputText.vue'
import InputNumber from '~/components/InputNumber.vue'

export default {
  components: { InputText, InputNumber },
  layout: 'original',
  middleware({ $cookies, redirect }) {
    if ($cookies.get('agree') !== 'Agree') {
      return redirect('/')
    }
  },
  data() {
    return {
      tab: null,
      items: ['New', 'Import SitemapJson'],
      isflag: true,
      isURL: false,
      url: '',
      crawlURLRule: [(value) => !!value || '必須項目です'],
      delay: null,
      delayRule: [(value) => Number(value) > 0 || '0以上を入力してください'],

      formdatas: null,
      file: null,
      checkedScanOption: [],

      loginflag: false,
      loginReferer: '',
      loginURL: '',

      loginOptions: null,
      loginmethods: ['GET', 'POST'],

      inputfileflag: false,
      fileUploadFlag: true,

      crawlingFlag: false,
      alertFlag: false,
    }
  },
  created() {
    this.formdatas = cloneDeep(this.$store.state.crawlParams.crawlParams)
    this.url = this.$store.state.crawlURL.crawlURL
    this.delay = this.$store.state.delay.delay

    this.loginReferer = this.$store.state.loginPath.loginRef
    this.loginURL = this.$store.state.loginPath.loginURL

    this.loginOptions = cloneDeep(this.$store.state.loginParams.loginParams)
  },
  methods: {
    flagOn(item) {
      if (item === 'New') {
        this.isflag = true
      } else if (item === 'Import SitemapJson') {
        this.isflag = false
        // New から Sitemap Input のページに切り替わった際にログインオプションの欄をOFFに
        this.loginflag = false
      }
    },
    printFlag() {
      console.log('ok')
    },
    doCrawl() {
      this.crawlingFlag = true
      this.$store.commit('crawlParams/changecrawlParams', this.formdatas)
      this.$store.commit('crawlURL/changecrawlURL', this.url)
      this.$store.commit('delay/changeDelay', this.delay)

      const forms = new FormData()

      forms.append('url', this.url)
      forms.append('delay', this.delay)

      if (this.loginflag) {
        forms.append('loginReferer', this.loginReferer)
        this.$store.commit('loginPath/changeloginRef', this.loginReferer)

        forms.append('loginURL', this.loginURL)
        this.$store.commit('loginPath/changeloginURL', this.loginURL)

        this.$store.commit('loginParams/changeloginParams', this.loginOptions)

        for (const i in this.loginOptions) {
          forms.append('loginKey[]', this.loginOptions[i].key)
          console.log(this.loginOptions[i].key)
          forms.append('loginValue[]', this.loginOptions[i].value)
          forms.append('loginMethod[]', this.loginOptions[i].method)
        }
      }

      for (const j in this.formdatas) {
        // names.push(this.formdatas[i].name)
        forms.append('name[]', this.formdatas[j].name)
        forms.append('value[]', this.formdatas[j].value)
      }

      this.$axios
        .$post('/api/crawl', forms)
        .then((response) => {
          console.log(this.url)
          console.log(response)

          this.transitionsitemap()
        })
        .catch((err) => {
          this.alertFlag = true
          console.log(err)
        })
      // const params = new URLSearchParams()
    },
    doUpload() {
      const data = new FormData()
      data.append('sitemap', this.file, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })
      this.$axios
        .$post('/sitemap/upload', data)
        .then((response) => {
          console.log(this.url)
          console.log(response)
          this.transitionsitemap()
        })
        .catch((err) => {
          console.log('err:', err)
        })
    },
    addForm() {
      const adddata = {
        name: '',
        value: '',
      }
      this.formdatas.push(adddata)
    },
    deleteForm(index) {
      this.formdatas.splice(index, 1)
    },
    onLoding() {
      this.inputfileflag = true
    },
    offLoding() {
      this.inputfileflag = false
    },
    finishUpload() {
      this.inputfileflag = false
      this.fileUploadFlag = false
    },
    transitionsitemap() {
      this.$router.push('/scan')
    },
  },
}
</script>
