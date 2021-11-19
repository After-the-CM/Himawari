<template>
  <li>
    <div :class="{ bold: isFolder }">
      <span @click="toggle">
        {{ item.path }}
        <span v-if="isFolder">[{{ isOpen ? '-' : '+' }}]</span>
      </span>
    </div>
    <ul v-show="isOpen" v-if="isFolder">
      <SitemapNode
        v-for="(child, index) in item.children"
        :key="index"
        :item="child"
        class="item"
      ></SitemapNode>
    </ul>
  </li>
</template>

<script>
import SitemapNode from '~/components/SitemapNode'
export default {
  name: 'SitemapNode',
  components: {
    SitemapNode,
  },
  props: {
    item: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      isOpen: true,
    }
  },
  computed: {
    isFolder() {
      return this.item.children && this.item.children.length
    },
  },
  methods: {
    toggle() {
      if (this.isFolder) {
        this.isOpen = !this.isOpen
      }
    },
  },
}
</script>
