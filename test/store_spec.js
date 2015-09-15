import Immutable from 'immutable';
import { expect } from 'chai';

import { ActionTypes } from '../ui/constants';
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

  it('deletes a photo', () => {

    const store = configureStore();

    store.dispatch({
      type: ActionTypes.FETCH_PHOTO_DETAIL_SUCCESS,
      payload: {
        title: "test"
      }
    });

    expect(store.getState().photoDetail.get("isDeleted")).to.equal(false);
    expect(store.getState().photoDetail.get("isLoaded")).to.equal(true);

    store.dispatch({
      type: ActionTypes.DELETE_PHOTO_PENDING
    });

    expect(store.getState().photoDetail.get("isDeleted")).to.equal(true);

  });

});
