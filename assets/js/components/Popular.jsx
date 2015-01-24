var React = require('react');
var PhotoList = require('./PhotoList.jsx')

var Popular = React.createClass({

    componentWillMount: function() {
        console.log("componentWillMount")
    },
    render: function() {
        return (
            <PhotoList photos={this.props.data.photos.photos} />
        )
    }
});
module.exports = Popular;
