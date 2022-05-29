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
              clearable
            >
            </v-text-field>

            <v-btn
              color="primary"
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
            <player :url="embedUrl" />
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
        <v-card>
          <v-toolbar color="primary" dark>
            <!-- <v-app-bar-nav-icon></v-app-bar-nav-icon> -->
            <v-toolbar-title>Files</v-toolbar-title>
            <!-- <v-spacer></v-spacer> -->
            <!-- <v-btn icon>                   -->
            <!--   <v-icon>mdi-magnify</v-icon> -->
            <!-- </v-btn>                       -->

            <!-- <v-btn icon>                       -->
            <!--   <v-icon>mdi-view-module</v-icon> -->
            <!-- </v-btn>                           -->
          </v-toolbar>

          <!-- <v-list subheader two-line> -->
          <v-list two-line>
            <!-- <v-subheader inset>Folders</v-subheader> -->
            <v-list-item v-for="done in doneList" :key="done.url">
              <v-list-item-avatar>
                <v-icon class="grey lighten-1" dark>
                  mdi-folder
                </v-icon>
              </v-list-item-avatar>
              <v-list-item-content>
                <v-list-item-title v-text="done.tag.title"></v-list-item-title>
                <v-list-item-subtitle v-text="done.tag.artist"></v-list-item-subtitle>
              </v-list-item-content>
              <v-list-item-action>
                <v-btn icon>
                  <v-icon color="grey lighten-1">mdi-information</v-icon>
                </v-btn>
              </v-list-item-action>
            </v-list-item>
            <!-- <v-divider inset></v-divider>                                            -->
            <!-- <v-subheader inset>Files</v-subheader>                                   -->
            <!-- <v-list-item v-for="file in files" :key="file.title">                    -->
            <!--   <v-list-item-avatar>                                                   -->
            <!--     <v-icon :class="file.color" dark v-text="file.icon"></v-icon>        -->
            <!--   </v-list-item-avatar>                                                  -->
            <!--   <v-list-item-content>                                                  -->
            <!--     <v-list-item-title v-text="file.title"></v-list-item-title>          -->
            <!--     <v-list-item-subtitle v-text="file.subtitle"></v-list-item-subtitle> -->
            <!--   </v-list-item-content>                                                 -->
            <!--   <v-list-item-action>                                                   -->
            <!--     <v-btn icon>                                                         -->
            <!--       <v-icon color="grey lighten-1">mdi-information</v-icon>            -->
            <!--     </v-btn>                                                             -->
            <!--   </v-list-item-action>                                                  -->
            <!-- </v-list-item>                                                           -->
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
import Player from '@/components/Player.vue';
import client from '@/api/client.js';
const { mapActions: mapActionsDone, mapGetters: mapGettersDone } = createNamespacedHelpers('done');
const LOCAL_STRAGE_KEY_CACHE = 'local_strage_key_cache';
export default {
  name: 'Search',
  components: {
    Player,
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
      localStorage.setItem(LOCAL_STRAGE_KEY_CACHE, JSON.stringify(data));
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
      localStorage.getItem(LOCAL_STRAGE_KEY_CACHE) || JSON.stringify(inputInit)
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
    async submit() {
      this.loader = 'loading';
      if (!this.isValidId) {
        return;
      }
      // const res = await client.list();
      // console.log({ res });
      const res = await client.download({
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
