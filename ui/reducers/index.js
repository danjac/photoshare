import { combineReducers } from 'redux';

import photos from './photos';
import photoDetail from './photoDetail';
import auth from './auth';

export default combineReducers({
  photos,
  photoDetail,
  auth
});


