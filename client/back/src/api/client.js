'use strict';

import Const from '../constants/constants.js';
import BaseClient from './baseClient.js';

const ENDPOINT = {
  LIST: {
    url: `${Const.API_PREFIX}`,
    method: 'get'
  }
};

const ApiClient = class ApiClient extends BaseClient {
  constructor() {
    super('ApiClient');
  }

  async list(params) {
    const res = await super.request({
      ...ENDPOINT.LIST,
      params
    });
    return res.data;
  }
};
export default new ApiClient();
