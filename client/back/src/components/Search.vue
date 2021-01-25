<template>
  <v-container>
    <v-row class="text-center">
      <v-col cols="12">
        <v-card>
          <v-progress-linear
            :active="loading"
            :indeterminate="loading"
            top
            color="blue accent-4"></v-progress-linear>

          <v-card-title>
            <!-- YouTube Downloader -->
            <v-text-field
              v-model="url"
              :rules="rules.url"
              label="youtube url or id を入力"
              clearable>
            </v-text-field>

            <v-btn
              color="primary"
              dark
              class="ma-2 white--text"
              href="https://www.youtube.com/"
              target="_blank">
              探しにいく
              <v-icon>mdi-youtube</v-icon>
            </v-btn>

            <v-chip class="ma-2" input-value="true" filter v-if="youtubeId">
              ID: {{ youtubeId }}
            </v-chip>
          </v-card-title>

          <v-card-text>
            <player :url="embedUrl" />
          </v-card-text>

          <v-divider></v-divider>

          <v-card-actions>
            <v-spacer></v-spacer>
            <v-radio-group v-model="type" row>
              <v-radio label="mp3" value="mp3"></v-radio>
              <v-radio label="mp4" value="mp4"></v-radio>
            </v-radio-group>
            <!-- elevation="4" -->
            <v-btn
              class="ma-2"
              color="error"
              fab
              :loading="loading"
              :disabled="loading"
              @click="submit">
              <v-icon>mdi-coffee-to-go</v-icon>
              <template v-slot:loader>
                <span class="custom-loader">
                  <v-icon light>mdi-cached</v-icon>
                </span>
              </template>
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import _ from 'lodash';
import qs from 'query-string';
import Player from '@/components/Player.vue';
import client from '@/api/client.js';
const LOCAL_STRAGE_KEY_PREV_ID = 'local_strage_key_prev_id';
export default {
  name: 'Search',
  components: {
    Player
  },

  watch: {
    loader() {
      const l = this.loader;
      this[l] = !this[l];
      setTimeout(() => this[l] = false, 3000);
      this.loader = null;
    }
  },

  data() {
    return {
      url: localStorage.getItem(LOCAL_STRAGE_KEY_PREV_ID) || '',
      rules: {
        url: [
          () => {
            const res = this.isValidId || 'Invalid id or url.';
            return res;
          }
        ]
      },
      type: 'mp3',
      loader: null,
      loading: false
    };
  },
  computed: {
    youtubeId() {
      const val = this.url;
      if (_.isEmpty(val)) return '';
      if (this.validate(val)) return val;
      if (val.match(/^http.*\/embed\/.*/)) {
        const _id = val.replace(/^http.*\/embed\//, '');
        if (this.validate(_id)) return _id;
      }
      const list = val.split('?');
      const data = qs.parse(list[1]);
      const id = data.v;
      if (this.validate(id)) return id;
      return '';
    },
    isValidId() {
      return this.validate(this.youtubeId);
    },
    embedUrl() {
      if (!this.isValidId) {
        return '';
      }
      return `https://www.youtube.com/embed/${this.youtubeId}`;
    }
  },
  methods: {
    validate(id) {
      const _id = '' + id;
      const res = _id.match(/^[a-zA-Z0-9_-]{11}$/);
      return !!res;
    },
    async submit() {
      this.loader = 'loading';
      if (this.isValidId) {
        localStorage.setItem(LOCAL_STRAGE_KEY_PREV_ID, this.youtubeId);
        await client.list();
      }
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.custom-loader {
  animation: loader 1s infinite;
  display: flex;
}
@-moz-keyframes loader {
  from {
    transform: rotate(0);
  }
  to {
    transform: rotate(360deg);
  }
}
@-webkit-keyframes loader {
  from {
    transform: rotate(0);
  }
  to {
    transform: rotate(360deg);
  }
}
@-o-keyframes loader {
  from {
    transform: rotate(0);
  }
  to {
    transform: rotate(360deg);
  }
}
@keyframes loader {
  from {
    transform: rotate(0);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
