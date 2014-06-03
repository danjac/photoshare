'use strict';


// Declare app level module which depends on filters, and services
angular.module('photoshare', [
    'ngRoute',
    'photoshare.filters',
    'photoshare.services',
    'photoshare.directives',
    'photoshare.controllers'
]).
    config(['$routeProvider', '$locationProvider', '$httpProvider', function ($routeProvider, $locationProvider, $httpProvider) {
        $routeProvider.when('/list', {templateUrl: 'partials/list.html', controller: 'ListCtrl'});
        $routeProvider.when('/upload', {templateUrl: 'partials/upload.html', controller: 'UploadCtrl'});
        $routeProvider.when('/login', {templateUrl: 'partials/login.html', controller: 'LoginCtrl'});
        $routeProvider.otherwise({redirectTo: '/list'});
        //$locationProvider.html5Mode(true);

        $httpProvider.defaults.transformRequest = function (data) {
            if (data === undefined) {
                return data;
            }
            var fd = new FormData();
            angular.forEach(data, function (value, key) {
                if (value instanceof FileList) {
                    if (value.length === 1) {
                        fd.append(key, value[0]);
                    } else {
                        angular.forEach(value, function (file, index) {
                            fd.append(key + "_" + index, file);
                        });
                    }
                } else {
                    fd.append(key, value);
                }
            });
            return fd;
        };
        $httpProvider.defaults.headers.post['Content-Type'] = undefined;
    }]);
