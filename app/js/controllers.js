'use strict';

/* Controllers */

var user = null;

angular.module('photoshare.controllers', ['photoshare.services'])
    .controller('AppCtrl', ['$scope', '$location', 'Authenticator', function ($scope, $location, Authenticator) {
        $scope.auth = Authenticator;
        $scope.logout = function () {
            Authenticator.logout();
            $location.path("/");
        };
    }])
    .controller('ListCtrl', ['$scope', 'Photo', function ($scope, Photo) {
        $scope.photos = Photo.query();
    }])
    .controller('UploadCtrl', ['$scope', '$location', 'Authenticator', function ($scope, $location, Authenticator) {
        if (!Authenticator.isLoggedIn()) {
            $location.path("/login");
            return;
        }
        $scope.newPhoto = {};
        $scope.uploadPhoto = function () {
            $location.path("/list");
        };

    }])
    .controller('LoginCtrl', ['$scope', '$location', 'Authenticator', function ($scope, $location, Authenticator) {
        $scope.loginCreds = {};
        $scope.login = function () {
            user = Authenticator.login($scope.loginCreds.email, $scope.loginCreds.password);
            if (user) {
                $scope.loginCreds = {};
                $location.path("/");
                console.log("OK");
            }
        };
    }]);
