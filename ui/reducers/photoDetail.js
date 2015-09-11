
import Immutable from 'immutable';

import { ActionTypes } from '../constants';

const {
  GET_PHOTO_DETAIL,
  EDIT_PHOTO_TITLE,
  EDIT_PHOTO_TAGS,
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
  isEditingTitle: false,
  isEditingTags: false
});

export default function(state=initialState, action) {
  switch(action.type) {
    case GET_PHOTO_DETAIL:
      return state.set("photo", Immutable.Map(action.photo));

    case VOTE_UP:
    case VOTE_DOWN:
      return state.setIn(["photo", "perms", "canVote"], false);

    case EDIT_PHOTO_TITLE:
      return state.set("isEditingTitle", !state.get("isEditingTitle"));

    case EDIT_PHOTO_TAGS:
      return state.set("isEditingTags", !state.get("isEditingTags"));

    case UPDATE_PHOTO_TITLE:

      return state
        .set("isEditingTitle", false)
        .setIn(["photo", "title"], action.title);

    case UPDATE_PHOTO_TAGS:
      return state
        .set("isEditingTags", false)
        .setIn(["photo", "tags"], action.tags);

    case DELETE_PHOTO:
      return initialState;

    default:
      return state;
  }
}
