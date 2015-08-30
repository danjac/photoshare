import * as api from '../api';
import { ActionTypes } from '../constants';

const {
  LOGIN_SUCCESS,
  LOGIN_FAILURE,
  GET_USER
} = ActionTypes;

export function getUser() {
  console.log("getUser")
  return dispatch => {
    api.getUser()
    .then(user => dispatch(getUserComplete(user)));
  }
}

export function getUserComplete(user) {
  console.log("getUserComplete", user);
  return {
    type: GET_USER,
    user: user || {}
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
