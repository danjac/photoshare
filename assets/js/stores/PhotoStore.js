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
var _deleted = false;
var _editMode = false;
var _uploadStarted = false;

var PhotoStore = assign({}, EventEmitter.prototype, {

    isUploadStarted: function() {
        return _uploadStarted;
    },

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

    isDeleted: function() {
        return _deleted;
    },

    isEditMode: function(){
        return _editMode;
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

PhotoStore.dispatchToken = AppDispatcher.register(function(action){

    switch(action.actionType){
        case Constants.GET_PHOTOS:
            _photos = action.photos;
            break;
        case Constants.UPLOAD_STARED:
            _uploadStarted = true;
            break;
        case Constants.GET_PHOTO_DETAIL:
            _photoDetail = action.photo;
            _newPhoto = null;
            _deleted = false;
            _editMode = false;
            break;
        case Constants.NEW_PHOTO:
            _newPhoto = action.photo;
            _uploadStarted = false;
            break;
        case Constants.NEW_PHOTO_PREVIEW:
            _previewUrl = action.url;
            break;
        case Constants.PHOTO_DELETED:
            _deleted = true;
            break;
        case Constants.PHOTO_EDIT_MODE:
            _editMode = !(_editMode);
            break;
        case Constants.PHOTO_EDIT_DONE:
            _editMode = false;
            _photoDetail.title = action.title;
            break;
        default:
    }
    PhotoStore.emitChange();
});

module.exports = PhotoStore;
