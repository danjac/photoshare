/* jslint ignore:start */
import React, { PropTypes } from 'react';
import _ from 'lodash';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Navbar,
         Nav,
         Alert,
         Input,
         NavDropdown,
         NavItem,
         MenuItem } from 'react-bootstrap';
import { Link } from 'react-router';

import * as ActionCreators from '../actions';
import { Facon } from './util';

class Messages extends React.Component {

  static propTypes = {
    messages: PropTypes.array.isRequired,
    actions: PropTypes.object.isRequired
  }

  render() {
    return (
    <div>
    {this.props.messages.map((msg, index) => {

      const handleDelete = () => {
        this.props.actions.deleteMessage(index);
      }

      return <Alert key={index}
                    onDismiss={handleDelete}
                    dismissAfter={2000}
                    bsStyle={msg.level}>{msg.msg}</Alert>;
    })}
    </div>
    );
  }

}

class Navigation extends React.Component {

  static propTypes = {
    auth: PropTypes.object.isRequired
  }

  static contextTypes = {
    router: PropTypes.object.isRequired
  }

  handleLogout(event) {
    event.preventDefault();
    this.props.actions.newMessage("Bye for now!", "success");
    this.props.actions.logout();
    this.context.router.transitionTo("/");
  }

  handleSearch(event) {
    event.preventDefault();
    const query = this.refs.query.getValue();
    this.refs.query.getInputDOMNode().value = "";
    this.context.router.transitionTo('/search/', { q: query });
  }

  rightNav() {
    const { name, loggedIn } = this.props.auth;
    const makeHref = this.context.router.makeHref;
    const handleLogout  = this.handleLogout.bind(this);

    if (loggedIn) {
      return (
        <Nav right>
          <NavDropdown title={name}>
            <MenuItem>My photos</MenuItem>
            <MenuItem>Change my password</MenuItem>
            <MenuItem onSelect={handleLogout}>Logout</MenuItem>
          </NavDropdown>
        </Nav>
      );
    }
    return (
      <Nav right>
        <NavItem href={makeHref('/login/')}><Facon name='sign-in' /> Login</NavItem>
        <NavItem href="/"><Facon name='user' /> Signup</NavItem>
      </Nav>
    );
  }

  render() {

    const brand = <Link to="/"><Facon name='camera' /> Wallshare</Link>;
    const searchIcon = <Facon name='search' />
    const makeHref = this.context.router.makeHref;
    const handleSearch = this.handleSearch.bind(this);

    const isActive = (path, q) => this.context.router.isActive(path, q);

    return (

      <Navbar fixedTop inverse brand={brand}>

        <Nav>
          <NavItem active={isActive('/')} href={makeHref("/")}><Facon name='fire' /> Popular</NavItem>
          <NavItem active={isActive('/latest/')} href={makeHref("/latest/")}><Facon name='clock-o' /> Latest</NavItem>
          <NavItem active={isActive('/tags/')} href="/"><Facon name='tags' /> Tags</NavItem>
          <NavItem active={isActive('/upload/')} href={makeHref("/upload/")}><Facon name='upload' /> Upload</NavItem>
        </Nav>

        <Nav>
          <form onSubmit={handleSearch} className="navbar-form navbar-left" role="search" name="searchForm">
            <Input type="text" ref="query" addonAfter={searchIcon} bsSize="small" placeholder="Search" />
          </form>
        </Nav>
        {this.rightNav()}
      </Navbar>
    );
  }


}

@connect(state => {
  return {
    auth: state.auth.toJS(),
    messages: state.messages.toJS()
  }
})
export default class App extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = Object.assign({},
      bindActionCreators(ActionCreators.auth, dispatch),
      bindActionCreators(ActionCreators.messages, dispatch)
    );
  }

  componentDidMount() {
    this.actions.getUser();
  }

  render() {
    return (
    <div className="container-fluid">
      <Navigation actions={this.actions} {...this.props} />
      <Messages actions={this.actions} {...this.props} />
      {this.props.children}
    </div>
    );
  }
}
