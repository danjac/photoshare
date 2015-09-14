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
    this.handleSearch = this.handleSearch.bind(this);
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
  let { source, filter, orderBy } = state.tags.toJS();
  if (filter) {
    source = source.filter(tag => tag.name.match(filter))
  }
  source.sort((left, right) => {
    const up = orderBy === 'numPhotos' ? -1 : 1;
    const down = orderBy === 'numPhotos' ? 1 : -1;
    return left[orderBy] > right[orderBy] ? up : ( left[orderBy] === right[orderBy] ? 0 : down );
  });
  return {
    tags: source,
    filter: filter,
    orderBy: orderBy
  };
})
export default class TagList extends React.Component {

  static propTypes = {
    tags: PropTypes.array.isRequired,
    dispatch: PropTypes.func.isRequired
  }

  static contextTypes = {
    router: PropTypes.object.isRequired
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

  handleSearch(event) {
    event.preventDefault();
    event.stopPropagation();

    const filter = this.refs.filterStr.getValue().trim();

    if (filter) {
      this.context.router.transitionTo("/search/", { q: filter });
    }

  }

  componentDidMount() {
    this.actions.getTags();
  }

  shouldComponentUpdate(nextProps) {
    if (nextProps.tags.length !== this.props.tags.length && nextProps.tags.length === 1) {
      this.context.router.transitionTo("/search/", { q: nextProps.tags[0].name });
      return true;
    }
    return this.props !== nextProps;
  }

  render() {
    return (
      <div>
          <div className="tag-control-box">
            <form className="form-inline" onSubmit={this.handleSearch.bind(this)}>
                <div className="form-group">
                    <Input ref="filterStr" type="text" placeholder="Find a tag" onChange={this.handleFilter.bind(this)} />
                    <Button bsStyle={this.props.orderBy === 'numPhotos' ? 'primary': 'default'}
                            onClick={this.handleOrderBy.bind(this, "numPhotos")}><i className="fa fa-sort-numeric-desc"></i></Button>
                    <Button bsStyle={this.props.orderBy === 'name' ? 'primary': 'default'}
                            onClick={this.handleOrderBy.bind(this, "name")}><i className="fa fa-sort-alpha-desc"></i></Button>
                </div>
            </form>
          </div>

          {this.props.tags.map(tag => <Tag tag={tag} />)}
      </div>
    );
  }

}
