<template>
  <v-container class="responsive-style">
    <iframe
      width="560"
      height="315"
      :src="url"
      frameborder="0"
      allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
      allowfullscreen
    ></iframe>
  </v-container>
</template>

<script>
export default {
  name: 'EmpbedPlayer',
  props: {
    url: String,
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.responsive-style {
  position: relative;
  width: 100%;
  height: 0;
  padding-top: 50%;
}
.responsive-style iframe {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}
</style>
