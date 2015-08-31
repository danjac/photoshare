import React, { PropTypes } from 'react';
import _ from 'lodash';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Navbar, 
         Nav, 
         Alert,
         NavDropdown, 
         NavItem, 
         MenuItem } from 'react-bootstrap';
import { Link } from 'react-router';

import * as ActionCreators from './actions';

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
        <NavItem href={makeHref('/login/')}><i className="fa fa-sign-in"></i> Login</NavItem>
        <NavItem href="/"><i className="fa fa-user"></i> Signup</NavItem>
      </Nav>
    );
  }

  render() {

    const brand = <Link to="/"><i className="fa fa-camera"></i> Wallshare</Link>;
    const searchIcon = <i className="fa fa-search"></i>;
    const makeHref = this.context.router.makeHref; 

    return (

      <Navbar fixedTop inverse brand={brand}>

        <Nav>
          <NavItem href={makeHref("/")}><i className="fa fa-fire"></i> Popular</NavItem>
          <NavItem href={makeHref("/latest/")}><i className="fa fa-clock-o"></i> Latest</NavItem>
          <NavItem href="/"><i className="fa fa-tags"></i> Tags</NavItem>
          <NavItem href="/"><i className="fa fa-upload"></i> Upload</NavItem>
        </Nav>
        <Nav>
          <form className="navbar-form navbar-left" role="search" name="searchForm">
            <Input type="text" addonAfter={searchIcon} bsSize="small" placeholder="Search" />
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


