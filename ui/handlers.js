import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { Provider, connect } from 'react-redux';
import { Navbar, Nav, Input, Pagination } from 'react-bootstrap';
import { NavItemLink } from 'react-router-bootstrap';
import { RouteHandler, Link } from 'react-router';

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
          <NavItemLink to="app"><i className="fa fa-sign-in"></i> Login</NavItemLink>
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
            return (
            <div className="col-xs-6 col-md-3">
                <div className="thumbnail">
                    <img alt={photo.title} className="img-responsive" src={`uploads/thumbnails/${photo.photo}`} />
                    <div className="caption">
                        <h3>{photo.title.substring(0, 20)}</h3>
                    </div>
                </div>
            </div>
            );
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
  }
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


