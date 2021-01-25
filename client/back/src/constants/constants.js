'use strict';

const CommonConstants = {};
const EnvConstants = {
  dev: {
    API_PREFIX: 'http://192.168.33.10'
  },
  prd: {
    API_PREFIX:
      'https://5hsnc6y80l.execute-api.ap-northeast-1.amazonaws.com/prd'
  }
};

const genConstants = () => {
  const env = process.env.VUE_APP_STAGE || 'dev';
  let res = {
    STAGE: env,
    ...CommonConstants,
    ...EnvConstants[env]
  };
  return res;
};

const Const = Object.freeze(genConstants());
export default Const;
