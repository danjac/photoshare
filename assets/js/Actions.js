var AppDispatcher = require('./AppDispatcher');
var Constants = require('./Constants');
var API = require('./API');

var Actions = {

    alertMessage: function(msg, msgType) {
        AppDispatcher.dispatch({
            actionType: Constants.NEW_ALERT_MESSAGE,
            message: {
                message: msg,
                type: msgType
            }
        });
    },

    previewPhoto: function(photo) {

        if (window.FileReader === null) {
            return;
        }

        var reader = new window.FileReader();
        reader.onload = function(){
            AppDispatcher.dispatch({
                actionType: Constants.NEW_PHOTO_PREVIEW,
                url: reader.result
            });
        };
        reader.readAsDataURL(photo);
    },

    uploadPhoto: function(title, tags, photo){
        var self = this;
        API.uploadPhoto(title, tags, photo, function(data) {
            AppDispatcher.dispatch({
                actionType: Constants.NEW_PHOTO,
                photo: data
            });
            self.alertMessage('Your photo has been uploaded', Constants.ALERT_SUCCESS);
        });
    },

    getPhotos: function(orderBy) {
        API.getPhotos(orderBy, function(data){
            AppDispatcher.dispatch({
                actionType: Constants.GET_PHOTOS,
                photos: data
            });
        });
    },

    getPhotoDetail: function(photoId) {
        API.getPhoto(photoId, function(data){
            AppDispatcher.dispatch({
                actionType: Constants.GET_PHOTO_DETAIL,
                photo: data
            })
        });
    },

    deletePhoto: function(photoId) {
        var self = this;
        API.deletePhoto(photoId, function() {
            self.alertMessage("Your photo has been deleted", Constants.ALERT_SUCCESS);
            AppDispatcher.dispatch({
                actionType: Constants.PHOTO_DELETED
            });
        });
    },

    getUser: function() {
        API.getUser(function(data) {
            AppDispatcher.dispatch({
                actionType: Constants.LOGIN_SUCCESSFUL,
                user: data
            });
        });
    },

    login: function(identifier, password) {
        var self = this;
        API.login(identifier, password, function(data, authToken) {
            self.alertMessage("Welcome back!", Constants.ALERT_SUCCESS);
            AppDispatcher.dispatch({
                actionType: Constants.LOGIN_SUCCESSFUL,
                user: data
            });
        }, function(err) {
            self.alertMessage(err, Constants.ALERT_DANGER);
        });
    },

    logout: function() {
        AppDispatcher.dispatch({
            actionType: Constants.LOGOUT
        });
        this.alertMessage("Bye for now!", Constants.ALERT_SUCCESS);
        API.logout();
    }
};

module.exports = Actions;
