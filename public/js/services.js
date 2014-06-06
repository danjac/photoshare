'use strict';

/* Services */

angular.module('photoshare.services', [])
    .service('Authenticator', ['$resource', 'urls', function ($resource, urls) {

        return {
            loggedIn: false,
            currentUser: null,
            resource: $resource(urls.auth)
        };

    }])
    .service('Photo', ['$resource', 'urls', function ($resource, urls) {
        return $resource(urls.photos, {id: '@id'});
    }]);
    
