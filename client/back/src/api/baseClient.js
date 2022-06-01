'use strict';

import axios from 'axios';
import FormData from 'form-data';

import util from '../util/util.js';

export default class BaseClient {
  constructor(name) {
    this.name = name;
    this.defaultHeaders = {};
  }

  setDefaultHeaders(headers) {
    this.defaultHeaders = headers || {};
  }

  formatUrl(url, pathParams) {
    let ret = `${url}`;
    for (const key in pathParams) {
      const val = pathParams[key];
      ret = ret.replace(new RegExp(`{${key}}`), val);
    }
    return ret;
  }

  newData(params, urlSearchParams) {
    const data = urlSearchParams ? new URLSearchParams() : new FormData();
    if (!util.isEmpty(params)) {
      for (const key in params) {
        const val = params[key];
        // console.debug(`appending form-data.. key:${key}, val:${val}`);
        data.append(key, val);
      }
    }
    return data;
  }

  async request(params, pathParams) {
    if (!util.isEmpty(pathParams)) {
      params.url = this.formatUrl(params.url, pathParams);
    }
    console.info(`>>> Requesting.. ${params.url}`);
    if (!util.isEmpty(this.defaultHeaders)) {
      let _headers = params.headers || {};
      params.headers = {
        ...this.defaultHeaders,
        ..._headers,
      };
    }

    // this._debug('>>>      params:', util.stringify(params));
    try {
      const res = await axios(params);
      console.info('>>> Requested. status:', res.status);
      this._debug(`>>> Status: ${res.status} Params:`, params, 'Response:', res.data);
      // this._debug('>>> Response:', res.data);
      this._debug(`>>> res:`, res);
      return res;
    } catch (e) {
      console.error(`>>> Failed to Request. ${params.url}`, params, e);
      if (e.response && e.response.data) {
        console.error('>>> Error response.data:', e.response.data);
      }
      throw e;
    }
  }

  _debug(...message) {
    if (process.env.VUE_APP_DEBUG_API_CLIENT !== '1') {
      return;
    }
    console.debug(...message);
  }
}
