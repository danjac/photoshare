/** @jsx React.DOM */

var React = require('react');
var Router = require('react-router');
var Routes = require('./Routes.jsx');
var data = JSON.parse(document.getElementById("initData").innerHTML);

Router.run(Routes, Router.HistoryLocation, function (Handler) {
    React.render(<Handler photos={data} />, document.body);
});
