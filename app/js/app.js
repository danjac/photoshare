'use strict';


// Declare app level module which depends on filters, and services
angular.module('photoshare', [
    'ngRoute',
    'photoshare.filters',
    'photoshare.services',
    'photoshare.directives',
    'photoshare.controllers'
]).
    config(['$routeProvider', '$locationProvider', function ($routeProvider, $locationProvider) {
        $routeProvider.when('/list', {templateUrl: 'partials/list.html', controller: 'ListCtrl'});
        $routeProvider.when('/upload', {templateUrl: 'partials/upload.html', controller: 'UploadCtrl'});
        $routeProvider.when('/login', {templateUrl: 'partials/login.html', controller: 'LoginCtrl'});
        $routeProvider.otherwise({redirectTo: '/list'});
        //$locationProvider.html5Mode(true);
    }]);
