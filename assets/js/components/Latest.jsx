var React = require('react');
var Actions = require('../Actions');
var PhotoStore = require('../stores/PhotoStore');
var PhotoList = require('./PhotoList.jsx')

var Latest = React.createClass({

    getInitialState: function() {
        return {
            photos: {
                photos: []
            }
        }
    },

    componentWillMount: function() {
        PhotoStore.addChangeListener(this._onChange);
        Actions.getPhotos();
    },

    componentWillUnmount: function() {
        PhotoStore.removeChangeListener(this._onChange);
    },

    render: function() {
        return (
            <PhotoList photos={this.state.photos.photos} />
        )
    },

    _onChange: function() {
        this.setState({
            photos: PhotoStore.getPhotos()
        })
    }
});

module.exports = Latest;
