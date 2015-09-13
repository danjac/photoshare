import * as api from '../api';
import { ActionTypes } from '../constants';

const {
  FETCH_PHOTOS_PENDING,
  FETCH_PHOTOS_SUCCESS,
  FETCH_PHOTOS_FAILURE
} = ActionTypes;


function fetchPhotos(promise) {
  return {
    types: [
      FETCH_PHOTOS_PENDING,
      FETCH_PHOTOS_SUCCESS,
      FETCH_PHOTOS_FAILURE
    ],
    payload: {
      promise: promise
    }
  }
}

export function getPhotos(page, orderBy) {
  return fetchPhotos(api.getPhotos(page, orderBy));
}

export function searchPhotos(page, query) {
  return fetchPhotos(api.searchPhotos(page, query));
}
