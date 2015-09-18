/* jslint ignore:start */

import * as api from '../api';
import ActionTypes from '../actionTypes';

const {
  UPLOAD_PENDING,
  UPLOAD_SUCCESS,
  UPLOAD_FAILURE,
  OK,
  ERROR,
  RESET
} = ActionTypes;


export function checkTitle(title) {
  let error;

  if (!title) {
    error = "You must provide a title";
  } else if (title.length > 100) {
    error = "Title is too long(100 chars max)";
  }
  if (error) {
    return {
      type: ERROR,
      form: "upload",
      field: "title",
      error: error
    }
  }
  return {
    type: OK,
    form: "upload",
    field: "title"
  }
}

export function checkPhoto(photo) {
  let error;

  if (!photo) {
    error = "You must provide a photo";
  } else if (!photo.type.match('image.*')) {
    error = "Photo must be an image";
  }

  if (error) {
    return {
      type: ERROR,
      form: "upload",
      field: "photo",
      previewURL: null,
      error: error
    }
  };

  return {
    type: OK,
    form: "upload",
    previewURL: URL.createObjectURL(photo),
    field: "photo"
  };
}


export function resetForm() {
  return { type: RESET, form: "upload" };
}

export function upload(title, tags, photo) {

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
