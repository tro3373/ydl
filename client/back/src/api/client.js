'use strict';

import Const from '../constants/constants.js';
import BaseClient from './baseClient.js';

const ENDPOINT = {
  LIST: {
    url: `${Const.API_PREFIX}`,
    method: 'get',
  },
  DL_REQUEST: {
    url: `${Const.API_PREFIX}`,
    method: 'post',
  },
};

const ApiClient = class ApiClient extends BaseClient {
  constructor() {
    super('ApiClient');
  }

  setUuid(uuid) {
    const headers = { 'x-uuid': uuid };
    this.setDefaultHeaders(headers);
  }

  async list(params) {
    const res = await super.request({
      ...ENDPOINT.LIST,
      params,
    });
    return res.data;
  }

  async downloadRequest(params) {
    const res = await super.request({
      ...ENDPOINT.DL_REQUEST,
      data: params,
    });
    return res.data;
  }

  async youtubeOembedInfo(key) {
    const url = `https://www.youtube.com/oembed?url=https://www.youtube.com/watch?v=${key}&format=json`;
    const res = await super.request({
      url,
      method: 'get',
    });
    return res.data;
  }
};
export default new ApiClient();
