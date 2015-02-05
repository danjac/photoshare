var AppDispatcher = require('../AppDispatcher');
var Constants = require('../Constants');
var EventEmitter = require('events').EventEmitter;
var assign = require('object-assign');

var _tags = [];


var TagStore = assign({}, EventEmitter.prototype, {

    getTags: function() {
        return _tags;
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
            _tags = action.tags;
            TagStore.emitChange();
            break;
    }

});

module.exports = TagStore;
