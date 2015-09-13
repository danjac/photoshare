import Immutable from 'immutable';

import { ActionTypes } from '../constants';

const {
  FETCH_PHOTOS_SUCCESS
} = ActionTypes;


const initialState = Immutable.fromJS({
  currentPage: 1,
  numPages: 0,
  total: 0,
  photos: [],
});

export default function(state=initialState, action) {
  switch(action.type) {
    case FETCH_PHOTOS_SUCCESS:
      return Immutable.fromJS(action.payload);
    default:
      return state;
  }
}




