<template>
  <div>
    <v-text-field
      :label="labelText"
      v-model="inputtext"
      :rules="inputRule"
      :class="textClass"
      :id="textId"
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
    inputtext: {
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

<style></style>
