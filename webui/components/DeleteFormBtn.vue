<template>
  <v-btn @click="deleteLoginOptionFormData(i)">{{ btnText }}</v-btn>
</template>

<script>
export default {
  props: { deleteform: Array, btnText: String, i: Number },
  computed: {
    _deleteform: {
      get() {
        return this.deleteform
      },
      set(value) {
        this.$emit('change', value)
      },
    },
  },
  methods: {
    deleteLoginOptionFormData(i) {
      this._deleteform.splice(i, 1)
    },
  },
}
</script>

<style></style>
