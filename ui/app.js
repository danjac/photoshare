import React from 'react';
import { Provider } from 'react-redux';
import HashHistory from 'react-router/lib/HashHistory';

import routes from './routes';
import configureStore from './store';


const store = configureStore(),
      history = new HashHistory();


class Container extends React.Component {
  render() {
    return (
    <div>
    <Provider store={store}>
    {() => routes(store, history)}
    </Provider>
    </div>
    );
  }
}

React.render(<Container />, document.body);
