
import Immutable from 'immutable';

import ActionTypes from '../actionTypes/photoDetail';

const {
  FETCH_PHOTO_DETAIL_PENDING,
  FETCH_PHOTO_DETAIL_SUCCESS,
  FETCH_PHOTO_DETAIL_FAILURE,
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
    photo: "",
    tags: [],
    perms: {
      vote: false,
      edit: false,
      delete: false
    },
    upVotes: 0,
    downVotes: 0
  },
  isDeleted: false,
  isLoaded: false,
  isNotFound: false,
  isEditingTitle: false,
  isEditingTags: false
});

export default function(state=initialState, action) {
  switch(action.type) {
    case FETCH_PHOTO_DETAIL_PENDING:
      return state.set('isLoaded', false);

    case FETCH_PHOTO_DETAIL_SUCCESS:

      return initialState
        .set("photo", Immutable.fromJS(action.payload))
        .set("isLoaded", true);

    case FETCH_PHOTO_DETAIL_FAILURE:

      return initialState
        .set("isLoaded", true)
        .set("isNotFound", true);

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
      return state.set("isDeleted", true);

    default:
      return state;
  }
}
