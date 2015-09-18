import Immutable from 'immutable';

import ActionTypes from '../actionTypes';

const {
  OK,
  ERROR,
  ASYNC_SUCCESS,
  RESET,
} = ActionTypes;

const initialState = Immutable.Map();

export default function(state=initialState, action) {
  const defaultForm = Immutable.fromJS({
    checked: [],
    errors: {}
  });

  let form = defaultForm;

  switch(action.type) {
    case RESET:
      return state.set(action.form, defaultForm);
    case OK:
      form = (state.get(action.form) || defaultForm)
      .update("errors", errors => errors.delete(action.field))
      .update("checked", checked => checked.push(action.field));
      return state.set(action.form, form);
    case ERROR:
      form = (state.get(action.form) || defaultForm)
      .update("errors", errors => errors.set(action.field, action.error))
      .update("checked", checked => checked.push(action.field));
      return state.set(action.form, form);
    case ASYNC_SUCCESS:
      let error = action.meta.resolve(action.payload) ? action.meta.error : false;
      form = state.get(action.meta.form) || defaultForm;
      if (error) {
        form = form
          .update("errors", errors => errors.set(action.meta.field, error))
          .update("checked", checked => checked.push(action.meta.field));
      } else {
        form = form
          .update("errors", errors => errors.delete(action.meta.field))
          .update("checked", checked => checked.push(action.meta.field));
      }
      return state.set(action.meta.form, form);

    default:
      return state;
  }
}




