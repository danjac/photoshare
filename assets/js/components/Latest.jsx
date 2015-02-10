var React = require('react');
var Actions = require('../Actions');
var PhotoStore = require('../stores/PhotoStore');
var PhotoList = require('./PhotoList.jsx')

var Latest = React.createClass({

    statics: {
        willTransitionTo: function(transition, params) {
            Actions.getPhotos();
        }
    },

    getInitialState: function() {
        return {
            photos: {
                photos: []
            }
        }
    },

    componentWillMount: function() {
        PhotoStore.addChangeListener(this._onChange);
    },

    componentWillUnmount: function() {
        PhotoStore.removeChangeListener(this._onChange);
    },

    handlePaginationLink: function(page) {
        Actions.getPhotos(null, page);
    },

    render: function() {
        return (
            <PhotoList photos={this.state.photos} handlePaginationLink={this.handlePaginationLink}/>
        )
    },

    _onChange: function() {
        this.setState({
            photos: PhotoStore.getPhotos()
        })
    }
});

module.exports = Latest;
