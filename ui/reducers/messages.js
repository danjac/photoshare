import Immutable from 'immutable';

import { ActionTypes }  from '../constants';

const { NEW_MESSAGE, DELETE_MESSAGE } = ActionTypes;
const initialState = Immutable.List();

export default function(state=initialState, action) {
  switch(action.type) {
    case NEW_MESSAGE:
      return state.unshift(Immutable.Map(action.message));
    case DELETE_MESSAGE:
      return state.delete(action.key);
    default:
      return state; 
  }
}

