import Immutable from 'immutable';
import sinon from 'sinon';
import {expect} from 'chai';

import {ActionTypes} from '../ui/constants';
import configureStore from '../ui/store';

describe('store', () => {

  it('handles login', () => {
    const store = configureStore();
    expect(store.getState().auth.get("loggedIn")).to.equal(false);

    store.dispatch({
      type: ActionTypes.LOGIN_SUCCESS,
      payload: {
        loggedIn: true,
        id: 1,
        name: 'test',
        email: 'test'
      }
    });

    expect(store.getState().auth.get("loggedIn")).to.equal(true);

  });

});
