var AppDispatcher = require('../AppDispatcher');
var Constants = require('../Constants');
var EventEmitter = require('events').EventEmitter;
var assign = require('object-assign');


var _messages = [];

var AlertStore = assign({}, EventEmitter.prototype, {

    getMessages: function (){
        return _messages;
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

AlertStore.dispatchToken = AppDispatcher.register(function(action){

    switch(action.actionType){
        case Constants.NEW_ALERT_MESSAGE:
            _messages.push(action.message);
            AlertStore.emitChange();
            break;
        default:
    }
});

module.exports = AlertStore;
