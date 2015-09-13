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
  progress: 0,
  formSubmitted: false,
  errors: new Map()
});

export default function(state=initialState, action) {
  switch(action.type) {
    case UPLOAD_PENDING:
      return state.set('formSubmitted', true);
    case UPLOAD_FORM_INVALID:
      return state.merge({'errors': action.errors, 'formSubmitted': true});
    case PHOTO_PREVIEW:
      return state.set('previewURL', action.url);
    case UPLOAD_SUCCESS:
      return state.set('uploadedPhoto', action.photo);
    default:
      return state;
  }
}
