<template>
  <div>
    <v-btn @click="addLoginOptionFormData()">{{ this.btnText }} </v-btn>
  </div>
</template>

<script>
export default {
  props: { addform: Array, btnText: String, adddata: Object },
  computed: {
    _addform: {
      get() {
        return this.addform
      },
      set(value) {
        this.$emit('change', value)
      },
    },
  },
  methods: {
    addLoginOptionFormData() {
      const add = this.adddata
      this._addform.push(add)
    },
  },
}
</script>

<style></style>
