import * as api from '../api';
import { ActionTypes } from '../constants';

const {
  RECOVER_PASSWORD_RESET,
  RECOVER_PASSWORD_PENDING,
  RECOVER_PASSWORD_SUCCESS,
  RECOVER_PASSWORD_FAILURE
} = ActionTypes;


export function resetForm() {
  return {
    type: RECOVER_PASSWORD_RESET
  }
}

export function recoverPassword(email) {
  return {
    types: [
      RECOVER_PASSWORD_PENDING,
      RECOVER_PASSWORD_SUCCESS,
      RECOVER_PASSWORD_FAILURE
    ],
    payload: {
      promise: api.recoverPassword(email)
    }
  }
}
