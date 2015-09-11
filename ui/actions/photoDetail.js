
import * as api from '../api';
import { ActionTypes } from '../constants';

const {
  GET_PHOTO_DETAIL,
  DELETE_PHOTO,
  EDIT_PHOTO_TITLE,
  EDIT_PHOTO_TAGS,
  UPDATE_PHOTO_TITLE,
  UPDATE_PHOTO_TAGS
} = ActionTypes;

export function getPhotoDetail(id) {
  return dispatch => {
    api.getPhotoDetail(id)
    .then(photo => {
      dispatch(getPhotoDetailDone(photo));
    });
  }
}

export function toggleEditTitle() {
  return {
    type: EDIT_PHOTO_TITLE
  }
}

export function toggleEditTags() {
  return {
    type: EDIT_PHOTO_TAGS
  }
}

export function updateTitle(id, title) {
  api.updatePhotoTitle(id, title);
  return {
    type: UPDATE_PHOTO_TITLE,
    title: title
  }
}

export function updateTags(id, tags) {
  api.updatePhotoTags(id, tags);
  return {
    type: UPDATE_PHOTO_TAGS,
    tags: tags
  }
}


export function deletePhoto(id) {
  api.deletePhoto(id);
  return {
    type: DELETE_PHOTO
  }
}

export function getPhotoDetailDone(photo) {
  return {
    type: GET_PHOTO_DETAIL,
    photo: photo
  };
}
