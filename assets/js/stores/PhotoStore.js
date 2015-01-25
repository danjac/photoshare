var AppDispatcher = require('../AppDispatcher');
var Constants = require('../Constants');
var EventEmitter = require('events').EventEmitter;
var assign = require('object-assign');

var _photos = {
    photos: []
};

var _photoDetail = {};

var PhotoStore = assign({}, EventEmitter.prototype, {

    getPhotos: function (){
        return _photos;
    },

    getPhotoDetail: function(){
        return _photoDetail;
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
        case Constants.GET_PHOTO_DETAIL:
            _photoDetail = action.photo;
            PhotoStore.emitChange();
            break;
        default:
    }
});

module.exports = PhotoStore;
