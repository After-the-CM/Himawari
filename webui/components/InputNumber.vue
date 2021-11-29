<template>
  <div>
    <v-text-field
      :label="labelText"
      v-model="inputnum"
      :rules="inputRule"
      class="align-right"
      :class="textClass"
      :id="textId"
      type="number"
    ></v-text-field>
  </div>
</template>

<script>
export default {
  model: { prop: 'propinput', event: 'change' },
  props: {
    propinput: String,
    labelText: String,
    inputRule: Array,
    textId: String,
    textClass: String,
  },
  computed: {
    inputnum: {
      get() {
        return this.propinput
      },
      set(value) {
        this.$emit('change', value)
      },
    },
  },
}
</script>

<style>
.align-right input {
  text-align: right;
}
</style>
