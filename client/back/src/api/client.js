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

  async list(params) {
    const res = await super.request({
      ...ENDPOINT.LIST,
      params,
    });
    return res.data;
  }

  async download(params) {
    const res = await super.request({
      ...ENDPOINT.DL_REQUEST,
      data: params,
    });
    return res.data;
  }
};
export default new ApiClient();
