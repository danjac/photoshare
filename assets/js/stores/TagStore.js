var AppDispatcher = require('../AppDispatcher');
var Constants = require('../Constants');
var EventEmitter = require('events').EventEmitter;
var assign = require('object-assign');

var _tags = [];
var _filteredTags = [];
var _tagFilter = null;

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

var TagStore = assign({}, EventEmitter.prototype, {

    getTags: function() {
        return _filteredTags;
    },

    emitChange: function() {
        this.emit(Constants.CHANGE_EVENT);
    },

    addChangeListener: function(callback) {
        this.on(Constants.CHANGE_EVENT, callback);
    },

    removeChangeListener: function(callback) {
        this.removeListener(Constants.CHANGE_EVENT, callback);
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
