import * as api from '../api';
import { ActionTypes } from '../constants';

const {
  LOGIN_SUCCESS,
  LOGIN_FAILURE,
  GET_USER,
  LOGOUT
} = ActionTypes;

export function getUser() {
  return dispatch => {
    api.getUser()
    .then(user => dispatch(getUserComplete(user)));
  }
}

export function getUserComplete(user) {
  return {
    type: GET_USER,
    user: user || {}
  }
}

export function logout() {
  api.logout();
  return {
    type: LOGOUT
  }
}

export function login(identifier, password) {
  return dispatch => {
    api.login(identifier, password)
    .then(user => {
      dispatch(loginSuccess(user));
    })
    .catch(err => {
      dispatch(loginFailure(err));
    });
  };
}

export function loginSuccess(user) {
  return {
    type: LOGIN_SUCCESS,
    user: user
  }
}

export function loginFailure(err) {
  return {
    type: LOGIN_FAILURE,
    error: err
  }
}
