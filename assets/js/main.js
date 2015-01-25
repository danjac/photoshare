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
var PhotoDetail = require('./components/PhotoDetail.jsx');

var routes = (
    <Route handler={App}>
        <DefaultRoute name="popular" handler={Popular} />
        <Route name="latest" path="latest" handler={Latest} />
        <Route name="login" path="login" handler={Login} />
        <Route name="photoDetail" path="photo/:id" handler={PhotoDetail} />
    </Route>
    );

Router.run(routes, function (Handler) {
    React.render(<Handler />, document.body);
});
