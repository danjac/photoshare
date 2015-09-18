import * as api from '../api';
import ActionTypes from '../actionTypes';

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

export function getPhotosForOwner(ownerID, page) {
  return fetchPhotos(api.getPhotosForOwner(ownerID, page));
}

export function searchPhotos(page, query) {
  return fetchPhotos(api.searchPhotos(page, query));
}
