var AppDispatcher = require('../AppDispatcher');
var Constants = require('../Constants');
var EventEmitter = require('events').EventEmitter;
var assign = require('object-assign');

var _photos = {
    photos: []
};

var _photoDetail = {};
var _newPhoto = null;
var _previewUrl = null;

var PhotoStore = assign({}, EventEmitter.prototype, {

    getPhotos: function (){
        return _photos;
    },

    getPhotoDetail: function(){
        return _photoDetail;
    },

    getNewPhoto: function(){
        return _newPhoto;
    },

    getPreviewUrl: function() {
        return _previewUrl;
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
            break;
        case Constants.GET_PHOTO_DETAIL:
            _photoDetail = action.photo;
            break;
        case Constants.NEW_PHOTO:
            _newPhoto = action.photo;
            break;
        case Constants.NEW_PHOTO_PREVIEW:
            _previewUrl = action.url;
            break;
        case Constants.PHOTO_DELETED:
            _photoDetail = null;
            break;
        default:
    }
    PhotoStore.emitChange();
});

module.exports = PhotoStore;
