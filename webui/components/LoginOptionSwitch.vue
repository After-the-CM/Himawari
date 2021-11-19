<template>
  <div class="mx-3">
    <v-switch
      v-model="loginflag"
      :label="`ログインオプション: ` + doLogin()"
    ></v-switch>
  </div>
</template>

<script>
export default {
  model: { prop: 'propLoginFlag', event: 'change' },
  props: { propLoginFlag: Boolean },
  computed: {
    loginflag: {
      get() {
        return this.propLoginFlag
      },
      set(value) {
        this.$emit('change', value)
      },
    },
  },
  methods: {
    doLogin() {
      if (this.loginflag) {
        return 'ON'
      } else {
        return 'OFF'
      }
    },
  },
}
</script>

<style></style>
