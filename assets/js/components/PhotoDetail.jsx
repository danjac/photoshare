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

    handleDelete: function() {
        if (window.confirm("Are you sure you want to delete this photo?")){
            Actions.deletePhoto(this.props.photo.id);
        }
    },

    handleEdit: function () {

    },

    handleVoteUp: function() {

    },


    handleVoteDown: function() {

    },

    render: function(){

        var buttons = [];
        var photo = this.props.photo;
        if (!photo) {
            return <div />;
        }

        if (photo.perms.edit) {
            buttons.push({
                className: "btn-default",
                icon: "fa-pencil",
                onClick: this.handleEdit
            });
        }

        if (photo.perms.vote) {
            buttons.push({
                className: "btn-default",
                icon: "fa-thumbs-up",
                onClick: this.handleVoteUp
            });
            buttons.push({
                className: "btn-default",
                icon: "fa-thumbs-down",
                onClick: this.handleVoteDown
            });
        }

        if (photo.perms.delete) {
            buttons.push({
                className: "btn-danger",
                icon: "fa-trash",
                onClick: this.handleDelete
            });
        }

        return (
            <div className="button-group col-md-3 pull-right">
                {buttons.map(function(btn) {
                    var className = "btn " + btn.className;
                    var icon = "fa " + btn.icon;
                    return (
                <button onClick={btn.onClick} type="button" className={className}> <i className={icon}></i></button>
                    )
                })}
            </div>
        );
    }

});

var PhotoDetail = React.createClass({

    mixins: [Router.State, Router.Navigation],

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
            <PhotoDetailToolbar photo={photo} user={user} />
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
                <img alt={photo.title} src={'uploads/thumbnails/' + photo.photo} />
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
        var photo = PhotoStore.getPhotoDetail();
        if (photo) {
            this.setState({
                photo: photo
            });
        } else {
            this.transitionTo("popular");
        }
    }

});

module.exports = PhotoDetail;
