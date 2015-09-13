/* jslint ignore:start */
import Immutable from 'immutable';

import { ActionTypes } from '../constants';

const {
  PHOTO_PREVIEW,
  UPLOAD_FORM_INVALID,
  UPLOAD_SUCCESS,
  UPLOAD_PENDING
} = ActionTypes;


const initialState = Immutable.fromJS({
  previewURL: null,
  uploadedPhoto: null,
  formSubmitted: false,
  errors: new Map()
});

export default function(state=initialState, action) {
  switch(action.type) {
    case PHOTO_PREVIEW:
      return state.set('previewURL', action.url);
    case UPLOAD_PENDING:
      return state.merge({
          errors: new Map(),
          formSubmitted: true,
          previewURL: null
      });
    case UPLOAD_FORM_INVALID:
      return state.merge({
          errors: action.errors,
          formSubmitted: true
      });
    case UPLOAD_SUCCESS:
      return state.set('uploadedPhoto', action.payload);
    default:
      return state;
  }
}
