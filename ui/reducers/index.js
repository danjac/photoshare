import { combineReducers } from 'redux';

import photos from './photos';
import photoDetail from './photoDetail';

export default combineReducers({
  photos,
  photoDetail
});


