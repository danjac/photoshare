import React, { PropTypes } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import {
  Input,
  Button
} from 'react-bootstrap';

import * as ActionCreators from '../actions';

class Tag extends React.Component {
  static propTypes = {
    tag: PropTypes.object.isRequired
  }

  static contextTypes = {
    router: PropTypes.object.isRequired
  }

  constructor(props) {
    super(props);
    this.handleClick = this.handleClick.bind(this);
  }

  handleSearch(event) {
    event.preventDefault();
    this.context.router.transitionTo('/search/', { q: this.props.tag.name });
  }

  render() {
    return (
        <div className="col-xs-6 col-md-3 ">
            <div className="thumbnail " onClick={this.handleSearch}>
                <img alt={this.props.tag.name} class="img-responsive " src={`uploads/thumbnails/${this.props.tag.photo}`} />
                <div className="caption ">
                    <h3>#{this.props.tag.name}</h3>
                </div>
            </div>
        </div>
    );

  }

}

@connect(state => {
  let tags = state.tags.toJS().source;
  if (state.tags.filter) {
    tags = tags.filter(tag => tag.name.match(state.tags.filter))
  }
  tags = tags.sort((left, right) => {
    const orderBy = state.tags.orderBy;
    return left[orderBy] > right[orderBy] ? 1 : ( left[orderBy] === right[orderBy] ? 0 : -1 );
  });
  return {
    tags: tags
  };
})
export default class TagList extends React.Component {

  static propTypes = {
    tags: PropTypes.array.isRequired,
    dispatch: PropTypes.func.isRequired
  }

  constructor(props) {
    super(props);
    const { dispatch } = this.props;
    this.actions = bindActionCreators(ActionCreators.tags, dispatch);
  }

  handleOrderBy(orderBy, event) {
    event.preventDefault();
    this.actions.orderTags(orderBy);
  }

  handleFilter(event) {
    event.preventDefault();
    const filter = this.refs.filterStr.getValue();
    this.actions.filterTags(filter);
  }

  componentDidMount() {
    this.actions.getTags();
  }

  render() {
    return (
      <div>
          <div className="tag-control-box">
            <form className="form-inline">
                <div className="form-group">
                    <Input ref="filterStr" type="text" placeholder="Find a tag" onChange={this.handleFilter.bind(this)} />
                    <Button onClick={this.handleOrderBy.bind(this, "numPhotos")}><i className="fa fa-sort-numeric-desc"></i></Button>
                    <Button onClick={this.handleOrderBy.bind(this, "name")}><i className="fa fa-sort-alpha-desc"></i></Button>
                </div>
            </form>
          </div>

          {this.props.tags.map(tag => <Tag tag={tag} />)}
      </div>
    );
  }

}
