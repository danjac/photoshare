var request = require('superagent');

var API = {

    getPhoto: function(photoId, callback) {
        request
            .get("/api/photos/" + photoId)
            .end(function(res){
                callback(res.body);
            });
    },

    getPhotos: function(orderBy, callback) {
        request
            .get("/api/photos/")
            .query({
                orderBy: orderBy || ''
            })
            .end(function(res) {
                callback(res.body);
            });
    },

    login: function(identifier, password, onSuccess, onError) {
        request
            .post("/api/auth/")
            .send({
                identifier: identifier,
                password: password
            })
            .end(function(res){
                if (res.badRequest) {
                    onError(res.text);
                } else {
                    onSuccess(res.body);
                }
            });
    },

    logout: function() {
        request.del('/api/auth/');
    }
};

module.exports = API;
