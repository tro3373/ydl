'use strict';

const CommonConstants = {
  API_PREFIX: process.env.VUE_APP_API_PREFIX,
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
