var Constants = require('../Constants').Events;
var EventEmitter = require('events').EventEmitter;
var assign = require('object-assign');

var BaseStore = assign({}, EventEmitter.prototype, {

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


module.exports = BaseStore;
