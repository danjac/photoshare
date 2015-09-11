import Immutable from 'immutable';

import { ActionTypes } from '../constants';

const {
  GET_TAGS,
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
    case GET_TAGS:
      return state.set("source", Immutable.List(action.tags));
    case FILTER_TAGS:
      return state.set("filter", action.filter);
    case ORDER_TAGS:
      return state.set("orderBy", action.orderBy);
    default:
      return state;
  }
}
