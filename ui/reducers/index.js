import { combineReducers } from 'redux';

import photos from './photos';
import photoDetail from './photoDetail';
import auth from './auth';
import messages from './messages';

export default combineReducers({
  photos,
  photoDetail,
  auth,
  messages
});


