import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Pagination } from 'react-bootstrap';

import * as ActionCreators from './actions';


class PhotoListItem extends React.Component {
  static propTypes = {
    photo: PropTypes.object.isRequired
  }

  static contextTypes = {
    router: PropTypes.object.isRequired
  }

  constructor(props) {
    super(props);
    this.handleClick = this.handleClick.bind(this);
  }

  handleClick(event) {
    event.preventDefault();
    this.context.router.transitionTo('/detail/' + this.props.photo.id);
  }

  render() {

    const photo = this.props.photo;
    const src = photo.photo ? `/uploads/thumbnails/${photo.photo}` : '/img/ajax-loader.gif';

    return (
      <div className="col-xs-6 col-md-3">
          <div className="thumbnail" onClick={this.handleClick}>
              <img alt={photo.title} className="img-responsive" src={src} />
              <div className="caption">
                  <h3>{photo.title.substring(0, 20)}</h3>
              </div>
          </div>
      </div>
      );
  } 
}

class PhotoList extends React.Component {

  static propTypes = {
    photos: PropTypes.array.isRequired,
    total: PropTypes.number.isRequired,
    numPages: PropTypes.number.isRequired,
    currentPage: PropTypes.number.isRequired,
    handlePageClick: PropTypes.func.isRequired
  }

  render() {
    const { handlePageClick, total, numPages, currentPage, photos } = this.props;
    const pagination = (
      <Pagination onSelect={handlePageClick}
                  items={numPages}
                  ellipsis={true}
                  first={true}
                  last={true}
                  next={true}
                  prev={true}
                  maxButtons={12}
                  activePage={currentPage} />
    );
    return (
    <div>
      {pagination}
      <div className="row">
          {photos.map(photo => {

            return <PhotoListItem key={photo.id} photo={photo} />
          })}
      </div>
      {pagination}
    </div>
    );
  }

}


@connect(state => {
  return {
    photos: state.photos.toJS()
  }
})
export class Popular extends React.Component {

  static propTypes = {
    photos: PropTypes.object.isRequired,
    dispatch: PropTypes.func.isRequired
  }

  constructor(props) {
    super(props);
    const {dispatch} = this.props;
    this.actions = bindActionCreators(ActionCreators.photos, dispatch);
    this.handlePageClick = this.handlePageClick.bind(this);
  }

  handlePageClick(event, selectedEvent) {
    event.preventDefault();
    const page = selectedEvent.eventKey;
    this.actions.getPhotos(page, "votes");
  }

  componentDidMount() {
    this.actions.getPhotos(1, "votes");
  }

  render() {
    return <PhotoList handlePageClick={this.handlePageClick} {...this.props.photos} />;
  }

}

@connect(state => {
  return {
    photos: state.photos.toJS()
  };
})
export class Latest extends React.Component {

  static propTypes = {
    photos: PropTypes.object.isRequired,
    dispatch: PropTypes.func.isRequired
  }

  constructor(props) {
    super(props);
    const {dispatch} = this.props;
    this.actions = bindActionCreators(ActionCreators.photos, dispatch);
    this.handlePageClick = this.handlePageClick.bind(this);
  }

  handlePageClick(event, selectedEvent) {
    event.preventDefault();
    const page = selectedEvent.eventKey;
    this.actions.getPhotos(page, "created");
  }

  componentDidMount() {
    this.actions.getPhotos(1, "created");
  }

  render() {
    return <PhotoList handlePageClick={this.handlePageClick} {...this.props.photos} />;
  }

}