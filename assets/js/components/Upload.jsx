var React = require('react');

var Auth = require('./Auth.jsx');


var Upload = React.createClass({

    mixins: [Auth],

    handleUpload: function() {

    },

    render: function(){

        return (
        <div>
            <div className="spinner"></div>

            <div className="row">
                <div className="col-md-6">
                    <form name="form" role="form" enctype="multipart/form-data" onSubmit={this.handleUpload}>

                        <div className="form-group">
                            <label for="">Title</label>
                            <input type="text" required="required" className="form-control" placeholder="Title" />
                            <span className="help-block"></span>
                        </div>

                        <div className="form-group">
                            <label for="">Tags</label>
                            <input type="text" className="form-control" placeholder="Tags (separate with spaces)" />
                        </div>

                        <div className="form-group">
                            <label for="photo">Photo</label>
                            <input type="file" id="photo" className="form-control" />
                        </div>

                        <button className="btn-submit btn" type="submit">Upload</button>

                        <button className="btn-submit btn" type="submit">Upload &amp; add another</button>
                    </form>
                </div>
                <div className="col-md-6">
                    <div className="thumbnail">
                        <img src="" />
                    </div>
                </div>
            </div>
        </div>
        );
    }
});

module.exports = Upload;
