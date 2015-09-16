/* jslint ignore:start */

import * as api from '../api';
import { ActionTypes } from '../constants';

const {
  PHOTO_PREVIEW,
  UPLOAD_FORM_INVALID,
  UPLOAD_RESET,
  UPLOAD_PENDING,
  UPLOAD_SUCCESS,
  UPLOAD_FAILURE
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

export function resetForm() {
  return { type: UPLOAD_RESET };
}

export function upload(title, tags, photo) {

  const errors = validate(title, tags, photo);

  if (errors.size > 0) {
    return dispatch => dispatch({
      type: UPLOAD_FORM_INVALID,
      errors: errors
    });
  }

  return {
    types: [
      UPLOAD_PENDING,
      UPLOAD_SUCCESS,
      UPLOAD_FAILURE
    ],
    payload: {
      promise: api.upload(title, tags, photo)
    }
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
