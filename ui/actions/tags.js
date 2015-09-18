import * as api from '../api';
import ActionTypes from '../actionTypes';

const {
  FETCH_TAGS_PENDING,
  FETCH_TAGS_SUCCESS,
  FETCH_TAGS_FAILURE,
  FILTER_TAGS,
  ORDER_TAGS
} = ActionTypes;

export function getTags() {
  return {
    types: [
      FETCH_TAGS_PENDING,
      FETCH_TAGS_SUCCESS,
      FETCH_TAGS_FAILURE
    ],
    payload: {
      promise: api.getTags()
    }
  }
}

export function filterTags(filter) {
  return {
    type: FILTER_TAGS,
    filter: filter
  };
}

export function orderTags(orderBy) {
  return {
    type: ORDER_TAGS,
    orderBy: orderBy
  };
}


