import Immutable from 'immutable';

import { ActionTypes }  from '../constants';

const {
  LOGIN_SUCCESS,
  LOGIN_FAILURE,
  GET_USER,
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
      console.log("login success:", action.user);
      return Immutable.fromJS(action.user);
    case GET_USER:
      return Immutable.fromJS(action.user || {});
    case LOGOUT:
      return initialState;
    default:
      return state;
  }
}
