'use strict';


// Declare app level module which depends on filters, and services
angular.module('photoshare', [
    'ngRoute',
    'ngResource',
    'ngAnimate',
    'infinite-scroll',
    'photoshare.filters',
    'photoshare.services',
    'photoshare.directives',
    'photoshare.controllers'
]).
    constant('urls', {
        auth: '/api/auth/',
        photos: '/api/photos/:id',
        users: '/api/user/',
        tags: '/api/tags/'
    }).
    constant('pageSize', 32).
    constant('authToken', 'X-Auth-Token').
    config(['$routeProvider',
            '$locationProvider',
            '$httpProvider',
            '$resourceProvider', function (
        $routeProvider,
        $locationProvider,
        $httpProvider,
        $resourceProvider
    ) {
        $routeProvider.
        
            when('/list', {templateUrl: 'partials/list.html', controller: 'ListCtrl'}).

            when('/tags', {templateUrl: 'partials/tags.html', controller: 'TagsCtrl'}).

            when('/search/:q', {templateUrl: 'partials/list.html', controller: 'ListCtrl'}).

            when('/owner/:ownerID/:ownerName', {templateUrl: 'partials/list.html', controller: 'ListCtrl'}).

            when('/detail/:id', {templateUrl: 'partials/detail.html', controller: 'DetailCtrl'}).

            when('/upload', {templateUrl: 'partials/upload.html', controller: 'UploadCtrl'}).

            when('/login', {templateUrl: 'partials/login.html', controller: 'LoginCtrl'}).

            when('/signup', {templateUrl: 'partials/signup.html', controller: 'SignupCtrl'}).

            otherwise({redirectTo: '/list'});
        //$locationProvider.html5Mode(true);
        //
        $resourceProvider.defaults.stripTrailingSlashes = false;

        //$httpProvider.defaults.xsrfCookieName = "csrf_token";
        //$httpProvider.defaults.xsrfHeaderName = "X-CSRF-Token";

        // handle file uploads

        $httpProvider.defaults.transformRequest = function (data, headersGetter) {

            if (data === undefined) {
                return data;
            }

            var fd = new FormData(),
                isFileUpload = false,
                headers = headersGetter();

            angular.forEach(data, function (value, key) {
                if (value instanceof FileList) {
                    isFileUpload = true;
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
            if (isFileUpload) {
                headers["Content-Type"] = undefined;
                return fd;
            }

            return JSON.stringify(data);
        };

        var interceptors = ['AuthInterceptor', 'ErrorInterceptor'];

        angular.forEach(interceptors, function (interceptor) {
            $httpProvider.interceptors.push([
                '$injector', function ($injector) {
                    return $injector.get(interceptor);
                }
            ]);
        });

    }]).factory('AuthInterceptor', function ($window) {
        
        return {
            request: function (config) {
                config.headers = config.headers || {};
                if ($window.sessionStorage.token) {
                    config.headers['X-Auth-Token'] = $window.sessionStorage.token;
                }
                return config;
            }
        };

    }).factory('ErrorInterceptor', function ($q, $location, $window, Session, Alert) {
        return {

            response: function (response) {
                return response;
            },

            responseError: function (response) {
                var rejection = $q.reject(response);

                if (response.status === 401) {
                    Alert.danger(angular.fromJson(response.data));
                    Session.clear();
                    Session.setLastLoginUrl();
                    $location.path("/login");
                }
                if (response.status === 400) {
                    if (response.data.errors) {
                        // TBD: render the specific form errors  
                        Alert.danger("Sorry, your form contains errors, please try again");
                    } else {
                        Alert.danger(angular.fromJson(response.data));
                    }
                }
                if (response.status === 500) {
                    Alert.danger("Sorry, an error has occurred");
                }
                return rejection;
            }
        };
    });
