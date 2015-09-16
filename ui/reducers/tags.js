import Immutable from 'immutable';

import ActionTypes from '../actionTypes/tags';

const {
  FETCH_TAGS_SUCCESS,
  FILTER_TAGS,
  ORDER_TAGS
} = ActionTypes;

const initialState = Immutable.fromJS({
  source: [],
  filter: "",
  orderBy: "numPhotos"
});

export default function(state=initialState, action) {
  switch(action.type) {
    case FETCH_TAGS_SUCCESS:
      return initialState.set("source", Immutable.fromJS(action.payload));
    case FILTER_TAGS:
      return state.set("filter", action.filter);
    case ORDER_TAGS:
      return state.set("orderBy", action.orderBy);
    default:
      return state;
  }
}
