var React = require('react');
var API = require('../API.js');
var PhotoList = require('./PhotoList.jsx')

var Latest = React.createClass({
    getInitialState: function() {
        return {
            photos: []
        }
    },
    componentWillMount: function() {
        var self = this;
        API.getPhotos(null, function(data){
            self.setState({
                photos: data.photos
            })
        });
    },

    render: function() {
        return (
            <PhotoList photos={this.state.photos} />
        )
    }
});
module.exports = Latest;
