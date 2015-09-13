import Immutable from 'immutable';

import { ActionTypes }  from '../constants';

const {
  DELETE_MESSAGE,
  LOGOUT,
  LOGIN_SUCCESS,
  SIGNUP_SUCCESS,
  SIGNUP_FAILURE,
  DELETE_PHOTO_SUCCESS
} = ActionTypes;
const initialState = Immutable.List();

function newMessage(state, msg, level) {
  return state.unshift(new Immutable.Map({
    level: level,
    msg: msg
  }));
}

const MessageLevel = {
  INFO: "info",
  SUCCESS: "success",
  WARNING: "warning",
  DANGER: "danger"
}

export default function(state=initialState, action) {
  switch(action.type) {
    case DELETE_MESSAGE:
      return state.delete(action.key);

    case LOGIN_SUCCESS:

      return newMessage(
          state,
          `Welcome back, ${action.payload.name}`,
          MessageLevel.SUCCESS);

    case LOGOUT:

      return newMessage(
          state,
          "Bye for now",
          MessageLevel.INFO
          );

    case DELETE_PHOTO_SUCCESS:

      return newMessage(
          state,
          "Your photo has been deleted",
          MessageLevel.INFO
          );

    case SIGNUP_SUCCESS:

      return newMessage(
          state,
          `Welcome ${action.payload.name}`,
          MessageLevel.SUCCESS
          );

    case SIGNUP_FAILURE:

      return newMessage(
          state,
          `Signup failed: ${action.payload.message}`,
          MessageLevel.WARNING
          );

    default:
      return state;
  }
}

