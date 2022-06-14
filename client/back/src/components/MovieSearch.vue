<template>
  <v-container>
    <v-layout wrap row>
      <v-flex xs12 sm12 md4 class="pa-1">
        <v-card>
          <v-progress-linear
            :active="serverQueueing"
            :indeterminate="true"
            top
            color="blue accent-4"
          ></v-progress-linear>

          <v-form v-model="valid" ref="form">
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
            </v-card-title>

            <v-card-actions>
              <v-chip class="ma-2" input-value="true" filter v-if="youtubeId" @click="seeInYoutube">
                ID: {{ youtubeId }}
              </v-chip>
              <v-spacer></v-spacer>
              <v-btn class="ma-2" color="error" fab @click="findInYoutube">
                <v-icon>mdi-magnify</v-icon>
                <template v-slot:loader>
                  <span class="custom-loader">
                    <v-icon light>mdi-cached</v-icon>
                  </span>
                </template>
              </v-btn>

              <v-btn
                class="ma-2"
                color="error"
                fab
                :loading="serverQueueing"
                :disabled="serverQueueing || !valid"
                @click="submit"
              >
                <span class="material-icons">save</span>
                <template v-slot:loader>
                  <span class="custom-loader">
                    <v-icon light>mdi-cached</v-icon>
                  </span>
                </template>
              </v-btn>
            </v-card-actions>

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
              <v-text-field v-model="uuid" :rules="rules.uuid" label="uuid" clearable>
              </v-text-field>
            </v-card-text>
          </v-form>
        </v-card>
      </v-flex>

      <v-flex xs12 sm12 md8 class="pa-1">
        <v-container fill-height v-if="!requestResults.length">
          <v-spacer></v-spacer>
          No Result
          <v-spacer></v-spacer>
        </v-container>
        <v-card v-if="requestResults.length">
          <!-- <v-toolbar color="primary" dark>           -->
          <!--   <v-toolbar-title>Files</v-toolbar-title> -->
          <!-- </v-toolbar>                               -->

          <v-list two-line>
            <v-progress-linear
              :active="serverDoing"
              :indeterminate="true"
              top
              color="blue accent-4"
            ></v-progress-linear>
            <v-list-item
              v-for="rr in requestResults"
              :key="`${rr.url}.${rr.doing}`"
              :disabled="rr.doing"
              @click.stop="onItemSelected(rr)"
            >
              <v-list-item-avatar size="80" width="160" rounded>
                <v-progress-circular
                  indeterminate
                  color="grey"
                  v-show="!rr.thumbnail"
                ></v-progress-circular>
                <v-img :src="rr.thumbnail" v-show="rr.thumbnail"></v-img>
              </v-list-item-avatar>
              <v-list-item-content>
                <v-list-item-title v-text="rr.tag.title"></v-list-item-title>
                <v-list-item-subtitle v-text="rr.tag.artist"></v-list-item-subtitle>
              </v-list-item-content>
              <v-list-item-action>
                <v-btn icon>
                  <v-icon :disabled="rr.doing" @click.stop="download(rr, 1)" color="red"
                    >mdi-movie</v-icon
                  >
                </v-btn>
                <div class="v-icon notranslate mdi theme--light text-caption">
                  {{ humanSize(rr.movieSize) }}
                </div>
                <v-btn icon>
                  <v-icon :disabled="rr.doing" @click.stop="download(rr, 0)" color="red"
                    >mdi-music</v-icon
                  >
                </v-btn>
                <div class="v-icon notranslate mdi theme--light text-caption">
                  {{ humanSize(rr.audioSize) }}
                </div>
              </v-list-item-action>
              <v-list-item-action>
                <v-btn color="red lighten-2" icon>
                  <v-icon :disabled="rr.doing" @click.stop="confirmDelete(rr)">
                    mdi-delete-forever
                  </v-icon>
                </v-btn>
              </v-list-item-action>
            </v-list-item>
          </v-list>
        </v-card>
      </v-flex>
      <ConfirmDialog
        ref="confirmDialog"
        title="削除"
        message="保存されたデータを消す"
        buttonOk="消す"
      />
    </v-layout>
  </v-container>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import _ from 'lodash';
