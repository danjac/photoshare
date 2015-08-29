import React from 'react';
import { Navbar, Nav, Input } from 'react-bootstrap';
import { NavItemLink } from 'react-router-bootstrap';
import { RouteHandler, Link } from 'react-router';



class Navigation extends React.Component {

  render() {

    const brand = <Link to="app"><i className="fa fa-camera"></i> Wallshare</Link>;
    const searchIcon = <i className="fa fa-search"></i>;

    return (

      <Navbar fixedTop inverse brand={brand}>

        <Nav>
          <NavItemLink to="app"><i className="fa fa-fire"></i> Popular</NavItemLink>
          <NavItemLink to="latest"><i className="fa fa-clock-o"></i> Latest</NavItemLink>
          <NavItemLink to="app"><i className="fa fa-tags"></i> Tags</NavItemLink>
          <NavItemLink to="app"><i className="fa fa-upload"></i> Upload</NavItemLink>
        </Nav>
        <Nav>
          <form className="navbar-form navbar-left" role="search" name="searchForm" novalidate>
            <Input type="text" addonAfter={searchIcon} bsSize="small" placeholder="Search" />
          </form>
        </Nav>
        <Nav right>
          <NavItemLink to="app"><i className="fa fa-sign-in"></i> Login</NavItemLink>
          <NavItemLink to="app"><i className="fa fa-user"></i> Signup</NavItemLink>
        </Nav>
      </Navbar>
    );
  }


}

export class App extends React.Component {
  render() {
    return (
      <div className="container-fluid">
        <Navigation />
        <RouteHandler />
      </div>
    );
  }
}

export class Popular extends React.Component {
  render() {
    return (
      <div>Popular photos go here...</div>
    );
  }

}

export class Latest extends React.Component {
  render() {
    return (
      <div>Latest photos go here...</div>
    );
  }

}
