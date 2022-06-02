const { defineConfig } = require('@vue/cli-service');
module.exports = defineConfig({
  devServer: {
    // host: '0.0.0.0'
    host: `${process.env.EXTERNAL_IP ? process.env.EXTERNAL_IP : '0.0.0.0'}`,
  },
  // transpileDependencies: true,
  transpileDependencies: ['vuetify'],
});
