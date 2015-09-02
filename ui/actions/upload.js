/* jslint ignore:start */

import * as api from '../api';
import { ActionTypes } from '../constants';

const {
  PHOTO_PREVIEW,
  PHOTO_UPLOAD,
  UPLOAD_RESET,
  UPLOAD_PROGRESS,
  UPLOAD_ERRORS,
  UPLOAD_SUBMITTED
} = ActionTypes;

function validate(title, tags, photo) {
  const errors = new Map();

  if (!title) {
    errors.set("title", "You must provide a title");
  }

  if (!photo) {
    errors.set("photo", "You must provide a photo");
  } else if (!photo.type.match('image.*')) {
    errors.set("photo", "Photo must be an image")
  }

  return errors;
}

export function upload(title, tags, photo) {

  const errors = validate(title, tags, photo);

  if (errors.size > 0) {
    return dispatch => dispatch(formErrors(errors));
  }

  return dispatch => {
    api.upload(title, tags, photo)
    .then(photo => dispatch(uploadDone(photo)));
  }

}

export function formSubmitted() {
    return {
        type: UPLOAD_SUBMITTED
    }
}

export function formErrors(errors){
    return {
        type: UPLOAD_ERRORS,
        errors: errors
    };
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

  let url = null;

  if (file.type.match('image.*')) {
    url = URL.createObjectURL(file);
  }

  return {
    type: PHOTO_PREVIEW,
    url: url
  }
}
