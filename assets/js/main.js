/** @jsx React.DOM */

var React = require('react');
var Router = require('react-router');
var Route = Router.Route;
var RouteHandler = Router.RouteHandler;
var DefaultRoute = Router.DefaultRoute;

var App = require('./components/App.jsx');
var Popular = require('./components/Popular.jsx');
var Latest = require('./components/Latest.jsx');

var routes = (
    <Route handler={App}>
        <DefaultRoute name="popular" handler={Popular} />
        <Route name="latest" path="latest" handler={Latest} />
    </Route>
    );

Router.run(routes, function (Handler) {
  React.render(<Handler/>, document.body);
});
