import Immutable from 'immutable';

import { ActionTypes } from '../constants';

const {
  GET_PHOTOS
} = ActionTypes;


const initialState = Immutable.List();

export default function(state=initialState, action) {
  switch(action.type) {
    case GET_PHOTOS:
      return Immutable.fromJS(action.photos);
    default:
      return state;
  }
}




