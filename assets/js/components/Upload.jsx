var React = require('react');
var Router = require('react-router');

var Actions = require('../Actions');
var PhotoStore = require('../stores/PhotoStore');

var Auth = require('./Auth.jsx');

var Upload = React.createClass({

    mixins: [Auth, Router.Navigation],

    getInitialState: function() {
        return {
            previewUrl: null
        }
    },

    handleUpload: function(event) {
        event.preventDefault();
        var title = this.refs.title.getDOMNode().value;
        var tags = this.refs.tags.getDOMNode().value;
        var photo = this.refs.photo.getDOMNode().files[0];
        Actions.uploadPhoto(title, tags, photo);
    },

    handlePhotoPreview: function(event) {
        event.preventDefault();
        var photo = this.refs.photo.getDOMNode().files[0];

        if (!photo || window.FileReader === null) {
            return;
        }

        Actions.previewPhoto(photo);

    },

    componentWillMount: function() {
        PhotoStore.addChangeListener(this._onChange);
    },

    componentWillUnmount: function() {
        PhotoStore.removeChangeListener(this._onChange);
    },

    render: function(){

        var preview = "";
        if(this.state.previewUrl) {
            preview = (
                <div className="thumbnail">
                    <img src={this.state.previewUrl} />
                </div>
            );
        }

        return (
        <div>
            <div className="spinner"></div>

            <div className="row">
                <div className="col-md-6">
                    <form name="form" role="form" encType="multipart/form-data" onSubmit={this.handleUpload}>

                        <div className="form-group">
                            <label htmlFor="">Title</label>
                            <input ref="title" type="text" required="required" className="form-control" placeholder="Title" />
                            <span className="help-block"></span>
                        </div>

                        <div className="form-group">
                            <label htmlFor="">Tags</label>
                            <input ref="tags" type="text" className="form-control" placeholder="Tags (separate with spaces)" />
                        </div>

                        <div className="form-group">
                            <label htmlFor="photo">Photo</label>
                            <input ref="photo" type="file" id="photo" className="form-control" onChange={this.handlePhotoPreview} />
                        </div>

                        <button className="btn-submit btn" type="submit">Upload</button>

                        <button className="btn-submit btn" type="submit">Upload &amp; add another</button>
                    </form>
                </div>
                <div className="col-md-6">
                    {preview}
                </div>
            </div>
        </div>
        );
    },

    _onChange: function() {
        var photo = PhotoStore.getNewPhoto();
        if (photo) {
            this.transitionTo('photoDetail', {id: photo.id});
            return;
        }
        this.setState({
            previewUrl: PhotoStore.getPreviewUrl()
        });
    }
});

module.exports = Upload;
