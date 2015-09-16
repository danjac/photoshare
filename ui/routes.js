import React from 'react';
import { Router, Route } from 'react-router';

import {
  App,
  Popular,
  Latest,
  Search,
  User,
  PhotoDetail,
  Login,
  Signup,
  RecoverPassword,
  Upload,
  TagList
} from './components';

export default function(store, history) {
  const requireAuth = (nextState, replaceState) => {
    const auth = store.getState().auth.toJS();
    if (!auth.loggedIn) {
      replaceState.to('/login/', { nextPath: nextState.location.pathname });
    }
  }
  return (
    <Router history={history}>
      <Route component={App}>
        <Route path="/" component={Popular} />
        <Route path="/upload/" component={Upload} onEnter={requireAuth} />
        <Route path="/latest/" component={Latest} />
        <Route path="/search/" component={Search} />
        <Route path="/tags/" component={TagList} />
        <Route path="/login/" component={Login} />
        <Route path="/signup/" component={Signup} />
        <Route path="/recoverpass/" component={RecoverPassword} />
        <Route path="/detail/:id" component={PhotoDetail} />
        <Route path="/user/:userID/:username" component={User} />
      </Route>
    </Router>
  );

}
