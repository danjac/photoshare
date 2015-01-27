var request = require('superagent');
var Constants = require('./Constants');

var X_AUTH_HEADER = "X-Auth-Token";

var _setAuthToken = function(token) {
    window.localStorage.setItem(Constants.AUTH_TOKEN_STORAGE_KEY, token);
};

var _getAuthToken = function() {
    return window.localStorage.getItem(Constants.AUTH_TOKEN_STORAGE_KEY);
};

var _delAuthToken = function() {
    window.localStorage.removeItem(Constants.AUTH_TOKEN_STORAGE_KEY);
};

var API = {

    uploadPhoto: function(title, tags, photo, callback) {
        request
            .post("/api/photos/")
            .set(X_AUTH_HEADER, _getAuthToken())
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
            .set(X_AUTH_HEADER, _getAuthToken())
            .end(function(res){
                callback();
            });
    },

    editPhotoTitle: function(photoId, title, callback) {

        request
            .patch("/api/photos/" + photoId + "/title")
            .set(X_AUTH_HEADER, _getAuthToken())
            .send({
                title: title
            })
            .end(function(res){
                if (callback) {
                    callback();
                }
            });

    },

    editPhotoTags: function(photoId, tags, callback) {

        request
            .patch("/api/photos/" + photoId + "/tags")
            .set(X_AUTH_HEADER, _getAuthToken())
            .send({
                taglist: tags
            })
            .end(function(res){
                if (callback) {
                    callback();
                }
            });

    },

    getPhoto: function(photoId, callback) {
        request
            .get("/api/photos/" + photoId)
            .set(X_AUTH_HEADER, _getAuthToken())
            .end(function(res){
                callback(res.body);
            });
    },

    getPhotos: function(orderBy, page, callback) {
        request
            .get("/api/photos/")
            .query({
                orderBy: orderBy || '',
                page: page
            })
            .end(function(res) {
                callback(res.body);
            });
    },

    searchPhotos: function(search, page, callback) {
        request
            .get("/api/photos/search")
            .query({
                q: search,
                page: page
            })
            .end(function(res){
                console.log(res)
                callback(res.body);
            });
    },

    getUser: function(callback) {
        var token = _getAuthToken();
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
                    _setAuthToken(res.headers["x-auth-token"]);
                    onSuccess(res.body); ;
                }
            });
    },

    logout: function(callback) {
        request
            .del('/api/auth/')
            .set(X_AUTH_HEADER, _getAuthToken())
            .end(function(){
                _delAuthToken();
                if (callback) {
                    callback();
                }
            });
    }
};

module.exports = API;
