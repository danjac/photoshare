var AppDispatcher = require('./AppDispatcher');
var Constants = require('./Constants');
var API = require('./API')

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

    login: function(identifier, password) {
        var self = this;
        API.login(identifier, password, function(data) {
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
