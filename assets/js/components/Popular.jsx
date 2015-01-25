var React = require('react');
var Actions = require('../Actions');
var PhotoStore = require('../stores/PhotoStore');
var PhotoList = require('./PhotoList.jsx');

var Popular = React.createClass({

    getInitialState: function() {
        return {
            photos: {
                photos: []
            }
        }
    },

    componentWillMount: function() {
        PhotoStore.addChangeListener(this._onChange);
        Actions.getPhotos("votes");
    },

    componentWillUnmount: function() {
        PhotoStore.removeChangeListener(this._onChange);
    },

    componentWillReceiveProps: function(nextProps) {
        if (nextProps) {
            this.setState({
                photos: nextProps
            });
        }
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
module.exports = Popular;
