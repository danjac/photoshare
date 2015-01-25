var React = require('react');
var Router = require('react-router');

var Constants = require('../Constants');
var Actions = require('../Actions');
var UserStore = require('../stores/UserStore');
var AlertStore = require('../stores/AlertStore');

var RouteHandler = Router.RouteHandler;
var Link = Router.Link;

var Alert = React.createClass({
    render: function (){
        var className = "alert alert-dismissable alert-";
        switch (this.props.message.type) {
            case Constants.ALERT_SUCCESS:
                className += "success";
                break;
            case Constants.ALERT_WARNING:
                className += "warning";
                break;
            case Constants.ALERT_DANGER:
                className += "danger";
                break;
            default:
                className += "info";
          }
        return (
            <div className={className} role="alert">
            <button type="button" className="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>
            {this.props.message.message}
            </div>
        );
    }
});

var NavbarLoggedIn = React.createClass({

    handleLogout: function() {
        Actions.logout();
    },

    render: function() {
        return (
                <ul className="nav navbar-nav navbar-right">
                    <li className="dropdown">
                        <a className="dropdown-toggle" data-toggle="dropdown">
                            {this.props.user.name} <i className="caret"></i>
                        </a>
                        <ul className="dropdown-menu" role="menu">
                            <li><a href="">My photos</a>
                            </li>
                            <li><a href="#/changepass">Change my password</a>
                            </li>
                            <li><a onClick={this.handleLogout}>Logout</a>
                            </li>
                        </ul>
                    </li>
                </ul>
        );
    }
});


var NavbarLoggedOut = React.createClass({

    render: function () {

        return (
                <ul className="nav navbar-nav navbar-right">
                    <li><Link to="login"><i className="fa fa-log-in"></i> Login</Link>
                    </li>
                    <li><a href="#/signup"><i className="fa fa-user"></i> Signup</a>
                    </li>
                </ul>
        );
    }
});

var Navbar = React.createClass({
    render: function(){

        if (this.props.user) {
            loginButtons = <NavbarLoggedIn user={this.props.user} />;
        } else {
            loginButtons = <NavbarLoggedOut />
        }

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
                <Link to="popular" className="navbar-brand"><i className="fa fa-camera"></i> Wallshare</Link>
            </div>
            <div className="collapse navbar-collapse" id="navbar-links">
                <ul className="nav navbar-nav navbar-left">
                    <li><Link to="popular"><i className="fa fa-fire"></i> Popular</Link>
                    </li>
                    <li><Link to="latest"><i className="fa fa-clock-o"></i> Latest</Link>
                    </li>
                    <li><a href="tags"><i className="fa fa-tags"></i> Tags</a>
                    </li>
                    <li><a href="upload"><i className="fa fa-upload"></i> Upload</a>
                    </li>
                </ul>
                <form className="navbar-form navbar-left" role="search" name="searchForm">
                    <div className="form-group">
                        <input type="text" className="form-control input-sm" placeholder="Search" data-toggle="tooltip" data-placement="bottom" title="Prefix search with '#' for tags and '@' for users" required />
                        <button type="submit" className="btn btn-default btn-sm"><i className="fa fa-search"></i>
                        </button>
                    </div>
                </form>

                {loginButtons}

            </div>
        </div>
    </nav>
    );
    }

});

var App = React.createClass({

    getInitialState: function() {
        return {
            user: null,
            messages: []
       };
    },

    componentWillMount: function() {
        UserStore.addChangeListener(this._onChange);
        AlertStore.addChangeListener(this._onChange);
        Actions.getUser();
    },

    componentWillUnmount: function () {
        UserStore.removeChangeListener(this._onChange);
        AlertStore.removeChangeListener(this._onChange);
    },

    render: function() {
        return (
    <div>
        <Navbar user={this.state.user}/>
        <div className="container-fluid">
        <div>
        {this.state.messages.map(function(msg, num){
            return <Alert key={num} message={msg} />
        })}
        </div>
        <RouteHandler user={this.state.user} />
        </div>
    </div>
        );
    },

    _onChange: function() {
        this.setState({
            user: UserStore.getUser(),
            messages: AlertStore.getMessages()
        });
    }

});

module.exports = App;
