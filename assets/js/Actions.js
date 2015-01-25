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
    }
};

module.exports = Actions;
