var AppDispatcher = require('../AppDispatcher');
var Constants = require('../Constants').Actions;
var assign = require('object-assign');

var BaseStore = require('./BaseStore');

var _tags = [];
var _filteredTags = [];

function _filterTags(tagFilter){
    if (!tagFilter) {
        return _tags;
    }

    var rv = [];

    _tags.forEach(function(tag) {
        if (tag.name.indexOf(tagFilter) !== -1) {
            rv.push(tag);
        }
    });
    return rv;

}

var TagStore = assign({}, BaseStore, {

    getTags: function() {
        return _filteredTags;
    }

});


TagStore.dispatchToken = AppDispatcher.register(function(action) {
    switch(action.actionType) {
        case Constants.GET_TAGS:
            _tags = _filteredTags = action.tags;
            TagStore.emitChange();
            break;
        case Constants.FILTER_TAGS:
            _filteredTags = _filterTags(action.tagFilter);
            TagStore.emitChange();
            break;
    }

});

module.exports = TagStore;
