'use strict';

import Const from '../constants/constants.js';
import BaseClient from './baseClient.js';

const ApiClient = class ApiClient extends BaseClient {
  constructor() {
    super('ApiClient');
    const apiPrefix = this.getApiPrefix();
    this.ENDPOINT = {
      LIST: {
        url: apiPrefix,
        method: 'get',
      },
      DL_REQUEST: {
        url: apiPrefix,
        method: 'post',
      },
      DELETE_REQUEST: {
        url: `${apiPrefix}/{key}`,
        method: 'delete',
      },
    };
  }

  getApiPrefix() {
    let apiPrefix = Const.API_PREFIX;
    if (apiPrefix) {
      return apiPrefix;
    }
    const host = location.host;
    const url = location.toString();
    const scheme = url.split(':')[0];
    return `${scheme}://${host}/api`;
  }

  setUuid(uuid) {
    const headers = { 'x-uuid': uuid };
    this.setDefaultHeaders(headers);
  }

  async list(params) {
    const res = await super.request({
      ...this.ENDPOINT.LIST,
      params,
    });
    return res.data;
  }

  async downloadRequest(params) {
    const res = await super.request({
      ...this.ENDPOINT.DL_REQUEST,
      data: params,
    });
    return res.data;
  }

  async deleteRequest(key) {
    const res = await super.request(
      {
        ...this.ENDPOINT.DELETE_REQUEST,
      },
      { key }
    );
    return res.data;
  }

  // async youtubeOembedInfo(key) {
  //   const url = `https://www.youtube.com/oembed?url=https://www.youtube.com/watch?v=${key}&format=json`;
  //   const res = await super.request({
  //     url,
  //     method: 'get',
  //   });
  //   return res.data;
  // }
};
export default new ApiClient();
