var AppDispatcher = require('./AppDispatcher');
var Constants = require('./Constants');
var API = require('./API')

var Actions = {

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
    }
};

module.exports = Actions;
