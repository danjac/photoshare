import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Input,
         ButtonInput
        } from 'react-bootstrap';

import * as ActionCreators from '../actions';


@connect(state => state.auth.toJS())
export default class Signup extends React.Component {

  static propTypes = {
    dispatch: PropTypes.func.isRequired
  }

  static contextTypes = {
    router: PropTypes.object.isRequired
  }

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(ActionCreators.auth, dispatch);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  shouldComponentUpdate(nextProps) {
    if (nextProps.loggedIn) {
      this.context.router.transitionTo("/upload/");
      return true;
    }
    return nextProps !== this.props;
  }

  handleSubmit(event) {
    event.preventDefault();

    const name = this.refs.name.getValue(),
          email = this.refs.email.getValue(),
          password = this.refs.password.getValue();

    this.actions.signup(name, email, password);

  }

  render() {
    return (
    <div className="col-md-6 col-md-offset-3">
      <form role="form"
            name="form"
            method="POST"
            onSubmit={this.handleSubmit}>

          <Input ref="name"
                 type="text"
                 placeholder="Name"
                 required />

          <Input ref="email"
                 type="email"
                 placeholder="Email address"
                 required />

          <Input ref="password"
                 type="password"
                 placeholder="Password"
                 required />

          <ButtonInput type="submit" bsStyle="primary">Continue</ButtonInput>

      </form>
    </div>
    );
  }

}
