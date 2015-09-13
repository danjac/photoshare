
import Immutable from 'immutable';

import { ActionTypes } from '../constants';

const {
  FETCH_PHOTO_DETAIL_SUCCESS,
  UPDATE_PHOTO_TAGS_PENDING,
  UPDATE_PHOTO_TITLE_PENDING,
  DELETE_PHOTO_PENDING,
  TOGGLE_PHOTO_TITLE_EDIT,
  TOGGLE_PHOTO_TAGS_EDIT,
  VOTE_UP_PENDING,
  VOTE_DOWN_PENDING
} = ActionTypes;


const initialState = Immutable.Map({
  photo: {
    title: "",
    tags: [],
    perms: {}
  },
  isLoaded: false,
  isEditingTitle: false,
  isEditingTags: false
});

export default function(state=initialState, action) {
  switch(action.type) {
    case FETCH_PHOTO_DETAIL_SUCCESS:

      return state
        .set("photo", Immutable.Map(action.payload))
        .set("isLoaded", true);

    case TOGGLE_PHOTO_TITLE_EDIT:
      return state
        .set("isEditingTitle",
            !state.get("isEditingTitle"));

    case TOGGLE_PHOTO_TAGS_EDIT:
      return state
        .set("isEditingTags",
            !state.get("isEditingTags"));

    case VOTE_UP_PENDING:
    case VOTE_DOWN_PENDING:
      return state.setIn(["photo", "perms", "canVote"], false);

    case UPDATE_PHOTO_TITLE_PENDING:
      return state
        .set("isEditingTitle", false)
        .setIn(["photo", "title"], action.payload);

    case UPDATE_PHOTO_TAGS_PENDING:
      return state
        .set("isEditingTags", false)
        .setIn(["photo", "tags"], action.payload);

    case DELETE_PHOTO_PENDING:
      return initialState;

    default:
      return state;
  }
}
