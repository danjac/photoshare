/* jslint ignore:start */
import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import moment from 'moment';

import {
  Button
} from 'react-bootstrap';

import * as ActionCreators from '../actions';
import { Facon } from './util';

@connect(state => {
  return state.photoDetail.toJS();
})
export default class PhotoDetail extends React.Component {

  static propTypes = {
    photo: PropTypes.object.isRequired,
    isEditing: PropTypes.bool.isRequired,
    dispatch: PropTypes.func.isRequired
  }

  static contextTypes = {
    router: PropTypes.object.isRequired
  }

  constructor(props) {
    super(props);
    const {dispatch} = this.props;
    this.actions = bindActionCreators(ActionCreators.photoDetail, dispatch);

    this.handleToggleEdit = this.handleToggleEdit.bind(this);
    this.handleVoteUp = this.handleVoteUp.bind(this);
    this.handleVoteDown = this.handleVoteDown.bind(this);
    this.handleDelete = this.handleDelete.bind(this);

  }

  componentDidMount() {
    this.actions.getPhotoDetail(this.props.params.id);
  }

  handleToggleEdit(event) {
    event.preventDefault();
  }

  handleVoteUp(event) {
    event.preventDefault();
  }

  handleVoteDown(event) {
    event.preventDefault();
  }

  handleDelete(event) {
    event.preventDefault();
    this.actions.deletePhoto(this.props.photo);
    this.context.router.goBack();
  }

  renderButtons() {

    if (!this.props.photo) {
      return '';
    }

    const buttons = [];

    if (this.props.photo.perms.edit) {
      buttons.push(<Button onClick={this.handleToggleEdit}><Facon name="pencil" /></Button>);
    }

    if (this.props.photo.perms.vote) {
      buttons.push(<Button onClick={this.handleVoteUp}><Facon name="thumbsUp" /></Button>);
      buttons.push(<Button onClick={this.handleVoteDown}><Facon name="thumbsDown" /></Button>);
    }

    if (this.props.photo.perms.delete) {
      buttons.push(<Button bsStyle="danger" onClick={this.handleDelete}><Facon name="trash" /></Button>);
    }

    if (!buttons) {
      return '';
    }

    return (
      <div className="button-group col-md-3 pull-right">
      {buttons}
      </div>
    );
  }

  render() {
    const photo = this.props.photo;
    const src = photo.photo ? `/uploads/thumbnails/${photo.photo}` : '/img/ajax-loader.gif';

    return (
    <div>
      <h3>{photo.title}</h3>
      {this.renderButtons()}
      <div className="row">
          <div className="col-xs-6 col-md-3">
              <a target="_blank" className="thumbnail" title={photo.title} href={`/uploads/${photo.photo}`}>
                  <img alt={photo.title} src={`/uploads/thumbnails/${photo.photo}`} />
              </a>
          </div>
          <div className="col-xs-6">
              <dl>
                  <dt>Score <span className="badge">{photo.score}</span></dt>
                  <dd>
                      <i className="fa fa-thumbs-up"></i> {photo.upVotes}
                      &nbsp;
                      <i className="fa fa-thumbs-down"></i> {photo.downVotes}
                  </dd>
                  <dt>Uploaded by</dt>
                  <dd>
                      <a href="#">{photo.ownerName}</a>
                  </dd>	<dt>Uploaded on</dt>
                  <dd>{moment(photo.createdAt).format('MMMM Do YYYY h:mm')}</dd>
              </dl>

          </div>
      </div>
    </div>
    );

  }
}
