import Immutable from 'immutable';

import ActionTypes from '../actionTypes';

const {
  CHANGE_PASSWORD_RESET,
  CHANGE_PASSWORD_FORM_INVALID,
  CHANGE_PASSWORD_PENDING,
  CHANGE_PASSWORD_SUCCESS,
  CHANGE_PASSWORD_FAILURE
} = ActionTypes;

const initialState = Immutable.fromJS({
  formSubmitted: false,
  isWaiting: false,
  errors: new Map(),
  isSuccess: false,
  isServerError: false
});


export default function(state=initialState, action) {

  switch(action.type) {
    case CHANGE_PASSWORD_RESET:
      return initialState;

    case CHANGE_PASSWORD_FORM_INVALID:
      return state
        .set("errors", action.errors)
        .set("formSubmitted", true);

    case CHANGE_PASSWORD_PENDING:
      return state
        .set("isWaiting", true)
        .set("errors", new Map());

    case CHANGE_PASSWORD_SUCCESS:
       return state
         .set("isWaiting", false)
         .set("formSubmitted", false)
         .set("isSuccess", true);

    case CHANGE_PASSWORD_FAILURE:
      return state
        .set("isWaiting", false)
        .set("formSubmitted", false)
        .set("isServerError", true);

    default:
      return state;
  }

}

