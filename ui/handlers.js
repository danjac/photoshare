import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { Provider, connect } from 'react-redux';
import { Navbar, Nav, Input, Pagination, Button } from 'react-bootstrap';
import { NavItemLink } from 'react-router-bootstrap';
import { RouteHandler, Link } from 'react-router';
import moment from 'moment';

import * as ActionCreators from './actions';
import configureStore from './store';

const store = configureStore();


class Navigation extends React.Component {

  render() {

    const brand = <Link to="app"><i className="fa fa-camera"></i> Wallshare</Link>;
    const searchIcon = <i className="fa fa-search"></i>;

    return (

      <Navbar fixedTop inverse brand={brand}>

        <Nav>
          <NavItemLink to="app"><i className="fa fa-fire"></i> Popular</NavItemLink>
          <NavItemLink to="latest"><i className="fa fa-clock-o"></i> Latest</NavItemLink>
          <NavItemLink to="app"><i className="fa fa-tags"></i> Tags</NavItemLink>
          <NavItemLink to="app"><i className="fa fa-upload"></i> Upload</NavItemLink>
        </Nav>
        <Nav>
          <form className="navbar-form navbar-left" role="search" name="searchForm">
            <Input type="text" addonAfter={searchIcon} bsSize="small" placeholder="Search" />
          </form>
        </Nav>
        <Nav right>
          <NavItemLink to="login"><i className="fa fa-sign-in"></i> Login</NavItemLink>
          <NavItemLink to="app"><i className="fa fa-user"></i> Signup</NavItemLink>
        </Nav>
      </Navbar>
    );
  }


}

export class App extends React.Component {
  render() {
    return (
      <Provider store={store}>
      {() => {
        return (
        <div className="container-fluid">
          <Navigation />
          <RouteHandler />
        </div>
        );
      }}
      </Provider>
    );
  }
}


export class Login extends React.Component {
  render() {
    return (
      <div>
        <form role="form" name="form">
            <Input type="text" ref="identifier" required placeholder="Name or email address" />
            <Input type="password" ref="password" required placeholder="Name or email address" />
            <Button bsStyle="primary">Login</Button>
        </form>

        <a href="#/recoverpass">Forgot your password?</a> 
      </div>
    );
  }
}

class PhotoListItem extends React.Component {
  static propTypes = {
    photo: PropTypes.object.isRequired
  }

  static contextTypes = {
    router: PropTypes.func.isRequired
  }

  constructor(props) {
    super(props);
    this.handleClick = this.handleClick.bind(this);
  }

  handleClick(event) {
    event.preventDefault();
    this.context.router.transitionTo('detail', {id: this.props.photo.id});
  }

  render() {

    const photo = this.props.photo;

    return (
      <div className="col-xs-6 col-md-3">
          <div className="thumbnail" onClick={this.handleClick}>
              <img alt={photo.title} className="img-responsive" src={`/uploads/thumbnails/${photo.photo}`} />
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

@connect(state => {
  return {
    photo: state.photoDetail.toJS()
  };
})
export class PhotoDetail extends React.Component {

  static propTypes = {
    photo: PropTypes.object.isRequired,
    dispatch: PropTypes.func.isRequired
  }

  constructor(props) {
    super(props);
    const {dispatch} = this.props; 
    this.actions = bindActionCreators(ActionCreators.photoDetail, dispatch);
  }

  componentDidMount() {
    this.actions.getPhotoDetail(this.props.params.id);
  }

  render() {
    const photo = this.props.photo;
    return (
    <div>
      <h3>{photo.title}</h3>
      <div className="row">
          <div className="col-xs-6 col-md-3">
              <a target="_blank" className="thumbnail" title={photo.title} href={`/uploads/${photo.photo}`}>
                  <img alt={photo.title} src={`/uploads/thumbnails/${photo.photo}`} />
              </a>
              <div class="btn-group">
              </div>
          </div>
          <div class="col-xs-6">
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
