'use strict';

/* Services */

angular.module('photoshare.services', [])
    .service('Authenticator', ['$resource', '$q', function ($resource, $q) {
        return $resource("/auth/");
    }])
    .service('Photo', ['$resource', function ($resource) {
        return $resource("/photos/");
    }]);
    
