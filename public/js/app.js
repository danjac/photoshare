'use strict';


// Declare app level module which depends on filters, and services
angular.module('photoshare', [
    'ngRoute',
    'ngResource',
    'infinite-scroll',
    'photoshare.filters',
    'photoshare.services',
    'photoshare.directives',
    'photoshare.controllers'
]).
    constant('urls', {
        auth: '/api/auth/',
        photos: '/api/photos/:id'
    }).
    constant('pageSize', 32).
    config(['$routeProvider',
            '$locationProvider',
            '$httpProvider',
            '$resourceProvider', function (
        $routeProvider,
        $locationProvider,
        $httpProvider,
        $resourceProvider
    ) {
        $routeProvider.when('/list', {templateUrl: 'partials/list.html', controller: 'ListCtrl'}).
            when('/detail/:id', {templateUrl: 'partials/detail.html', controller: 'DetailCtrl'}).
            when('/upload', {templateUrl: 'partials/upload.html', controller: 'UploadCtrl'}).
            when('/login', {templateUrl: 'partials/login.html', controller: 'LoginCtrl'}).
            otherwise({redirectTo: '/list'});
        //$locationProvider.html5Mode(true);
        //
        $resourceProvider.defaults.stripTrailingSlashes = false;

        $httpProvider.defaults.xsrfCookieName = "csrf_token";
        $httpProvider.defaults.xsrfHeaderName = "X-CSRF-Token";

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
    }]).factory('AuthInterceptor', function ($q, $location, Alert) {
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
                if (response.status === 400) {
                    Alert.addMessage(response, 'error');
                    return rejection;
                }
                return rejection;
            }
        };
    });
