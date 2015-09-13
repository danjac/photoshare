import * as api from '../api';
import { ActionTypes } from '../constants';

const {
  LOGOUT,
  LOGIN_PENDING,
  LOGIN_SUCCESS,
  LOGIN_FAILURE,
  FETCH_USER_PENDING,
  FETCH_USER_SUCCESS,
  FETCH_USER_FAILURE
} = ActionTypes;

export function getUser() {
  return {
    types: [
      FETCH_USER_PENDING,
      FETCH_USER_SUCCESS,
      FETCH_USER_FAILURE
    ],
    payload: {
      promise: api.getUser()
    }
  }
}

export function logout() {
  api.logout();
  return {
    type: LOGOUT
  }
}

export function login(identifier, password) {
  return {
    types: [
      LOGIN_PENDING,
      LOGIN_SUCCESS,
      LOGIN_FAILURE
    ],
    payload: {
      promise: api.login(identifier, password)
    }
  }

}
