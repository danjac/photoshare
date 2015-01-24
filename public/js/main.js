/** @jsx React.DOM */

var React = require('react');
// Here we put our React instance to the global scope. Make sure you do not put it
// into production and make sure that you close and open your console if the
// DEV-TOOLS does not display
// window.React = React;

var App = require('./components/App.jsx');
React.renderComponent(<App/>, document.body);
