var React = require('react');
var Router = require('react-router');
var Actions = require('../Actions');
var PhotoStore = require('../stores/PhotoStore');
var PhotoList = require('./PhotoList.jsx');


var Search = React.createClass({

    mixins: [Router.State],

    getInitialState: function() {
        return {
            photos: {
                photos: [],
            }
        }
    },

    getSearch: function(){
        return this.getQuery().q;
    },

    fetchData: function(page) {

        Actions.searchPhotos(this.getSearch(), page);
    },

    componentDidMount: function() {
        this.fetchData();
    },

    componentWillReceiveProps: function() {
        this.fetchData();
    },

    componentWillMount: function() {
        PhotoStore.addChangeListener(this._onChange);
    },

    componentWillUnmount: function() {
        PhotoStore.removeChangeListener(this._onChange);
    },

    handlePaginationLink: function(page) {
        this.fetchData(page);
    },

    render: function() {
        return (
            <div>
            <h3>{this.getSearch()}: {this.state.photos.total}</h3>
            <PhotoList photos={this.state.photos} handlePaginationLink={this.handlePaginationLink} />
            </div>
        )
    },

    _onChange: function() {
        this.setState({
            photos: PhotoStore.getPhotos()
       })
    }
});
module.exports = Search;
