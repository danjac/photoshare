
import Immutable from 'immutable';

import { ActionTypes } from '../constants';

const {
  GET_PHOTO_DETAIL
} = ActionTypes;


const initialState = Immutable.Map();

export default function(state=initialState, action) {
  switch(action.type) {
    case GET_PHOTO_DETAIL:
      return Immutable.fromJS(action.photo);
    default:
      return state;
  }
}




