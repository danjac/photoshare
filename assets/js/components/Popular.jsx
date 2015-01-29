var React = require('react');
var Actions = require('../Actions');
var PhotoStore = require('../stores/PhotoStore');
var PhotoList = require('./PhotoList.jsx');

var Popular = React.createClass({

    getInitialState: function() {
        return {
            photos: this.props.photos || {
                photos: []
            }
        }
    },

    componentWillMount: function() {
        PhotoStore.addChangeListener(this._onChange);
    },

    componentDidMount: function() {
        Actions.getPhotos("votes");
    },

    componentWillUnmount: function() {
        PhotoStore.removeChangeListener(this._onChange);
    },

    handlePaginationLink: function(page) {
        Actions.getPhotos("votes", page);
    },

    render: function() {
        console.log(this.state.photos)
        return (
            <PhotoList photos={this.state.photos} handlePaginationLink={this.handlePaginationLink} />
        )
    },

    _onChange: function() {
        this.setState({
            photos: PhotoStore.getPhotos()
        })
    }
});
module.exports = Popular;
