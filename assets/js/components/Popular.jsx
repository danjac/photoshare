var React = require('react');
var PhotoList = require('./PhotoList.jsx')

var Popular = React.createClass({

    render: function() {
        return (
            <PhotoList photos={this.props.data.photos.photos} />
        )
    }
});
module.exports = Popular;
