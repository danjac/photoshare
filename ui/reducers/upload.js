
import Immutable from 'immutable';

import { ActionTypes } from '../constants';

const {
  PHOTO_PREVIEW,
  PHOTO_UPLOAD,
  UPLOAD_RESET
} = ActionTypes;


const initialState = Immutable.fromJS({
  previewURL: null,
  uploadedPhoto: null
});

export default function(state=initialState, action) {
  switch(action.type) {
    case PHOTO_PREVIEW:
      return state.set('previewURL', action.url);
    case PHOTO_UPLOAD:
      return state.set('uploadedPhoto', action.photo);
    case UPLOAD_RESET:
      return initialState;
    default:
      return state;
  }
}




