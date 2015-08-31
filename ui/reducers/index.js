import { combineReducers } from 'redux';

import photos from './photos';
import photoDetail from './photoDetail';
import auth from './auth';
import messages from './messages';
import upload from './upload';

export default combineReducers({
  photos,
  photoDetail,
  auth,
  messages,
  upload
});


