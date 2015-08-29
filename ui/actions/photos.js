import * as api from '../api';
import { ActionTypes } from '../constants';

const {
  GET_PHOTOS 
} = ActionTypes;

export function getPhotos(page, orderBy) {
  return dispatch => {
    api.getPhotos(page, orderBy)
    .then(photos => {
      dispatch(getPhotosDone(photos));
    });
  }
}

export function getPhotosDone(photos) {
  return {
    type: GET_PHOTOS,
    photos: photos
  };
}
