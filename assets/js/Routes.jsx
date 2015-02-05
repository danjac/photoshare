/** @jsx React.DOM */

var React = require('react');
var Router = require('react-router');
var Route = Router.Route;
var RouteHandler = Router.RouteHandler;
var DefaultRoute = Router.DefaultRoute;

var App = require('./components/App.jsx');
var Popular = require('./components/Popular.jsx');
var Latest = require('./components/Latest.jsx');
var Login = require('./components/Login.jsx');
var Upload = require('./components/Upload.jsx');
var Search = require('./components/Search.jsx');
var PhotoDetail = require('./components/PhotoDetail.jsx');
var Tags = require('./components/Tags.jsx');
var User = require('./components/User.jsx');

var routes = (
    <Route handler={App}>
        <DefaultRoute name="popular" handler={Popular} />
        <Route name="latest" path="latest" handler={Latest} />
        <Route name="search" path="search" handler={Search} />
        <Route name="tags" path="tags" handler={Tags} />
        <Route name="upload" path="upload" handler={Upload} />
        <Route name="login" path="login" handler={Login} />
        <Route name="user" path="user/:id" handler={User} />
        <Route name="photoDetail" path="photo/:id" handler={PhotoDetail} />
    </Route>
    );

module.exports = routes;

