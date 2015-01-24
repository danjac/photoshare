var React = require('react');

var Photo = React.createClass({

    render: function(){
        var photo = this.props.photo;
        return (
    <div className="col-xs-6 col-md-3">
        <div className="thumbnail">
            <img alt={photo.title} className="img-responsive" src={'uploads/thumbnails/' + photo.photo} />
            <div class="caption">
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
                {this.props.photos.map(function(photo) {
                    return <Photo photo={photo} />;
                })};
            </div>
        );
    }

})

module.exports = PhotoList;
