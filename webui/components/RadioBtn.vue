<template>
  <div>
    <v-radio-group
      v-model="scanOption"
      v-for="(item, index) in _btnSetting"
      :key="index"
    >
      <v-row>
        <v-radio :label="item.text" :value="item.text" :color="item.color">
        </v-radio>
      </v-row>
    </v-radio-group>
  </div>
</template>

<script>
export default {
  model: { prop: 'propRadio', event: 'change' },
  props: { propRadio: String, btnSetting: Array },
  computed: {
    scanOption: {
      get() {
        return this.propRadio
      },
      set(value) {
        this.$emit('change', value)
      },
    },
    _btnSetting: {
      get() {
        return this.btnSetting
      },
      set(value) {
        this.$emit('change', value)
      },
    },
  },
}
</script>

<style></style>
