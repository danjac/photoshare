var AppDispatcher = require('../AppDispatcher');
var Constants = require('../Constants');
var EventEmitter = require('events').EventEmitter;
var assign = require('object-assign');

var _user = null;

var UserStore = assign({}, EventEmitter.prototype, {

    getUser: function() {
        return _user;
    },

    isLoggedIn: function(){
        return _user !== null;
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


UserStore.dispatchToken = AppDispatcher.register(function(action){

    switch(action.actionType){
        case Constants.LOGIN_SUCCESSFUL:
            _user = action.user;
            UserStore.emitChange();
            break;
        case Constants.LOGOUT:
            _user = null;
            UserStore.emitChange();
        default:
    }

});

module.exports = UserStore;
