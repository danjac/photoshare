import Immutable from 'immutable';

import ActionTypes from '../actionTypes/auth';

const {
  LOGIN_SUCCESS,
  SIGNUP_SUCCESS,
  SIGNUP_PENDING,
  SIGNUP_FIELD_ERROR,
  SIGNUP_FIELD_OK,
  SIGNUP_CHECK_ASYNC_SUCCESS,
  FETCH_USER_SUCCESS,
  LOGOUT
} = ActionTypes;


const initialState = Immutable.fromJS({

  id: null,
  name: null,
  email: null,
  isAdmin: false,
  loggedIn: false,

  signupErrors: new Immutable.Map(),
  signupPrechecks: new Immutable.List(),
  signupFormSubmitted: false

});

export default function(state=initialState, action) {
  switch(action.type) {

    case LOGIN_SUCCESS:
    case SIGNUP_SUCCESS:
    case FETCH_USER_SUCCESS:

      return state.merge(action.payload);

    case SIGNUP_FIELD_OK:
        return state
        .update("signupPrechecks", checks => checks.push(action.field))
        .update("signupErrors", errors => errors.delete(action.field));

    case SIGNUP_FIELD_ERROR:
        return state
        .update("signupPrechecks", checks => checks.push(action.field))
        .update("signupErrors", errors => errors.set(action.field, action.error));

    case SIGNUP_CHECK_ASYNC_SUCCESS:

      let error = action.meta.resolve(action.payload) ? action.meta.error: null;

      if (error) {
        return state
        .update("signupPrechecks", checks => checks.push(action.meta.field))
        .update("signupErrors", errors => errors.set(action.meta.field, error));
      }
      return state
      .update("signupPrechecks", checks => checks.push(action.meta.field))
      .update("signupErrors", errors => errors.delete(action.meta.field));

    case SIGNUP_PENDING:
      return state.set("signupFormSubmitted", true);

    case LOGOUT:
      return initialState;
    default:
      return state;
  }
}
