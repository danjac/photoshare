/* jslint ignore:start */
import Immutable from 'immutable';

import ActionTypes from '../actionTypes';

const {
  PHOTO_PREVIEW,
  UPLOAD_SUCCESS,
  UPLOAD_FAILURE,
  UPLOAD_PENDING,
  OK
} = ActionTypes;


const initialState = Immutable.fromJS({
  previewURL: null,
  uploadedPhoto: null,
  formSubmitted: false,
  isWaiting: false
});

export default function(state=initialState, action) {
  switch(action.type) {
    case OK:
      if(action.form === 'upload' && action.field === 'photo') {
        return state.set('previewURL', action.previewURL);
      }
      return state;
    case UPLOAD_PENDING:
      return state.merge({
          formSubmitted: true,
          isWaiting: true,
          previewURL: null
      });
    case UPLOAD_SUCCESS:
      return state
        .set("isWaiting", false)
        .set('uploadedPhoto', action.payload);
    case UPLOAD_FAILURE:
        return state.set("isWaiting", false)
    default:
      return state;
  }
}
