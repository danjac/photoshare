var request = require('superagent');

var API = {

    getPhoto: function(photoId, callback) {
        request
            .get("/api/photos/" + photoId)
            .end(function(response){
                callback(response.body);
            });
    },

    getPhotos: function(orderBy, callback) {
        request
            .get("/api/photos/")
            .query({
                orderBy: orderBy || ''
            })
            .end(function(response) {
                callback(response.body);
            });
    }
};

module.exports = API;
