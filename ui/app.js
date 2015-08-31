import React from 'react';
import { Router, Route } from 'react-router';
import HashHistory from 'react-router/lib/HashHistory';
import { Provider } from 'react-redux';

import { 
  App, 
  Popular, 
  Latest,
  PhotoDetail,
  Login
} from './components';

import configureStore from './store';


const store = configureStore(), 
      history = new HashHistory();


class Container extends React.Component {
  render() {
    return (
    <div>
    <Provider store={store}>
    {() =>
      <Router history={history}>
        <Route component={App}>
          <Route path="/" component={Popular} />
          <Route path="/latest/" component={Latest} />
          <Route path="/detail/:id" component={PhotoDetail} />
          <Route path="/login/" component={Login} />
        </Route>
      </Router>
    }
    </Provider>
    </div>
    );
  }
}

React.render(<Container />, document.body);
