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
        var self = this;
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
        var self = this;
        AppDispatcher.handleViewAction({
            actionType: Constants.Actions.UPLOAD_STARTED
        });
        API.uploadPhoto(title, tags, photo, function(data) {
            AppDispatcher.handleServerAction({
                actionType: Constants.Actions.NEW_PHOTO,
                photo: data
            });
            self.alertMessage('Your photo has been uploaded', Constants.Alerts.SUCCESS);
        });
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
        var self = this;
        API.getPhoto(photoId, function(data){
            AppDispatcher.handleServerAction({
                actionType: Constants.Actions.GET_PHOTO_DETAIL,
                photo: data
            });
        });
    },

    deletePhoto: function(photoId) {
        var self = this;
        API.deletePhoto(photoId, function() {
            self.alertMessage("Your photo has been deleted", Constants.Alerts.ALERT_SUCCESS);
            AppDispatcher.handleServerAction({
                actionType: Constants.Actions.PHOTO_DELETED
            });
        });
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
        var self = this;
        API.login(identifier, password, function(data, authToken) {
            self.alertMessage("Welcome back!", Constants.Alerts.ALERT_SUCCESS);
            AppDispatcher.handleServerAction({
                actionType: Constants.Actions.LOGIN_SUCCESSFUL,
                user: data
            });
        }, function(err) {
            self.alertMessage(err, Constants.Alerts.ALERT_DANGER);
        });
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
