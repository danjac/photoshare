import * as api from '../api';
import { ActionTypes } from '../constants';

const {
  PHOTO_PREVIEW,
  PHOTO_UPLOAD,
  UPLOAD_RESET,
  UPLOAD_PROGRESS
} = ActionTypes;


export function upload(title, tags, photo) {

  return dispatch => {
    api.upload(title, tags, photo)
    .then(photo => dispatch(uploadDone(photo)));
  }

}

export function uploadDone(photo) {
  return {
    type: PHOTO_UPLOAD,
    photo: photo
  };
}

export function reset() {
  return {
    type: UPLOAD_RESET
  }
}

export function progress() {
  return {
    type: UPLOAD_PROGRESS
  }
}

export function previewPhoto(file) {

  return dispatch => {
    const reader = new window.FileReader();
    reader.onload = ((f) => {
      return (event) => {
        dispatch(previewPhotoDone(event.target.result));
      };
    })(file);
    reader.readAsDataURL(file);
  }
}

export function previewPhotoDone(url) {
  return {
    type: PHOTO_PREVIEW,
    url: url
  };
}
