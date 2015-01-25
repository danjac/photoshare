var React = require('react');
var Router = require('react-router');
var PhotoStore = require('../stores/PhotoStore');
var Actions = require('../Actions');

var PhotoDetail = React.createClass({

    mixins: [Router.State],

    getInitialState: function() {
        return {
            photo: {}
        }
    },

    componentWillMount: function() {
        PhotoStore.addChangeListener(this._onChange);
        Actions.getPhotoDetail(this.getParams().id);
    },

    componentWillUnmount: function () {
        PhotoStore.removeChangeListener(this._onChange);
    },

    render: function() {
        var photo = this.state.photo;

        return (
<div>
    <div className="row">
        <h2 className="col-md-9">{photo.title}</h2>
        <div className="button-group col-md-3 pull-right">
            <button type="button" className="btn btn-default"> <i className="fa fa-pencil"></i>
            </button>
            <button type="button" className="btn btn-default"><i className="fa fa-thumbs-up"></i>
            </button>
            <button type="button" className="btn btn-default"><i className="fa fa-thumbs-down"></i>
            </button>
            <button type="button" className="btn btn-danger"><i className="fa fa-trash"></i>
            </button>
        </div>
    </div>

    <hr />
    <h3>
        tags go here...
    </h3>

    <div className="row">
        <div className="col-xs-6 col-md-3">
            <a target="_blank" className="thumbnail" title="{photo.title}" href="uploads/{photo.photo}">
                <img alt="{photo.title}" src={'uploads/thumbnails/' + photo.photo} />
            </a>
            <div className="btn-group">
            </div>
        </div>
        <div className="col-xs-6">
            <dl>
                <dt>Score <span className="badge">{photo.score}</span></dt>
                <dd>
                    <i className="fa fa-thumbs-up"></i> {photo.upVotes}
                    <i className="fa fa-thumbs-down"></i> {photo.downVotes}
                </dd>
                <dt>Uploaded by</dt>
                <dd>
                    <a href="#">{photo.ownerName}</a>
                </dd>   <dt>Uploaded on</dt>
                <dd>{photo.createdAt}</dd>
            </dl>

        </div>
    </div>
</div>
        );
    },

    _onChange: function() {
        console.log("onchange")
        console.log(PhotoStore.getPhotoDetail())
        this.setState({
            photo: PhotoStore.getPhotoDetail()
        });
    }

});

module.exports = PhotoDetail;
