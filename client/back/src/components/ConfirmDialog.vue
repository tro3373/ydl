<template>
  <div class="text-center">
    <v-dialog overlay-opacity="0.9" dark v-model="dialog" width="500">
      <v-card>
        <v-card-title class="headline">{{ title }}</v-card-title>
        <v-card-text>
          {{ message }}
          <br />
          <br />
          {{ subMessage }}
        </v-card-text>
        <v-divider></v-divider>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="error" outlined @click="onOk()">
            {{ buttonOk }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
export default {
  props: ['title', 'message', 'buttonOk'],
  data() {
    return {
      dialog: false,
      subMessage: '',
      handler: () => {},
    };
  },
  methods: {
    onOk() {
      this.dialog = false;
      this.handler();
    },
    show(subMessage, handler) {
      this.dialog = true;
      this.handler = handler;
      this.subMessage = subMessage;
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss"></style>
