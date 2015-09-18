import Immutable from 'immutable';
import * as api from '../api';
import ActionTypes from '../actionTypes/auth';

const {
  LOGOUT,
  LOGIN_PENDING,
  LOGIN_SUCCESS,
  LOGIN_FAILURE,
  SIGNUP_FIELD_ERROR,
  SIGNUP_FIELD_OK,
  SIGNUP_PENDING,
  SIGNUP_SUCCESS,
  SIGNUP_FAILURE,
  SIGNUP_CHECK_ASYNC_PENDING,
  SIGNUP_CHECK_ASYNC_SUCCESS,
  SIGNUP_CHECK_ASYNC_FAILURE,
  FETCH_USER_PENDING,
  FETCH_USER_SUCCESS,
  FETCH_USER_FAILURE
} = ActionTypes;



export function checkName(name) {
  let error = null;
  if (!name || name.length < 6) {
    error = "Name must be at least 6 characters"
  } else if (name.length > 30) {
    error = "Name must be max 30 characters"
  }
  if (error) {
    return {
      type: SIGNUP_FIELD_ERROR,
      error: error,
      field: "name"
    }
  }
  return {
    type: SIGNUP_FIELD_OK,
    field: "name"
  }
}

export function checkPassword(password) {
  let error = null;
  if (!password || password.length < 6) {
    error = "Password must be at least 6 characters"
  }
  if (error) {
    return {
      type: SIGNUP_FIELD_ERROR,
      error: error,
      field: "password"
    }
  }
  return {
    type: SIGNUP_FIELD_OK,
    field: "password"
  }
}


export function checkEmail(email) {

  if (!email || email.indexOf("@") === -1) {
    return {
      type: SIGNUP_FIELD_ERROR,
      field: "email",
      error: "You must provide a valid email address"
    };
  }

  return {
    types: [
      SIGNUP_CHECK_ASYNC_PENDING,
      SIGNUP_CHECK_ASYNC_SUCCESS,
      SIGNUP_CHECK_ASYNC_FAILURE,
    ],
    payload: {
      promise: api.emailExists(email),
    },
    meta: {
      resolve: result => result.exists,
      error: "Email already exists",
      field: "email"
    }
  }
}

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
  };
}

export function logout() {
  api.logout();
  return {
    type: LOGOUT
  };
}

export function signup(name, email, password) {
  return {
    types: [
      SIGNUP_PENDING,
      SIGNUP_SUCCESS,
      SIGNUP_FAILURE
    ],
    payload: {
      promise: api.signup(name, email, password)
    }
  };

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
  };

}
