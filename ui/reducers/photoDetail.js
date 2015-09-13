
import Immutable from 'immutable';

import { ActionTypes } from '../constants';

const {
  FETCH_PHOTO_DETAIL_SUCCESS,
  UPDATE_PHOTO_TAGS_PENDING,
  UPDATE_PHOTO_TITLE_PENDING,
  DELETE_PHOTO_PENDING,
  TOGGLE_PHOTO_TITLE_EDIT,
  TOGGLE_PHOTO_TAGS_EDIT,
  VOTE_UP_PHOTO_PENDING,
  VOTE_DOWN_PHOTO_PENDING
} = ActionTypes;


const initialState = Immutable.Map({
  photo: {
    title: "",
    tags: [],
    perms: {
      vote: false,
      edit: false,
      delete: false
    },
    upVotes: 0,
    downVotes: 0
  },
  isLoaded: false,
  isEditingTitle: false,
  isEditingTags: false
});

export default function(state=initialState, action) {
  switch(action.type) {
    case FETCH_PHOTO_DETAIL_SUCCESS:

      return state
        .set("photo", Immutable.fromJS(action.payload))
        .set("isLoaded", true);

    case TOGGLE_PHOTO_TITLE_EDIT:
      return state
        .set("isEditingTitle",
            !state.get("isEditingTitle"));

    case TOGGLE_PHOTO_TAGS_EDIT:
      return state
        .set("isEditingTags",
            !state.get("isEditingTags"));

    case VOTE_UP_PHOTO_PENDING:
      return state
        .setIn(["photo", "perms", "vote"], false)
        .updateIn(["photo", "upVotes"], value => value + 1);

    case VOTE_DOWN_PHOTO_PENDING:
      return state
        .setIn(["photo", "perms", "vote"], false)
        .updateIn(["photo", "downVotes"], value => value - 1);

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
