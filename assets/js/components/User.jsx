var React = require('react');
var Router = require('react-router');
var Actions = require('../Actions');
var PhotoStore = require('../stores/PhotoStore');
var PhotoList = require('./PhotoList.jsx')

var User = React.createClass({

    mixins: [Router.State],

    getInitialState: function() {
        return {
            photos: {
                photos: []
            },
            user: {
                id: null,
                name: null
            }
        }
    },

    componentWillMount: function() {
        PhotoStore.addChangeListener(this._onChange);
    },

    componentDidMount: function() {
        this.state.user = {
            id: this.getParams().id,
            name: this.getQuery().name
        }
        Actions.getPhotosForUser(this.state.user.id);
    },

    componentWillUnmount: function() {
        PhotoStore.removeChangeListener(this._onChange);
    },

    handlePaginationLink: function(page) {
        Actions.getPhotosForUser(this.state.user.id, page);
    },

    render: function() {
        return (
            <div>
            <h3>Photos for {this.state.user.name}</h3>
            <PhotoList photos={this.state.photos} handlePaginationLink={this.handlePaginationLink}/>
            </div>
        )
    },

    _onChange: function() {
        this.setState({
            photos: PhotoStore.getPhotos()
        })
    }
});

module.exports = User;
