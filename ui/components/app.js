/* jslint ignore:start */
import React, { PropTypes } from 'react';
//import CSSModules from 'react-css-modules';
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

// import styles from '../app.css';

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
    const { id, name, loggedIn } = this.props.auth;
    const router = this.context.router;
    const makeHref = router.makeHref;
    const handleLogout  = this.handleLogout.bind(this);
    const currentPath = router.state.location.pathname;

    if (loggedIn) {
      return (
        <Nav right>
          <NavDropdown title={name} id="userDropdown">
            <MenuItem href={makeHref(`/user/${id}/${name}`)}>My photos</MenuItem>
            <MenuItem>Change my password</MenuItem>
            <MenuItem onSelect={handleLogout}>Logout</MenuItem>
          </NavDropdown>
        </Nav>
      );
    }
    return (
      <Nav right>
        <NavItem active={this.isActive('/login/')} href={makeHref('/login/', { nextPath: currentPath })}><Facon name='sign-in' /> Login</NavItem>
        <NavItem active={this.isActive('/signup/')} href={makeHref('/signup/')}><Facon name='user' /> Signup</NavItem>
      </Nav>
    );
  }

  isActive(path) {
    return this.context.router.isActive(path);
  }

  render() {

    const brand = <Link to="/"><Facon name='camera' /> Wallshare</Link>;
    const searchIcon = <Facon name='search' />
    const makeHref = this.context.router.makeHref;
    const handleSearch = this.handleSearch.bind(this);


    return (

      <Navbar fixedTop inverse brand={brand}>

        <Nav>
          <NavItem active={this.isActive('/latest/')} href={makeHref("/latest/")}><Facon name='clock-o' /> Latest</NavItem>
          <NavItem active={this.isActive('/tags/')} href={makeHref("/tags/")}><Facon name='tags' /> Tags</NavItem>
          <NavItem active={this.isActive('/upload/')} href={makeHref("/upload/")}><Facon name='upload' /> Upload</NavItem>
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

//@CSSModules(styles)
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
