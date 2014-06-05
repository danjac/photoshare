'use strict';

/* Services */

angular.module('photoshare.services', [])
    .service('Authenticator', ['$resource', '$q', function ($resource, $q) {

        return {
            loggedIn: false,
            currentUser: null,
            resource: $resource("/auth/")
        };

    }])
    .service('Photo', ['$resource', function ($resource) {
        return $resource("/photos/");
    }]);
    