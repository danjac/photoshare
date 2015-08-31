import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Navbar, 
         Nav, 
         Input, 
         Pagination, 
         ButtonInput, 
         NavDropdown, 
         NavItem, 
         MenuItem } from 'react-bootstrap';
import { Link } from 'react-router';
import moment from 'moment';

import * as ActionCreators from './actions';

@connect(state => state.auth.toJS())
class Navigation extends React.Component {

  static propTypes = {
    loggedIn: PropTypes.bool.isRequired,
    name: PropTypes.string
  }

  static contextTypes = {
    router: PropTypes.func.isRequired
  }

  rightNav() {
    const { name, loggedIn } = this.props;
    const makeHref = this.context.router.makeHref; 

    if (loggedIn) {
      return (
        <Nav right>
          <NavDropdown title={name}>
            <MenuItem>My photos</MenuItem>
            <MenuItem>Change my password</MenuItem>
            <MenuItem>Logout</MenuItem>
          </NavDropdown>
        </Nav>
      );
    }
    return (
      <Nav right>
        <NavItem href={makeHref('/login/')}><i className="fa fa-sign-in"></i> Login</NavItem>
        <NavItem href="/"><i className="fa fa-user"></i> Signup</NavItem>
      </Nav>
    );
  }

  render() {

    const brand = <Link to="/"><i className="fa fa-camera"></i> Wallshare</Link>;
    const searchIcon = <i className="fa fa-search"></i>;
    const makeHref = this.context.router.makeHref; 

    return (

      <Navbar fixedTop inverse brand={brand}>

        <Nav>
          <NavItem href={makeHref("/")}><i className="fa fa-fire"></i> Popular</NavItem>
          <NavItem href={makeHref("/latest/")}><i className="fa fa-clock-o"></i> Latest</NavItem>
          <NavItem href="/"><i className="fa fa-tags"></i> Tags</NavItem>
          <NavItem href="/"><i className="fa fa-upload"></i> Upload</NavItem>
        </Nav>
        <Nav>
          <form className="navbar-form navbar-left" role="search" name="searchForm">
            <Input type="text" addonAfter={searchIcon} bsSize="small" placeholder="Search" />
          </form>
        </Nav>
        {this.rightNav()}
      </Navbar>
    );
  }


}

@connect(state => state.auth.toJS())
export class App extends React.Component {

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(ActionCreators.auth, dispatch);
  }

  componentDidMount() {
    this.actions.getUser();
  }
  
  render() {
    return (
    <div className="container-fluid">
      <Navigation auth={this.props.auth} />
      {this.props.children}
    </div>
    );
  }
}

@connect(state => state.auth.toJS())
export class Login extends React.Component {

  static propTypes = {
    dispatch: PropTypes.func.isRequired
  }

  static contextTypes = {
    router: PropTypes.func.isRequired
  }

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(ActionCreators.auth, dispatch);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  shouldComponentUpdate(nextProps) {
    if (nextProps.loggedIn) {
      this.context.router.transitionTo("/");
      return true;
    }
    return false;
  }

  handleSubmit(event) {
    event.preventDefault();
    const identifier = this.refs.identifier.getValue(),
          password = this.refs.password.getValue();

    this.refs.password.getInputDOMNode().value = "";
    this.actions.login(identifier, password);
  }

  render() {
    return (
      <div className="col-md-6 col-md-offset-3">
        <form role="form" method="POST" onSubmit={this.handleSubmit}>
            <Input type="text" ref="identifier" required placeholder="Name or email address" />
            <Input type="password" ref="password" required placeholder="Password" />
            <ButtonInput bsStyle="primary" type="submit">Login</ButtonInput>
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
    const src = photo.photo ? `/uploads/thumbnails/${photo.photo}` : '/img/ajax-loader.gif';

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
