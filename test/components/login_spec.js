import React from 'react/addons';
import { expect } from 'chai';
import Login from '../../ui/components/login.js';
import configureStore from '../../ui/store';

const { renderIntoDocument, scryRenderedDOMComponentsWithTag } = React.addons.TestUtils;

describe('Login', () => {
  it('should render a component', () => {

    const store = configureStore();
    const component = renderIntoDocument(<Login store={store} />);
    const forms = scryRenderedDOMComponentsWithTag(component, 'form');
    expect(forms.length).to.equal(1);

  });
});
