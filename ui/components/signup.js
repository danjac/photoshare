import React, { PropTypes } from 'react';
import _ from 'lodash';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Input,
         ButtonInput
        } from 'react-bootstrap';

import * as ActionCreators from '../actions';


@connect(state => {
  let props = state.auth.toJS();
  let forms = state.forms.toJS();
  props.signupPrechecks = forms.signup ? forms.signup.checked : [];
  props.signupErrors = forms.signup ? forms.signup.errors : {};
  return props;
})
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
    this.handleCheckName = this.handleCheckName.bind(this);
    this.handleCheckPassword = this.handleCheckPassword.bind(this);
    this.handleCheckEmail = this.handleCheckEmail.bind(this);
  }

  componentDidMount() {
    this.actions.resetSignupForm();
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

    if (!_.isEmpty(this.props.signupErrors)) {
      return;
    }

    const name = this.refs.name.getValue(),
          email = this.refs.email.getValue(),
          password = this.refs.password.getValue();

    this.actions.signup(name, email, password);

  }

  handleCheckEmail(event) {
    event.preventDefault();
    const email = this.refs.email.getValue().trim();
    if (email) {
      this.actions.checkEmail(email);
    }
  }

  handleCheckName(event) {
    event.preventDefault();
    const name = this.refs.name.getValue().trim();
    if (name) {
      this.actions.checkName(name);
    }
  }

  handleCheckPassword(event) {
    event.preventDefault();
    const password = this.refs.password.getValue().trim();
    if (password) {
      this.actions.checkPassword(password);
    }
  }

  errorStatus(name) {
    console.log(this.props.signupErrors);
    if (this.props.signupPrechecks.indexOf(name) === -1) {
        return;
    }
    return this.props.signupErrors[name] ? 'error' : 'success';
  }

  errorMsg(name) {
    return this.props.signupErrors[name] || '';
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
                 onBlur={this.handleCheckName}
                 hasFeedback
                 bsStyle={this.errorStatus('name')}
                 help={this.errorMsg('name')}
                 required />

          <Input ref="email"
                 type="email"
                 placeholder="Email address"
                 onBlur={this.handleCheckEmail}
                 hasFeedback
                 bsStyle={this.errorStatus('email')}
                 help={this.errorMsg('email')}
                 required />

          <Input ref="password"
                 type="password"
                 placeholder="Password"
                 onBlur={this.handleCheckPassword}
                 hasFeedback
                 bsStyle={this.errorStatus('password')}
                 help={this.errorMsg('password')}
                 required />

               <ButtonInput enabled={!this.props.signupErrors}
                            type="submit" bsStyle="primary">Continue</ButtonInput>

      </form>
    </div>
    );
  }

}
