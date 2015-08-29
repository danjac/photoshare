import { createStore, applyMiddleware } from 'redux';
import thunkMiddleware from 'redux-thunk';
import loggingMiddleware from 'redux-logger';
import reducer from './reducers';

const createStoreWithMiddleware = applyMiddleware(
  thunkMiddleware,
  loggingMiddleware
)(createStore);

export default function configureStore(initialState) {
  return createStoreWithMiddleware(reducer, initialState);
}
