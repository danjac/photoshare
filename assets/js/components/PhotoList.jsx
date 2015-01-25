var React = require('react');
var Router = require('react-router');

var PhotoListItem = React.createClass({

    mixins: [Router.Navigation],

    handleSelectPhoto: function() {
        this.transitionTo("photoDetail", { id: this.props.photo.id });
    },

    render: function(){
        var photo = this.props.photo;
        return (
    <div className="col-xs-6 col-md-3" onClick={this.handleSelectPhoto}>
        <div className="thumbnail">
            <img alt={photo.title} className="img-responsive" src={'uploads/thumbnails/' + photo.photo} />
            <div className="caption">
                <h3>{photo.title}</h3>
            </div>
        </div>
    </div>
        );
    }
});

var PhotoList = React.createClass({
    render: function (){

        return (
            <div className="row">
                {this.props.photos.map(function(photo){
                    return <PhotoListItem key={photo.id} photo={photo} />;
                })};
            </div>
        );
    }

})

module.exports = PhotoList;
