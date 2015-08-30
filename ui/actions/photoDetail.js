
import * as api from '../api';
import { ActionTypes } from '../constants';

const {
  GET_PHOTO_DETAIL 
} = ActionTypes;

export function getPhotoDetail(id) {
  return dispatch => {
    api.getPhotoDetail(id)
    .then(photo => {
      dispatch(getPhotoDetailDone(photo));
    });
  }
}

export function getPhotoDetailDone(photo) {
  return {
    type: GET_PHOTO_DETAIL,
    photo: photo
  };
}
