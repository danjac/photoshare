
import * as api from '../api';
import { ActionTypes } from '../constants';

const {
  GET_PHOTO_DETAIL,
  DELETE_PHOTO
} = ActionTypes;

export function getPhotoDetail(id) {
  return dispatch => {
    api.getPhotoDetail(id)
    .then(photo => {
      dispatch(getPhotoDetailDone(photo));
    });
  }
}

export function deletePhoto(photo) {
  api.deletePhoto(photo.id);
  return {
    type: DELETE_PHOTO,
    photo: photo
  }
}

export function getPhotoDetailDone(photo) {
  return {
    type: GET_PHOTO_DETAIL,
    photo: photo
  };
}
