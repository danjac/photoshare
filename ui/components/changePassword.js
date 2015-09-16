import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import { Input,
         ButtonInput,
         Alert,
         Well
        } from 'react-bootstrap';

import { Loader } from './util';

import * as ActionCreators from '../actions';

@connect(state => {
  const props = state.changePassword.toJS();
  props.loggedIn = state.auth.get("loggedIn");
  return props;
})
export default class ChangePassword extends React.Component {

  static propTypes = {
    dispatch: PropTypes.func.isRequired
  }

  static contextTypes = {
    router: PropTypes.object.isRequired
  }

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = Object.assign({},
      bindActionCreators(ActionCreators.changePassword, dispatch)
    );
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentDidMount() {
    this.actions.resetForm();
  }

  shouldComponentUpdate(nextProps) {
    if (nextProps.isSuccess) {
      const nextPath = this.props.loggedIn ? "/" : "/login/";
      this.context.router.transitionTo(nextPath);
    }
    return nextProps !== this.props;
  }

  handleSubmit(event) {
    event.preventDefault();

    const password = this.refs.password.getValue().trim();
    const passwordConfirm = this.refs.passwordConfirm.getValue().trim();

    const code = this.getRecoveryCode();

    this.refs.password.getInputDOMNode().value = "";
    this.refs.passwordConfirm.getInputDOMNode().value = "";

    if (password && passwordConfirm) {
        this.actions.submitForm(password, passwordConfirm, code, this.props.loggedIn);
    }
  }

  getRecoveryCode() {
    return this.props.location.query ? this.props.location.query.code : null;
  }

  errorStatus(name) {
    if (!this.props.formSubmitted) {
        return;
    }
    return this.props.errors.has(name) ? 'error' : 'success';
  }

  errorMsg(name) {
    return this.props.errors.get(name) || '';
  }

  render() {

    if (this.props.isWaiting) {
      return <Loader />;
    }

    const code = this.getRecoveryCode();

    if (!code && !this.props.loggedIn) {
      return (
        <Well className="col-md-6 col-md-offset-3">
          You must have a valid recovery code or be logged in to view this page. If you are not logged in please go to <Link to="/recoverpass/">this page</Link> to get a new code. Otherwise you can sign in <Link to="/login/?nextPath=/changepass/">here</Link>.
        </Well>
      )
    }

    return (
      <div className="col-md-6 col-md-offset-3">
          <form role="form" method="POST" onSubmit={this.handleSubmit}>
            <Input type="password"
              ref="password"
              placeholder="Password"
              required
              hasFeedback
              bsStyle={this.errorStatus('password')}
              help={this.errorMsg('password')} />

            <Input type="password"
               ref="passwordConfirm"
               placeholder="Repeat password"
               required
               hasFeedback
               bsStyle={this.errorStatus('passwordConfirm')}
               help={this.errorMsg('passwordConfirm')} />

            <ButtonInput bsStyle="primary" type="submit">Continue</ButtonInput>
          </form>
      </div>
    );
  }

}
