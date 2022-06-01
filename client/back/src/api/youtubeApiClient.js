'use strict';

import BaseClient from './baseClient.js';

const YoutubeApiClient = class YoutubeApiClient extends BaseClient {
  constructor() {
    super('YoutubeApiClient');
  }

  async getOembedInfo(key) {
    const url = `https://www.youtube.com/oembed?url=https://www.youtube.com/watch?v=${key}&format=json`;
    const res = await super.request({
      url,
      method: 'get',
    });
    return res.data;
  }
};
export default new YoutubeApiClient();
