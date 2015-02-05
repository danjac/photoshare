var React = require('react');
var Router = require('react-router');
var TagStore = require('../stores/TagStore');
var Actions = require('../Actions');

var Tag = React.createClass({
    mixins: [Router.Navigation],

    handleSelectTag: function() {
        this.transitionTo("search", null, {q: this.props.tag.name});
    },
    
    render: function() {
        var tag = this.props.tag;
        return (
        <div className="col-xs-6 col-md-3" onClick={this.handleSelectTag}>
            <div className="thumbnail">
                <img alt={tag.name} className="img-responsive" src={"uploads/thumbnails/" + tag.photo}  />
                <div className="caption">
                    <h3>#{tag.name}</h3>
                </div>
            </div>
        </div>
        );
    }
});

var Tags = React.createClass({
    getInitialState: function() {
        return {
            tags: []
        }
    },

    componentWillMount: function() {
        console.log(TagStore);
        TagStore.addChangeListener(this._onChange);
    },

    componentWillUnmount: function() {
        TagStore.removeChangeListener(this._onChange);
    },

    componentDidMount: function() {
        Actions.getTags();
    },

    onChangeFilter: function(event) {
        event.preventDefault();
        tagFilter = this.refs.tagFilter.getDOMNode().value;
        Actions.filterTags(tagFilter);
    },

    render: function() {

        var self = this;

        return (
                <div>
    <div className="tag-control-box">
        <form className="form-inline">
            <div className="form-group">
                <input ref="tagFilter" className="form-control" type="text" placeholder="Find a tag" onChange={self.onChangeFilter} />
                <button className="btn"><i className="fa fa-sort-numeric-desc"></i>
                </button>
                <button className="btn"><i className="fa fa-sort-alpha-asc"></i>
                </button>
            </div>
        </form>
    </div>

    <div>
        {self.state.tags.map(function(tag) {
            return <Tag tag={tag} />
        })}
    </div>
</div>
        );
    },

    _onChange: function() {
        this.setState({
            tags: TagStore.getTags()
        });
    },
});

module.exports = Tags;
