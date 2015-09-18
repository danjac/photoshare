import Immutable from 'immutable';

import ActionTypes  from '../actionTypes';

const  {
  RECOVER_PASSWORD_RESET,
  RECOVER_PASSWORD_PENDING,
  RECOVER_PASSWORD_SUCCESS,
  RECOVER_PASSWORD_FAILURE
} = ActionTypes;


const initialState = Immutable.fromJS({
  formSubmitted: false,
  isSuccess: false,
  isError: false
});


export default function(state=initialState, action) {
  switch(action.type) {
    case RECOVER_PASSWORD_RESET:

      return initialState;

    case RECOVER_PASSWORD_PENDING:

      return state.set("formSubmitted", true);

    case RECOVER_PASSWORD_SUCCESS:

      return state
        .set("formSubmitted", false)
        .set("isError", false)
        .set("isSuccess", true);

    case RECOVER_PASSWORD_FAILURE:

      return state
        .set("formSubmitted", false)
        .set("isSuccess", false)
        .set("isError", action.error);

    default:
      return state;
  }
}
