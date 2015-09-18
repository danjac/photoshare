import * as api from '../api';
import ActionTypes from '../actionTypes';

const {
  TOGGLE_PHOTO_TITLE_EDIT,
  TOGGLE_PHOTO_TAGS_EDIT,
  FETCH_PHOTO_DETAIL_PENDING,
  FETCH_PHOTO_DETAIL_SUCCESS,
  FETCH_PHOTO_DETAIL_FAILURE,
  DELETE_PHOTO_PENDING,
  DELETE_PHOTO_SUCCESS,
  DELETE_PHOTO_FAILURE,
  UPDATE_PHOTO_TITLE_PENDING,
  UPDATE_PHOTO_TITLE_SUCCESS,
  UPDATE_PHOTO_TITLE_FAILURE,
  UPDATE_PHOTO_TAGS_PENDING,
  UPDATE_PHOTO_TAGS_SUCCESS,
  UPDATE_PHOTO_TAGS_FAILURE,
  VOTE_UP_PHOTO_PENDING,
  VOTE_UP_PHOTO_SUCCESS,
  VOTE_UP_PHOTO_FAILURE,
  VOTE_DOWN_PHOTO_PENDING,
  VOTE_DOWN_PHOTO_SUCCESS,
  VOTE_DOWN_PHOTO_FAILURE
} = ActionTypes;

export function getPhotoDetail(id) {
  return {
    types: [
      FETCH_PHOTO_DETAIL_PENDING,
      FETCH_PHOTO_DETAIL_SUCCESS,
      FETCH_PHOTO_DETAIL_FAILURE
    ],
    payload: {
      promise: api.getPhotoDetail(id)
    }
  };
}

export function voteUp(id) {
  return {
    types: [
      VOTE_UP_PHOTO_PENDING,
      VOTE_UP_PHOTO_SUCCESS,
      VOTE_UP_PHOTO_FAILURE
    ],
    payload: {
      promise: api.votePhotoUp(id)
    }
  }
}

export function voteDown(id) {
  return {
    types: [
      VOTE_DOWN_PHOTO_PENDING,
      VOTE_DOWN_PHOTO_SUCCESS,
      VOTE_DOWN_PHOTO_FAILURE
    ],
    payload: {
      promise: api.votePhotoDown(id)
    }
  }
}


export function toggleEditTitle() {
  return {
    type: TOGGLE_PHOTO_TITLE_EDIT
  }
}

export function toggleEditTags() {
  return {
    type: TOGGLE_PHOTO_TAGS_EDIT
  }
}

export function updateTitle(id, title) {
  return {
    types: [
      UPDATE_PHOTO_TITLE_PENDING,
      UPDATE_PHOTO_TITLE_SUCCESS,
      UPDATE_PHOTO_TITLE_FAILURE
    ],
    payload: {
      promise: api.updatePhotoTitle(id, title),
      data: title
    }
  }
}

export function updateTags(id, tags) {
  return {
    types: [
      UPDATE_PHOTO_TAGS_PENDING,
      UPDATE_PHOTO_TAGS_SUCCESS,
      UPDATE_PHOTO_TAGS_FAILURE
    ],
    payload: {
      promise: api.updatePhotoTags(id, tags),
      data: tags
    }
  }
}


export function deletePhoto(id) {
  return {
    types: [
      DELETE_PHOTO_PENDING,
      DELETE_PHOTO_SUCCESS,
      DELETE_PHOTO_FAILURE
    ],
    payload: {
      promise: api.deletePhoto(id)
    }
  }
}
