var React = require('react');
var Router = require('react-router');
var PhotoStore = require('../stores/PhotoStore');
var Constants = require('../Constants');
var Actions = require('../Actions');

var moment = require('moment');

var Link = Router.Link;

var Tag = React.createClass({
    render: function (){

        return (
            <span>
                <Link to="search" query={{q: this.props.tag}}><span className="label label-md label-default">#{this.props.tag}</span></Link>&nbsp;
            </span>
        );
    }
});

var PhotoTitle = React.createClass({

    handleEdit: function(event) {
        event.preventDefault();
        if (!this.props.user) {
            return;
        }
        Actions.photoEditMode();
    },

    handleSubmit: function(event) {
        event.preventDefault();
        var title = this.refs.title.getDOMNode().value;
        Actions.photoEditDone(this.props.photo.id, title);
    },

    render: function(){
        var photo = this.props.photo;

        if (this.props.editMode) {
            return (
        <form role="form" name="form" onSubmit={this.handleSubmit}>
            <div className="form-group">
                <input ref="title" type="text" className="form-control" required="required" defaultValue={photo.title} />
            </div>
            <button type="submit" className="btn btn-primary">Save</button>
            <button type="cancel" className="btn btn-default" onClick={this.handleEdit}>Cancel</button>
        </form>
            );
        }

        return (
            <h2 className="col-md-9" onClick={this.handleEdit}>{photo.title}</h2>
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
        Actions.photoEditMode();
    },

    handleVoteUp: function() {

    },


    handleVoteDown: function() {

    },

    render: function(){

        var buttons = [];
        var photo = this.props.photo;

        if (!photo || !photo.perms) {
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

    fetchData: function() {
        Actions.getPhotoDetail(this.getParams().id);
    },

    componentWillReceiveProps: function() {
        this.fetchData();
    },

    componentDidMount: function() {
        this.fetchData();
    },

    componentWillMount: function() {
        PhotoStore.addChangeListener(this._onChange);
    },

    componentWillUnmount: function () {
        PhotoStore.removeChangeListener(this._onChange);
    },

    handleEditTitle: function() {
        if (this.state.photo && this.state.photo.perms.edit){
            Actions.photoEditMode();
        }
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
            <PhotoTitle photo={photo} user={user} editMode={this.state.editMode} />
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
            <a target="_blank" className="thumbnail" title={photo.title} href={'uploads/' + photo.photo}>
                <img alt={photo.title} src={'uploads/thumbnails/' + photo.photo} />
            </a>
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
                    <Link to="user" params={{id: photo.ownerId}} query={{name: photo.ownerName}}>{photo.ownerName}</Link>
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

        var isDeleted = PhotoStore.isDeleted();

        if (isDeleted){
            this.transitionTo("popular");
            return;
        }

        var photo = PhotoStore.getPhotoDetail();
        if (photo) {
            this.setState({
                photo: photo,
                editMode: PhotoStore.isEditMode()
            });
        } else {
            // TBD: transition to a 404 page
            this.transitionTo("popular");
        }
    }

});

module.exports = PhotoDetail;
