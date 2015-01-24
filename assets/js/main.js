/** @jsx React.DOM */

var React = require('react');
var Router = require('react-router');
var Route = Router.Route;
var RouteHandler = Router.RouteHandler;
var DefaultRoute = Router.DefaultRoute;

var App = require('./components/App.jsx');
var Popular = require('./components/Popular.jsx');
var Latest = require('./components/Latest.jsx');

var API = require('./API.js')

var routes = (
    <Route handler={App}>
        <DefaultRoute name="popular" handler={Popular} />
        <Route name="latest" path="latest" handler={Latest} />
    </Route>
    );

var fetchData = function(callback) {

    API.getPhotos("votes", function(response){
        callback({ photos: response });
    });
};


Router.run(routes, function (Handler) {
    fetchData(function(data){
        React.render(<Handler data={data} />, document.body);
    });
});
