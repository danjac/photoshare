import Immutable from 'immutable';

import ActionTypes from '../actionTypes';

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
  loggedIn: false,
});

export default function(state=initialState, action) {
  switch(action.type) {

    case LOGIN_SUCCESS:
    case SIGNUP_SUCCESS:
    case FETCH_USER_SUCCESS:

      return state.merge(action.payload);

      case LOGOUT:
      return initialState;
    default:
      return state;
  }
}
