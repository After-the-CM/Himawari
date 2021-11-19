<template>
  <div>
    <v-list-item v-for="(loginOption, index) in loginOptions" :key="index">
      <v-col :cols="methodcols">
        <v-select
          v-model="loginOption.method"
          label="メソッド"
          outlined
          :items="methods"
        ></v-select>
      </v-col>
      <v-text-field v-model="loginOption.key" label="LoginKey"></v-text-field>
      <v-text-field
        v-model="loginOption.value"
        label="LoginValue"
      ></v-text-field>
      <delete-form-btn :deleteform="loginOptions" btnText="削除" :i="index" />
    </v-list-item>
  </div>
</template>

<script>
import DeleteFormBtn from './DeleteFormBtn.vue'
export default {
  components: { DeleteFormBtn },
  model: { prop: 'propLoginOption', event: 'change' },
  props: { propLoginOption: Array, methods: Array, methodcols: Number },
  computed: {
    loginOptions: {
      get() {
        return this.propLoginOption
      },
      set(value) {
        this.$emit('change', value)
      },
    },
  },
}
</script>

<style></style>
