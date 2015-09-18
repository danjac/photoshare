import Immutable from 'immutable';

import ActionTypes from '../actionTypes';

const {
  FETCH_PHOTOS_PENDING,
  FETCH_PHOTOS_SUCCESS
} = ActionTypes;


const initialState = Immutable.fromJS({
  currentPage: 1,
  numPages: 0,
  total: 0,
  photos: [],
  isLoaded: false
});

export default function(state=initialState, action) {
  switch(action.type) {
    case FETCH_PHOTOS_PENDING:
      return state.set('isLoaded', false);
    case FETCH_PHOTOS_SUCCESS:
      return Immutable.fromJS(action.payload).set('isLoaded', true);
    default:
      return state;
  }
}




