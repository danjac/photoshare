import React from 'react';
import { DefaultRoute, Route } from 'react-router';

import { App, Popular } from './handlers';

export default (
  <Route name="app" path="/" handler={App}>
    <DefaultRoute handler={Popular} />
  </Route>
);
