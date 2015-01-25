var Constants = require('./Constants');

var Utils = {

    setAuthToken: function(token) {
        window.localStorage.setItem(Constants.AUTH_TOKEN_STORAGE_KEY, token);
    },

    getAuthToken: function() {
        return window.localStorage.getItem(Constants.AUTH_TOKEN_STORAGE_KEY);
    },

    delAuthToken: function() {
        window.localStorage.removeItem(Constants.AUTH_TOKEN_STORAGE_KEY);
    }

};

module.exports = Utils;
