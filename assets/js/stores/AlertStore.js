var AppDispatcher = require('../AppDispatcher');
var Constants = require('../Constants');
var assign = require('object-assign');

var BaseStore = require('./BaseStore');

var _messages = [];

var AlertStore = assign({}, BaseStore, {

    getMessages: function (){
        return _messages;
    }

});

AlertStore.dispatchToken = AppDispatcher.register(function(action){

    switch(action.actionType){
        case Constants.NEW_ALERT_MESSAGE:
            if (action.message) {
                _messages.push(action.message);
                AlertStore.emitChange();
            }
            break;
        default:
    }
});

module.exports = AlertStore;
