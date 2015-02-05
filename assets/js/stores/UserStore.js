var AppDispatcher = require('../AppDispatcher');
var Constants = require('../Constants');
var assign = require('object-assign');

var BaseStore = require('./BaseStore');

var _user = null;

var UserStore = assign({}, BaseStore, {

    getUser: function() {
        return _user;
    },

    isLoggedIn: function(){
        return _user !== null;
    },

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
