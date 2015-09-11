
import Immutable from 'immutable';

import { ActionTypes } from '../constants';

const {
  GET_PHOTO_DETAIL,
  EDIT_PHOTO,
  UPDATE_PHOTO_TAGS,
  UPDATE_PHOTO_TITLE,
  DELETE_PHOTO,
  VOTE_UP,
  VOTE_DOWN
} = ActionTypes;


const initialState = Immutable.Map({
  photo: {
    title: "",
    tags: [],
    perms: {}
  },
  isEditing: false
});

export default function(state=initialState, action) {
  switch(action.type) {
    case GET_PHOTO_DETAIL:
      return state.set("photo", action.photo);
    case VOTE_UP:
    case VOTE_DOWN:
      return state.setIn(["photo", "perms", "canVote"], false);
    case EDIT_PHOTO:
      return state.set("isEditing", !state.get("isEditing"));
    case UPDATE_PHOTO_TITLE:
      return state.setIn(["photo", "title"], action.title).set("isEditing", false);
    case UPDATE_PHOTO_TAGS:
      return state.setIn(["photo", "tags"], action.tags).set("isEditing", false);
    case DELETE_PHOTO:
      return initialState;
    default:
      return state;
  }
}




