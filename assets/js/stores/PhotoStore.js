var AppDispatcher = require('../AppDispatcher');
var Constants = require('../Constants');
var EventEmitter = require('events').EventEmitter;
var assign = require('object-assign');

var _photos = {
    photos: []
};

var PhotoStore = assign({}, EventEmitter.prototype, {

    getPhotos: function (){
        return _photos;
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

AppDispatcher.register(function(action){

    switch(action.actionType){
        case Constants.GET_PHOTOS:
            _photos = action.photos;
            PhotoStore.emitChange();
            break;
        default:
    }
});

module.exports = PhotoStore;
