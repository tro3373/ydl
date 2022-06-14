<template>
  <v-app>
    <v-main>
      <MovieSearch />
    </v-main>
  </v-app>
</template>

<script>
import MovieSearch from '@/components/MovieSearch.vue';

export default {
  name: 'App',
  components: {
    MovieSearch,
  },
  data: () => ({}),
  mounted() {
    this.$vuetify.theme.dark = true;
  },
};
</script>
