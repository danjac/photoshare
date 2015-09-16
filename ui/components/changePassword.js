import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Input,
         ButtonInput,
         Alert
        } from 'react-bootstrap';

import { Loader } from './util';

import * as ActionCreators from '../actions';

@connect(state => state.changePassword.toJS())
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
      this.context.router.transitionTo("/login/");
    }
  }

  handleSubmit(event) {
    event.preventDefault();

    const password = this.refs.password.getValue().trim();
    const passwordConfirm = this.refs.password.getValue().trim();

    const code = this.context.router.state.query.code;

    this.refs.password.getInputDOMNode().value = "";
    this.refs.passwordConfirm.getInputDOMNode().value = "";

    if (password && passwordConfirm) {
      this.actions.submitForm(password, passwordConfirm, code);
    }
  }

  render() {

    if (this.props.formSubmitted) {
      return <Loader />;
    }

    return (
      <div className="col-md-6 col-md-offset-3">
        {msg}
          <form role="form" method="POST" onSubmit={this.handleSubmit}>
              <Input type="password" ref="password" required placeholder="Password" />
              <Input type="password" ref="password" required placeholder="Repeat password" />
              <ButtonInput bsStyle="primary" type="submit">Continue</ButtonInput>
          </form>
      </div>
    );
  }

}
