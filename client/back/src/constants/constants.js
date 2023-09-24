'use strict';

const CommonConstants = {
  APP_VERSION: '0.0.2',
  API_PREFIX: process.env.VUE_APP_API_PREFIX,
  LOCAL_STRAGE_KEY: {
    CACHE: 'local_strage_key_cache',
    UUID: 'local_strage_key_uuid',
    VISITED: 'local_strage_key_visited',
  },
};
const EnvConstants = {};

const genConstants = () => {
  const env = process.env.VUE_APP_STAGE || 'dev';
  let res = {
    STAGE: env,
    ...CommonConstants,
    ...EnvConstants[env],
  };
  return res;
};

const Const = Object.freeze(genConstants());
export default Const;
