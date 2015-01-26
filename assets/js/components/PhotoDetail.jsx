var React = require('react');
var Router = require('react-router');
var PhotoStore = require('../stores/PhotoStore');
var Actions = require('../Actions');
var moment = require('moment');

var Tag = React.createClass({
    render: function (){
       return (
            <span>
                <a href="#"><span className="label label-md label-default">#{this.props.tag}</span></a>&nbsp;
            </span>
        );
    }
});

var PhotoDetailToolbar = React.createClass({

    render: function(){
        if (!this.props.user) {
            return (<div />);
        }
        return (
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
        );
    }

});

var PhotoDetail = React.createClass({

    mixins: [Router.State],

    getInitialState: function() {
        return {
            photo: null
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
        var user = this.props.user;

        if (photo === null) {
            return (<div />);
        }

        var tags = photo.tags || [];

        return (
    <div>
        <div className="row">
            <h2 className="col-md-9">{photo.title}</h2>
            <PhotoDetailToolbar user={user} />
        </div>

        <hr />
    <h3>
        {tags.map(function(tag){
            return <Tag key={tag} tag={tag} />;
        })}
    </h3>

    <div className="row">
        <div className="col-xs-6 col-md-3">
            <a target="_blank" className="thumbnail" title="{photo.title}" href={'uploads/' + photo.photo}>
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
                </dd>
                <dt>Uploaded on</dt>
                <dd>{moment(photo.createdAt).format("MMMM Do YYYY hh:mm")}</dd>
            </dl>

        </div>
    </div>
</div>
        );
    },

    _onChange: function() {
        this.setState({
            photo: PhotoStore.getPhotoDetail()
        });
    }

});

module.exports = PhotoDetail;
