'use strict';

import { v4 as uuidv4 } from 'uuid';

const core = {};

core.isEmpty = (arg) => {
  if (arg == null) return true;
  if (arg === void 0) return true;
  switch (typeof arg) {
    case 'object':
      if (Array.isArray(arg)) {
        // When object is array:
        return arg.length === 0;
      } else {
        // When object is not array:
        if (Object.keys(arg).length > 0 || Object.getOwnPropertySymbols(arg).length > 0) {
          return false;
        } else if (arg.valueOf().length !== undefined) {
          return arg.valueOf().length === 0;
        } else if (typeof arg.valueOf() !== 'object') {
          return this.isEmpty(arg.valueOf());
        } else {
          return true;
        }
      }
    default:
      break;
  }
  let tmp = '' + arg;
  return 0 === tmp.length;
};

core.uuid = () => {
  return uuidv4();
};

core.uniq = (list) => {
  return [...new Set(list)];
};

export default core;
