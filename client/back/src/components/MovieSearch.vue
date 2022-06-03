<template>
  <v-container>
    <v-layout wrap row>
      <v-flex xs12 sm12 md4 class="pa-1">
        <v-card>
          <v-progress-linear
            :active="loading"
            :indeterminate="loading"
            top
            color="blue accent-4"
          ></v-progress-linear>

          <v-card-title>
            <!-- YouTube Downloader -->
            <v-text-field
              v-model="url"
              :rules="rules.url"
              label="youtube url or id を入力"
              @focus="oembedGuard = false"
              @blur="oembedGuard = true"
              clearable
            >
            </v-text-field>

            <v-btn
              color="red"
              dark
              class="ma-2 white--text"
              href="https://www.youtube.com/"
              target="_blank"
            >
              探しにいく
              <v-icon>mdi-youtube</v-icon>
            </v-btn>

            <v-chip class="ma-2" input-value="true" filter v-if="youtubeId">
              ID: {{ youtubeId }}
            </v-chip>
          </v-card-title>

          <v-card-text>
            <EmbedPlayer :url="embedUrl" />
          </v-card-text>

          <v-divider></v-divider>

          <v-card-text>
            <v-text-field v-model="title" :rules="rules.title" label="タイトル" clearable>
            </v-text-field>
            <v-text-field v-model="artist" :rules="rules.artist" label="アーティスト" clearable>
            </v-text-field>
            <v-text-field v-model="album" :rules="rules.album" label="アルバム" clearable>
            </v-text-field>
            <v-text-field v-model="genre" :rules="rules.genre" label="ジャンル" clearable>
            </v-text-field>
          </v-card-text>

          <v-card-actions>
            <v-spacer></v-spacer>
            <!-- <v-radio-group v-model="type" row>            -->
            <!--   <v-radio label="mp3" value="mp3"></v-radio> -->
            <!--   <v-radio label="mp4" value="mp4"></v-radio> -->
            <!-- </v-radio-group>                              -->
            <!-- elevation="4" -->
            <v-btn
              class="ma-2"
              color="error"
              fab
              :loading="loading"
              :disabled="loading"
              @click="submit"
            >
              <v-icon>mdi-coffee-to-go</v-icon>
              <template v-slot:loader>
                <span class="custom-loader">
                  <v-icon light>mdi-cached</v-icon>
                </span>
              </template>
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-flex>

      <v-flex xs12 sm12 md8 class="pa-1">
        <v-container fill-height v-if="!doneList.length">
          <v-spacer></v-spacer>
          No Result
          <v-spacer></v-spacer>
        </v-container>
        <v-card v-if="doneList.length">
          <v-toolbar color="primary" dark>
            <v-toolbar-title>Files</v-toolbar-title>
          </v-toolbar>

          <v-list two-line>
            <v-list-item
              v-for="done in doneList"
              :key="done.url"
              @click.stop="onItemSelected(done)"
            >
              <v-list-item-avatar size="80" width="160" rounded>
                <v-img :src="done.thumbnail"></v-img>
              </v-list-item-avatar>
              <v-list-item-content>
                <v-list-item-title v-text="done.tag.title"></v-list-item-title>
                <v-list-item-subtitle v-text="done.tag.artist"></v-list-item-subtitle>
              </v-list-item-content>
              <v-list-item-action>
                <v-btn icon>
                  <v-icon @click.stop="download(done, 1)" color="red">mdi-movie</v-icon>
                </v-btn>
                <div class="v-icon notranslate mdi theme--light" style="font-size: 0.2rem">
                  {{ humanSize(done.movieSize) }}
                </div>
                <v-btn icon>
                  <v-icon @click.stop="download(done, 0)" color="red">mdi-music</v-icon>
                </v-btn>
                <div class="v-icon notranslate mdi theme--light" style="font-size: 0.2rem">
                  {{ humanSize(done.audioSize) }}
                </div>
              </v-list-item-action>
            </v-list-item>
          </v-list>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import _ from 'lodash';
import qs from 'query-string';
import EmbedPlayer from '@/components/EmbedPlayer.vue';
import client from '@/api/client.js';
import youtubeApilient from '@/api/youtubeApiClient.js';
import Const from '../constants/constants.js';
const { mapActions: mapActionsDone, mapGetters: mapGettersDone } = createNamespacedHelpers('done');
export default {
  name: 'MovieSearch',
  components: {
    EmbedPlayer,
  },
  watch: {
    loader() {
      const l = this.loader;
      this[l] = !this[l];
      setTimeout(() => (this[l] = false), 3000);
      this.loader = null;
    },
    inputCache() {
      const data = {
        url: this.url,
        title: this.title,
        artist: this.artist,
        album: this.album,
        genre: this.genre,
      };
      console.debug('==> Caching to local storage..', data);
      localStorage.setItem(Const.LOCAL_STRAGE_KEY.CACHE, JSON.stringify(data));
    },
    youtubeId() {
      this.onYoutubeIdChanged();
    },
  },
  data() {
    const inputInit = {
      url: '',
      title: '',
      artist: '',
      album: '',
      genre: '',
    };
    const input = JSON.parse(
      localStorage.getItem(Const.LOCAL_STRAGE_KEY.CACHE) || JSON.stringify(inputInit)
    );
    return {
      ...input,
      rules: {
        url: [
          () => {
            const res = this.isValidId || 'Invalid id or url.';
            return res;
          },
        ],
      },
      // type: 'mp3',
      loader: null,
      loading: false,
      oembedGuard: true,
    };
  },
  computed: {
    ...mapGettersDone(['doneList']),
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
    },
    inputCache() {
      return `${this.url}${this.title}${this.artist}${this.album}${this.genre}`;
    },
  },
  async mounted() {
    await this.getDone();
  },
  methods: {
    ...mapActionsDone(['getDone']),
    validate(id) {
      const _id = '' + id;
      const res = _id.match(/^[a-zA-Z0-9_-]{11}$/);
      return !!res;
    },
    async onYoutubeIdChanged() {
      if (!this.youtubeId || this.oembedGuard) {
        return;
      }
      const res = await youtubeApilient.getOembedInfo(this.youtubeId);
      this.title = res.title;
      this.artist = res.author_name;
      this.album = res.author_name;
      this.genre = '';
    },
    async submit() {
      this.loader = 'loading';
      if (!this.isValidId) {
        return;
      }
      // const res = await client.list();
      // console.log({ res });
      const res = await client.downloadRequest({
        url: this.youtubeId,
        tag: {
          title: this.title,
          artist: this.artist,
          album: this.album,
          genre: this.genre,
        },
      });
      console.log({ res });
    },
    onItemSelected(done) {
      console.log({ done });
      this.url = done.url;
      this.title = done.tag.title;
      this.artist = done.tag.artist;
      this.album = done.tag.album;
      this.genre = done.tag.genre;
    },
    download(done, movie) {
      const title = done.tag.title;
      let url = done.audio;
      if (movie) {
        url = done.movie;
      }
      const ext = this.ext(url);
      url = `${url}?f=${title}.${ext}`;
      window.open(url, '_self');
    },
    ext(file) {
      return file.substr(file.lastIndexOf('.') + 1);
    },
    humanSize(size) {
      if (!size || size === -1) return '';
      return `${(size / 1024 / 1024).toFixed(1)}Mb`;
    },
  },
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