import React from 'react';
import { DefaultRoute, Route } from 'react-router';

import { 
  App, 
  Popular, 
  Latest,
  PhotoDetail,
  Login
} from './handlers';

export default (
  <Route name="app" path="/" handler={App}>
    <DefaultRoute handler={Popular} />
    <Route name="latest" path="/latest/" handler={Latest} />
    <Route name="detail" path="/detail/:id" handler={PhotoDetail} />
    <Route name="login" path="/login/" handler={Login} />
  </Route>
);
