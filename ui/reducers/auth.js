import Immutable from 'immutable';

import ActionTypes from '../actionTypes/auth';

const {
  LOGIN_SUCCESS,
  SIGNUP_SUCCESS,
  FETCH_USER_SUCCESS,
  LOGOUT
} = ActionTypes;


const initialState = Immutable.fromJS({

  id: null,
  name: null,
  email: null,
  isAdmin: false,
  loggedIn: false

});

export default function(state=initialState, action) {
  switch(action.type) {
    case LOGIN_SUCCESS:
      return Immutable.fromJS(action.payload);
    case SIGNUP_SUCCESS:
      return Immutable.fromJS(action.payload);
    case FETCH_USER_SUCCESS:
      return Immutable.fromJS(action.payload || {});
    case LOGOUT:
      return initialState;
    default:
      return state;
  }
}
