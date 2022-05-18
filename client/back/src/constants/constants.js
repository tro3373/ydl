'use strict';

const CommonConstants = {};
const EnvConstants = {
  dev: {
    API_PREFIX: 'http://localhost/api',
  },
  prd: {
    API_PREFIX: 'https://ydl.chillfixx.work/api',
  },
};

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
