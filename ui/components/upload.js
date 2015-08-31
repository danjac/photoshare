import React, { PropTypes } from 'react';
import _ from 'lodash';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Input, 
         ButtonInput 
        } from 'react-bootstrap';

import * as ActionCreators from '../actions';


@connect(state => state.upload.toJS())
export default class Upload extends React.Component {

  static propTypes = {
    previewURL: PropTypes.string
  }

  static contextTypes = {
    router: PropTypes.object.isRequired
  }

  constructor(props) {
    super(props);
    this.actions = bindActionCreators(ActionCreators.upload, this.props.dispatch);
  }
  
  handlePhotoSelect(event) {

    event.preventDefault();
    const files = event.target.files;

    if (!files) {
      return;
    }

    const file = files[0];
    if (!file.type.match('image.*')) {
      return;
    }

    this.actions.previewPhoto(file);
  }

  previewPhoto() {
    if (this.props.previewURL) {
      return (
       <div className="thumbnail">
            <img src={this.props.previewURL} />
        </div>
      );
    }
    return '';
  }

  handleSubmit(event) {
    event.preventDefault();
    const title = this.refs.title.getValue().trim(),
          tags = this.refs.tags.getValue().trim(),
          photo = this.refs.photo.getInputDOMNode().files[0];

    if (!_.every([title, tags, photo])) {
      return;
    }

    this.refs.title.getInputDOMNode().value = "";
    this.refs.tags.getInputDOMNode().value = "";

    this.actions.upload(title, tags, photo);

  }

  shouldComponentUpdate(nextProps) {
    if (nextProps.uploadedPhoto) {
      const { id } = nextProps.uploadedPhoto;
      this.actions.reset();
      this.context.router.transitionTo("/detail/" + id);
      return true;
    }
    return nextProps !== this.props;
  }

  render() {
    const handlePhotoSelect = this.handlePhotoSelect.bind(this);
    const handleSubmit = this.handleSubmit.bind(this);

    return (
      <div className="row">
          <div className="col-md-6">
              <form name="form" role="form" encType="multipart/form-data" onSubmit={handleSubmit}>
                <Input name="title" type="text" ref="title" label="Title" />
                <Input name="tags" type="text" ref="tags" label="Tags" placeholder="Separate with spaces" />
                <Input name="photo" type="file" onChange={handlePhotoSelect} ref="photo" label="Photo" />
                <ButtonInput type="submit" bsStyle="primary">Upload</ButtonInput>
              </form>
          </div>
          <div className="col-md-6">
              {this.previewPhoto()} 
          </div>
      </div>
    );
  }
}
