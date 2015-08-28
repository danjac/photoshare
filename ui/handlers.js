import React from 'react';
import { RouteHandler } from 'react-router';



class NavBar extends React.Component {

  render() {
    return (
    <nav className="navbar navbar-inverse navbar-fixed-top" role="navigation">
        <div className="container-fluid">
            <div className="navbar-header">
                <button type="button" className="navbar-toggle" data-toggle="collapse" data-target="#navbar-links">
                    <span className="sr-only">Toggle navigation</span>
                    <span className="icon-bar"></span>
                    <span className="icon-bar"></span>
                    <span className="icon-bar"></span>
                </button>
                <a className="navbar-brand" href="front"><i className="fa fa-camera"></i> Wallshare</a>
            </div>
            <div className="collapse navbar-collapse" id="navbar-links">
                <ul className="nav navbar-nav navbar-left">
                    <li href-active="popular"><a href="popular"><i className="fa fa-fire"></i> Popular</a>
                    </li>
                    <li href-active="latest"><a href="latest"><i className="fa fa-clock-o"></i> Latest</a>
                    </li>
                    <li href-active="active"><a href="tags"><i className="fa fa-tags"></i> Tags</a>
                    </li>
                    <li href-active="active"><a href="upload"><i className="fa fa-upload"></i> Upload</a>
                    </li>
                </ul>
                <form className="navbar-form navbar-left" role="search" name="searchForm" novalidate>
                    <div className="form-group">
                        <input type="text" className="form-control input-sm" placeholder="Search" data-toggle="tooltip" data-placement="bottom" title="Prefix search with '#' for tags and '@' for users" required />
                        <button type="submit" className="btn btn-default btn-sm"><i className="fa fa-search"></i>
                        </button>
                    </div>
                </form>

                <ul className="nav navbar-nav navbar-right">
                    <li href-active="login"><a><i className="fa fa-log-in"></i> Login</a>
                    </li>
                    <li href-active="signup"><a href="#/signup"><i className="fa fa-user"></i> Signup</a>
                    </li>
                </ul>
                {/*
                <ul className="nav navbar-nav navbar-right">
                    <li className="dropdown">
                        <a className="dropdown-toggle" data-toggle="dropdown">
                            <img gravatar-src="session.email" gravatar-size="20" />&nbsp;{{session.name}} <i className="caret"></i>
                        </a>
                        <ul className="dropdown-menu" role="menu">
                            <li><a href="owner({ownerID: session.id, ownerName: session.name})">My photos</a>
                            </li>
                            <li><a ng-href="#/changepass">Change my password</a>
                            </li>
                            <li><a ng-click="logout()">Logout</a>
                            </li>
                        </ul>
                    </li>
                </ul>
                */}
            </div>
        </div>
    </nav>
    );
  }


}

export class App extends React.Component {
  render() {
    return (
      <div className="container-fluid">
        <NavBar />
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
