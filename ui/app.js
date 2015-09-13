import React from 'react';
import { Provider } from 'react-redux';
import { Router, Route } from 'react-router';
import HashHistory from 'react-router/lib/HashHistory';

import {
  App,
  Popular,
  Latest,
  Search,
  PhotoDetail,
  Login,
  Signup,
  Upload,
  TagList
} from './components';

import configureStore from './store';


const store = configureStore(),
      history = new HashHistory();


class Container extends React.Component {
  render() {
    return (
    <div>
    <Provider store={store}>
    {() => {
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
          <Route path="/detail/:id" component={PhotoDetail} />
        </Route>
      </Router>
      );
    }}
    </Provider>
    </div>
    );
  }
}

React.render(<Container />, document.body);
