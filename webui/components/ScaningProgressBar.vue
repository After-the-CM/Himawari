<template>
  <div v-if="flag">
    <v-progress-linear
      indeterminate
      color="cyan"
      fixed
      style="z-index: 4"
    ></v-progress-linear>
  </div>
</template>

<script>
export default {
  props: { flag: Boolean },
}
</script>

<style></style>
