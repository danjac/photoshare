/* jslint ignore:start */
import React, { PropTypes } from 'react';
import _ from 'lodash';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Input,
         ButtonInput,
        } from 'react-bootstrap';

import * as ActionCreators from '../actions';
import { Loader } from './widgets';


@connect(state => {
  let props = state.upload.toJS();
  let forms = state.forms.toJS();
  props.errors = forms.upload ? forms.upload.errors : {};
  props.checked = forms.upload ? forms.upload.checked: [];
  return props;
})
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

    this.actions.previewPhoto(files[0]);
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

    if(!_.isEmpty(this.props.errors)) {
      return;
    }

    const title = this.refs.title.getValue().trim(),
          tags = this.refs.tags.getValue().trim(),
          photo = this.refs.photo.getInputDOMNode().files[0];

    this.actions.upload(title, tags, photo);
  }

  componentDidMount() {
    this.actions.resetForm();
  }

  shouldComponentUpdate(nextProps) {
    if (nextProps.uploadedPhoto) {
      const { id } = nextProps.uploadedPhoto;
      this.context.router.transitionTo("/detail/" + id);
    }
    return nextProps !== this.props;
  }

  handleCheckTitle(event) {
    event.preventDefault();
    const title = this.refs.title.getValue().trim();
    this.actions.checkTitle(title);
  }

  handleCheckPhoto(event) {
    event.preventDefault();
    const photo = this.refs.photo.getInputDOMNode().files[0];
    this.actions.checkPhoto(photo);
  }

  errorStatus(field) {
    if (this.props.checked.indexOf(field) === -1) {
        return;
    }
    return this.props.errors[field] ? 'error' : 'success';
  }

  errorMsg(field) {
    return this.props.errors[field] || '';
  }

  render() {

    if (this.props.isWaiting) {
      return <Loader />;
    }

    const handlePhotoSelect = this.handlePhotoSelect.bind(this);
    const handleCheckTitle = this.handleCheckTitle.bind(this);
    const handleCheckPhoto = this.handleCheckPhoto.bind(this);
    const handleSubmit = this.handleSubmit.bind(this);

    return (
      <div className="row">
          <div className="col-md-6">
              <form name="form" role="form" encType="multipart/form-data" onSubmit={handleSubmit}>

                <Input name="title"
                       type="text"
                       ref="title"
                       label="Title"
                       hasFeedback
                       onBlur={handleCheckTitle}
                       bsStyle={this.errorStatus('title')}
                       help={this.errorMsg('title')} />

                <Input name="tags"
                       type="text"
                       ref="tags"
                       label="Tags"
                       placeholder="Separate with spaces" />

                <Input name="photo"
                       type="file"
                       onChange={handleCheckPhoto}
                       ref="photo" label="Photo"
                       hasFeedback
                       bsStyle={this.errorStatus('photo')}
                       help={this.errorMsg('photo')} />

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
