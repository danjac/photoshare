import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Input,
         ButtonInput,
         Alert
        } from 'react-bootstrap';

import { Loader } from './util';

import * as ActionCreators from '../actions';

@connect(state => state.recoverPassword.toJS())
export default class RecoverPassword extends React.Component {

  static propTypes = {
    dispatch: PropTypes.func.isRequired
  }

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = Object.assign({},
      bindActionCreators(ActionCreators.recoverPassword, dispatch)
    );
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentDidMount() {
    this.actions.resetForm();
  }

  handleSubmit(event) {
    event.preventDefault();
    const email = this.refs.email.getValue().trim();
    if (email) {
      this.actions.recoverPassword(email);
      this.refs.email.getInputDOMNode().value = "";
    }
  }

  render() {

    if (this.props.formSubmitted) {
      return <Loader />;
    }

    let msg = '';

    if (this.props.isError) {
      msg = <Alert bsStyle="warning">Sorry, we were not able to find your email address.</Alert>;
    } else if (this.props.isSuccess) {
      msg = <Alert bsStyle="success">Please check your email inbox for instructions on how to recover your password.</Alert>;
    } else {
      msg = <Alert bsStyle="info">Please enter your email address and we'll send you a link to recover your password.</Alert>;
    }

    return (
      <div className="col-md-6 col-md-offset-3">
        {msg}
          <form role="form" method="POST" onSubmit={this.handleSubmit}>
              <Input type="email" ref="email" required placeholder="Email address" />
              <ButtonInput bsStyle="primary" type="submit">Continue</ButtonInput>
          </form>
      </div>
    );
  }

}
