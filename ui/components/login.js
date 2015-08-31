import React, { PropTypes } from 'react';
import _ from 'lodash';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Input, 
         ButtonInput 
        } from 'react-bootstrap';

import * as ActionCreators from '../actions';


@connect(state => state.auth.toJS())
export default class Login extends React.Component {

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
      bindActionCreators(ActionCreators.messages, dispatch),
      bindActionCreators(ActionCreators.auth, dispatch)
    );
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  shouldComponentUpdate(nextProps) {
    if (nextProps.loggedIn) {
      this.actions.newMessage(`Welcome back, ${nextProps.name}`, "success");
      //const nextPath = this.props.query.nextPath || "/";
      //console.log("query", this.props.query);
      this.context.router.transitionTo("/upload/");
      return true;
    }
    return false;
  }

  handleSubmit(event) {
    event.preventDefault();
    const identifier = this.refs.identifier.getValue(),
          password = this.refs.password.getValue();

    this.refs.password.getInputDOMNode().value = "";
    this.actions.login(identifier, password);
  }

  render() {
    return (
      <div className="col-md-6 col-md-offset-3">
        <form role="form" method="POST" onSubmit={this.handleSubmit}>
            <Input type="text" ref="identifier" required placeholder="Name or email address" />
            <Input type="password" ref="password" required placeholder="Password" />
            <ButtonInput bsStyle="primary" type="submit">Login</ButtonInput>
        </form>

        <a href="#/recoverpass">Forgot your password?</a> 
      </div>
    );
  }
}


