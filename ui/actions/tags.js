import * as api from '../api';
import { ActionTypes } from '../constants';

const {
  GET_TAGS,
  FILTER_TAGS,
  ORDER_TAGS
} = ActionTypes;

export function getTags() {
  return dispatch => {
    api.getTags()
    .then(tags => {
      dispatch(getTagsDone(tags));
    });
  };
}

export function getTagsDone(tags) {
  return {
    type: GET_TAGS,
    tags: tags
  }
}

export function filterTags(filterStr) {
  return {
    type: FILTER_TAGS,
    filter: filterStr
  };
}

export function orderTags(orderBy) {
  return {
    type: ORDER_TAGS,
    orderBy: orderBy
  };
}


