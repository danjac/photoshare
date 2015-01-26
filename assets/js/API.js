var request = require('superagent');
var Utils = require('./Utils');

var X_AUTH_HEADER = "X-Auth-Token";

var API = {

    uploadPhoto: function(title, tags, photo, callback) {
        request
            .post("/api/photos/")
            .set(X_AUTH_HEADER, Utils.getAuthToken())
            .field("title", title)
            .field("taglist", tags)
            .attach("photo", photo)
            .end(function(res){
                callback(res.body);
            });
    },

    deletePhoto: function(photoId, callback) {
        request
            .del("/api/photos/" + photoId)
            .set(X_AUTH_HEADER, Utils.getAuthToken())
            .end(function(res){
                callback();
            });
    },

    editPhotoTitle: function(photoId, title, callback) {

        request
            .patch("/api/photos/" + photoId + "/title")
            .set(X_AUTH_HEADER, Utils.getAuthToken())
            .send({
                title: title
            })
            .end(function(res){
                callback();
            });

    },

    getPhoto: function(photoId, callback) {
        request
            .get("/api/photos/" + photoId)
            .set(X_AUTH_HEADER, Utils.getAuthToken())
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

    getUser: function(callback) {
        var token = Utils.getAuthToken();
        if (!token) {
            return;
        }

        request
            .get("/api/auth/")
            .set(X_AUTH_HEADER, token)
            .end(function(res){
                if (res.body.loggedIn) {
                    callback(res.body);
                }
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
                    Utils.setAuthToken(res.headers["x-auth-token"]);
                    onSuccess(res.body); ;
                }
            });
    },

    logout: function(callback) {
        request
            .del('/api/auth/')
            .set(X_AUTH_HEADER, Utils.getAuthToken())
            .end(function(){
                Utils.delAuthToken();
                if (callback) {
                    callback();
                }
            });
    }
};

module.exports = API;