import qs from 'query-string';
import EmbedPlayer from '@/components/EmbedPlayer.vue';
import ConfirmDialog from '@/components/ConfirmDialog.vue';
import util from '../util/util.js';
import client from '@/api/client.js';
import youtubeApilient from '@/api/youtubeApiClient.js';
import Const from '../constants/constants.js';
const { mapActions: mapActionsRequestResults, mapGetters: mapGettersRequestResults } =
  createNamespacedHelpers('requestResults');
export default {
  name: 'MovieSearch',
  components: {
    EmbedPlayer,
    ConfirmDialog,
  },
  watch: {
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
    uuid() {
      if (!this.uuid) return;
      localStorage.setItem(Const.LOCAL_STRAGE_KEY.UUID, this.uuid);
      this.getRequestResultsWithUuid();
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
    const uuid = localStorage.getItem(Const.LOCAL_STRAGE_KEY.UUID) || util.uuid();
    return {
      ...input,
      uuid,
      rules: {
        url: [
          () => {
            return this.isValidId || 'Invalid id or url.';
          },
        ],
        uuid: [
          () => {
            return !!this.uuid || 'Input any your id';
          },
        ],
      },
      // type: 'mp3',
      serverQueueing: false,
      oembedGuard: true,
      valid: this.valid,
    };
  },
  computed: {
    ...mapGettersRequestResults(['requestResults']),
    youtubeId() {
      const val = this.url;
      if (_.isEmpty(val)) return '';
      if (this.validate(val)) return val;
      if (val.match(/^http.*\/youtu.be\/.*/)) {
        const _id = val.replace(/^http.*\/youtu.be\//, '');
        if (this.validate(_id)) return _id;
      }
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
    serverDoing() {
      return !!this.requestResults.find((r) => r.doing);
    },
  },
  async mounted() {
    // this.$refs.form.validate(); // for submit icon not enable
    this.getRequestResultsWithUuid();
  },
  methods: {
    ...mapActionsRequestResults(['getRequestResults']),
    getRequestResultsWithUuid() {
      client.setUuid(this.uuid);
      this.getRequestResults();
    },
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
    seeInYoutube() {
      if (util.isEmpty(this.youtubeId)) {
        return;
      }
      let url = `https://www.youtube.com/watch?v=${this.youtubeId}`;
      window.open(url, '_blank');
    },
    findInYoutube() {
      let url = 'https://www.youtube.com';
      let q = this.title || this.youtubeId;
      if (!util.isEmpty(q)) {
        url = `https://www.youtube.com/results?search_query=${q}`;
      }
      window.open(url, '_blank');
    },
    async submit() {
      this.serverQueueing = true;
      if (!this.isValidId) {
        return;
      }
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
      console.log('==> Start interval!');
      const intervalId = setInterval(async () => {
        await this.getRequestResults();
        if (this.serverDoing) {
          this.serverQueueing = false;
          return;
        }
        console.log('==> Clear interval!');
        clearInterval(intervalId);
      }, 5000);
    },
    onItemSelected(rr) {
      console.log({ rr });
      this.url = rr.url;
      this.title = rr.tag.title;
      this.artist = rr.tag.artist;
      this.album = rr.tag.album;
      this.genre = rr.tag.genre;
    },
    download(rr, movie) {
      const title = rr.tag.title;
      let url = rr.audio;
      if (movie) {
        url = rr.movie;
      }
      const ext = this.ext(url);
      url = `${url}?f=${title}.${ext}`;
      window.open(url, '_self');
    },
    confirmDelete(rr) {
      const youtubeId = rr.url;
      const movieTitle = rr.tag.title;
      this.$refs.confirmDialog.show(movieTitle, async () => {
        await this.deleteRequest(youtubeId);
      });
    },
    async deleteRequest(youtubeId) {
      console.log(youtubeId);
      const res = await client.deleteRequest(youtubeId);
      console.log({ res });
      await this.getRequestResults();
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
