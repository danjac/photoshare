/* jslint ignore:start */
import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import moment from 'moment';

import {
  Button,
  ButtonInput,
  Input,
  Label
} from 'react-bootstrap';

import * as ActionCreators from '../actions';
import { Facon, Loader } from './util';

@connect(state => {
  return state.photoDetail.toJS();
})
export default class PhotoDetail extends React.Component {

  static propTypes = {
    photo: PropTypes.object.isRequired,
    isLoaded: PropTypes.bool.isRequired,
    isDeleted: PropTypes.bool.isRequired,
    isNotFound: PropTypes.bool.isRequired,
    isEditingTitle: PropTypes.bool.isRequired,
    isEditingTags: PropTypes.bool.isRequired,
    dispatch: PropTypes.func.isRequired
  }

  static contextTypes = {
    router: PropTypes.object.isRequired
  }

  constructor(props) {
    super(props);
    const {dispatch} = this.props;

    this.actions = bindActionCreators(ActionCreators.photoDetail, dispatch);

    this.handleVoteUp = this.handleVoteUp.bind(this);
    this.handleVoteDown = this.handleVoteDown.bind(this);
    this.handleDelete = this.handleDelete.bind(this);
    this.handleToggleEditTitle = this.handleToggleEditTitle.bind(this);
    this.handleToggleEditTags = this.handleToggleEditTags.bind(this);
    this.handleUpdateTitle = this.handleUpdateTitle.bind(this);
    this.handleUpdateTags = this.handleUpdateTags.bind(this);

  }

  componentDidMount() {
    this.actions.getPhotoDetail(this.props.params.id);
  }

  handleUpdateTitle(event) {
    event.preventDefault();
    const title = this.refs.title.getValue().trim();
    if (title) {
      this.actions.updateTitle(this.props.photo.id, title);
    }
  }

  handleUpdateTags(event) {
    event.preventDefault();
    const tags = this.refs.tags.getValue().trim().split(" ");
    this.actions.updateTags(this.props.photo.id, tags);
  }

  handleToggleEditTitle(event) {
    event.preventDefault();
    this.actions.toggleEditTitle();
  }

  handleToggleEditTags(event) {
    event.preventDefault();
    if (this.props.photo.perms.edit) {
      this.actions.toggleEditTags();
    }
  }

  handleVoteUp(event) {
    event.preventDefault();
    this.actions.voteUp(this.props.photo.id);
  }

  handleVoteDown(event) {
    event.preventDefault();
    this.actions.voteDown(this.props.photo.id);
  }

  handleDelete(event) {
    event.preventDefault();
    this.actions.deletePhoto(this.props.photo.id);
    this.context.router.goBack();
  }

  componentDidUpdate(prevProps) {
    if (!prevProps.isEditingTitle && this.props.isEditingTitle) {
      this.refs.title.getInputDOMNode().select();
    }
    if (!prevProps.isEditingTags && this.props.isEditingTags) {
      this.refs.tags.getInputDOMNode().focus();
    }
    return prevProps !== this.props;
  }

  shouldComponentUpdate(nextProps) {
    if (nextProps.isDeleted && !this.props.isDeleted) {
      this.context.router.transitionTo("/");
    }
    return nextProps !== this.props;
  }

  renderButtons() {

    const buttons = [];

    if (this.props.photo.perms.edit) {
      buttons.push(<Button key="edit" onClick={this.handleToggleEditTitle}><Facon name="pencil" /></Button>);
    }

    if (this.props.photo.perms.vote) {
      buttons.push(<Button key="voteUp" onClick={this.handleVoteUp}><Facon name="thumbs-up" /></Button>);
      buttons.push(<Button key="voteDown" onClick={this.handleVoteDown}><Facon name="thumbs-down" /></Button>);
    }

    if (this.props.photo.perms.delete) {
      buttons.push(<Button key="delete" bsStyle="danger" onClick={this.handleDelete}><Facon name="trash" /></Button>);
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

  renderTitle() {
    if (this.props.isEditingTitle) {
      return (
        <div>
          <form className="form-inline" onSubmit={this.handleUpdateTitle}>
            <Input ref="title" type="text" defaultValue={this.props.photo.title} />
            <Button type="submit"><Facon name="floppy-o" /></Button>
          </form>
        </div>
      );
    }
    return <h3 onClick={this.handleToggleEditTitle}>{this.props.photo.title}</h3>;
  }

  renderTags() {
    const tags = this.props.photo.tags || [];
    if (this.props.isEditingTags) {
      return (
          <div>
            <form className="form-inline" onSubmit={this.handleUpdateTags}>
              <Input ref="tags" type="text" defaultValue={tags.join(" ")} />
              <Button type="submit"><Facon name="floppy-o" /></Button>
            </form>
          </div>
      );
    }
    if (!tags) {
      return "";
    }
    return (
      <div>
        {tags.map(tag => {
        return <span key={tag}><Link to='/search/' query={{q: tag}}><Label default>#{tag}</Label></Link>&nbsp;</span>;
        })}
        {this.props.photo.perms.edit ? <Label bsStyle="primary" onClick={this.handleToggleEditTags}>Edit tags <Facon name="pencil" /></Label> : ''}
      </div>
    );

  }

  render() {

    const makeHref = this.context.router.makeHref;

    const photo = this.props.photo;
    const src = photo.photo ? `/uploads/thumbnails/${photo.photo}` : '/img/ajax-loader.gif';

    if (!this.props.isLoaded) {
      return <Loader />;
    }

    if (this.props.isNotFound) {
      return <div>Photo not found</div>;
    }

    return (
    <div>
      {this.renderTitle()}
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
                      <Facon name="thumbs-up" /> {photo.upVotes}
                      &nbsp;
                      <Facon name="thumbs-down" /> {photo.downVotes}
                  </dd>
                  <dt>Uploaded by</dt>
                  <dd>
                      <a href={makeHref(`/user/${photo.ownerId}/${photo.ownerName}`)}>{photo.ownerName}</a>
                  </dd>
                  <dt>Uploaded on</dt>
                  <dd>{moment(photo.createdAt).format('MMMM Do YYYY h:mm')}</dd>
              </dl>
              {this.renderTags()}
          </div>
      </div>
    </div>
    );

  }
}
