'use strict';


// Declare app level module which depends on filters, and services
angular.module('photoshare', [
    'ngRoute',
    'ngResource',
    'photoshare.filters',
    'photoshare.services',
    'photoshare.directives',
    'photoshare.controllers'
]).
    config(['$routeProvider', '$locationProvider', '$httpProvider', '$resourceProvider', function ($routeProvider, $locationProvider, $httpProvider, $resourceProvider) {
        $routeProvider.when('/list', {templateUrl: 'partials/list.html', controller: 'ListCtrl'});
        $routeProvider.when('/upload', {templateUrl: 'partials/upload.html', controller: 'UploadCtrl'});
        $routeProvider.when('/login', {templateUrl: 'partials/login.html', controller: 'LoginCtrl'});
        $routeProvider.otherwise({redirectTo: '/list'});
        //$locationProvider.html5Mode(true);
        //
        $resourceProvider.defaults.stripTrailingSlashes = false;

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
        $httpProvider.interceptors.push([
            '$injector', function ($injector) {
                return $injector.get('AuthInterceptor');
            }
        ]);
        $httpProvider.defaults.headers.post['Content-Type'] = undefined;
    }]).factory('AuthInterceptor', function ($q, $location) {
        return {
            response: function (response) {
                return response;
            },

            responseError: function (response) {
                var rejection = $q.reject(response);

                if (response.status === 401) {
                    $location.path("/login");
                    return rejection;
                }
                return rejection;
            }
        };
    });
