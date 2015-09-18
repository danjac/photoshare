import * as auth from './auth';
import * as changePassword from './changePassword';
import * as forms from './forms';
import * as messages from './messages';
import * as photoDetail from './photoDetail';
import * as photos from './photos';
import * as recoverPassword from './recoverPassword';
import * as tags from './tags';
import * as upload from './upload';

export default Object.assign({},
  auth,
  changePassword,
  forms,
  messages,
  photoDetail,
  photos,
  recoverPassword,
  tags,
  upload
);
