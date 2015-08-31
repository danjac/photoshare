import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Link } from 'react-router';
import moment from 'moment';

import * as ActionCreators from '../actions';

@connect(state => {
  return {
    photo: state.photoDetail.toJS()
  };
})
export default class PhotoDetail extends React.Component {

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
