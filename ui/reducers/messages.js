import Immutable from 'immutable';

import { ActionTypes }  from '../constants';

const {
  DELETE_MESSAGE,
  LOGOUT,
  LOGIN_SUCCESS,
  DELETE_PHOTO
} = ActionTypes;
const initialState = Immutable.List();

function newMessage(state, msg, level) {
  return state.unshift(new Immutable.Map({
    level: level,
    msg: msg
  }));
}

export default function(state=initialState, action) {
  switch(action.type) {
    case DELETE_MESSAGE:
      return state.delete(action.key);
    case LOGIN_SUCCESS:
      return newMessage(state, `Welcome back, ${action.user.name}`, "success");
    case LOGOUT:
      return newMessage(state, "Bye for now", "info");
    case DELETE_PHOTO:
      return newMessage(state, "Your photo has been deleted", "info");
    default:
      return state;
  }
}

