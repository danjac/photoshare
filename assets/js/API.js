var request = require('superagent');

var API = {

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
