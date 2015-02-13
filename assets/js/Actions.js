var AppDispatcher = require('./AppDispatcher');
var Constants = require('./Constants');
var API = require('./API');

var Actions = {

    alertMessage: function(msg, msgType) {
        AppDispatcher.handleViewAction({
            actionType: Constants.Actions.NEW_ALERT_MESSAGE,
            message: {
                message: msg,
                type: msgType
            }
        });
    },

    filterTags: function(tagFilter) {
        AppDispatcher.handleViewAction({
            actionType: Constants.Actions.FILTER_TAGS,
            tagFilter: tagFilter
        });
    },

    getTags: function() {
        API.getTags(function(data) {
            AppDispatcher.handleServerAction({
                actionType: Constants.Actions.GET_TAGS,
                tags: data
            });
        });
    },

    searchPhotos: function(search, page) {
        if (!search) {
            AppDispatcher.handleViewAction({
                actionType: Constants.Actions.GET_PHOTOS,
                photos: []
            });
        }
        page = page || 1;
        API.searchPhotos(search, page, function(data) {
            AppDispatcher.handleServerAction({
                actionType: Constants.Actions.GET_PHOTOS,
                photos: data
            });
        });
    },

    getPhotosForUser: function(userId, page) {
        API.getPhotosForUser(userId, page, function(data) {
            AppDispatcher.handleServerAction({
                actionType: Constants.Actions.GET_PHOTOS,
                photos: data
            });

        });
    },

    photoEditMode: function() {
        AppDispatcher.handleViewAction({
            actionType: Constants.Actions.PHOTO_EDIT_MODE
        });
    },

    photoEditDone: function(photoId, newTitle) {

        AppDispatcher.handleViewAction({
            actionType: Constants.Actions.PHOTO_EDIT_DONE,
            title: newTitle
        });
        API.editPhotoTitle(photoId, newTitle);
    },

    previewPhoto: function(photo) {

        if (window.FileReader === null) {
            return;
        }

        var reader = new window.FileReader();

        reader.onload = function(){
            AppDispatcher.handleViewAction({
                actionType: Constants.Actions.NEW_PHOTO_PREVIEW,
                url: reader.result
            });
        };
        reader.readAsDataURL(photo);
    },

    uploadPhoto: function(title, tags, photo){
        AppDispatcher.handleViewAction({
            actionType: Constants.Actions.UPLOAD_STARTED
        });
        API.uploadPhoto(title, tags, photo, function(data) {
            AppDispatcher.handleServerAction({
                actionType: Constants.Actions.NEW_PHOTO,
                photo: data
            });
            this.alertMessage('Your photo has been uploaded', Constants.Alerts.SUCCESS);
        }).bind(this);
    },

    getPhotos: function(orderBy, page) {
        page = page || 1;
        API.getPhotos(orderBy, page, function(data){
            AppDispatcher.handleServerAction({
                actionType: Constants.Actions.GET_PHOTOS,
                photos: data
            });
        });
    },

    getPhotoDetail: function(photoId) {
        API.getPhoto(photoId, function(data){
            AppDispatcher.handleServerAction({
                actionType: Constants.Actions.GET_PHOTO_DETAIL,
                photo: data
            });
        });
    },

    deletePhoto: function(photoId) {
        API.deletePhoto(photoId, function() {
            this.alertMessage("Your photo has been deleted", Constants.Alerts.ALERT_SUCCESS);
            AppDispatcher.handleServerAction({
                actionType: Constants.Actions.PHOTO_DELETED
            });
        }).bind(this);
    },

    getUser: function() {
        API.getUser(function(data) {
            AppDispatcher.handleServerAction({
                actionType: Constants.Actions.LOGIN_SUCCESSFUL,
                user: data
            });
        });
    },

    login: function(identifier, password) {
        API.login(identifier, password, function(data, authToken) {
            this.alertMessage("Welcome back!", Constants.Alerts.ALERT_SUCCESS);
            AppDispatcher.handleServerAction({
                actionType: Constants.Actions.LOGIN_SUCCESSFUL,
                user: data
            });
        }, function(err) {
            this.alertMessage(err, Constants.Alerts.ALERT_DANGER);
        }).bind(this);
    },

    logout: function() {
        AppDispatcher.handleServerAction({
            actionType: Constants.Actions.LOGOUT
        });
        this.alertMessage("Bye for now!", Constants.Alerts.ALERT_SUCCESS);
        API.logout();
    }
};

module.exports = Actions;
