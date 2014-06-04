'use strict';

/* Controllers */

var user = null;

angular.module('photoshare.controllers', ['photoshare.services'])
    .controller('AppCtrl', ['$scope', '$location', 'Authenticator', function ($scope, $location, Authenticator) {

        Authenticator.authenticate().then(function () {
            $scope.auth = Authenticator;
        });

        $scope.logout = function () {
            Authenticator.logout();
            $location.path("#/list");
        };
    }])
    .controller('ListCtrl', ['$scope', 'Photo', function ($scope, Photo) {
        Photo.query().$promise.then(function (photos) {
            $scope.photos = photos;
        });
    }])
    .controller('UploadCtrl', ['$scope', '$location', '$window', 'Authenticator', 'Photo', function ($scope, $location, $window, Authenticator, Photo) {
        if (!Authenticator.isLoggedIn()) {
            $location.path("#/login");
            return;
        }
        $scope.newPhoto = new Photo();
        $scope.upload = null;
        $scope.uploadPhoto = function () {
            $scope.newPhoto.$save(function () {
                $scope.newPhoto = new Photo();
                $location.path("#/list");
            });
        };

    }])
    .controller('LoginCtrl', ['$scope', '$location', 'Authenticator', function ($scope, $location, Authenticator) {
        $scope.loginCreds = {};
        $scope.login = function () {
            Authenticator.login($scope.loginCreds.email, $scope.loginCreds.password).then(function (user) {
                console.log("USER", user);
                if (user) {
                    $scope.loginCreds = {};
                    $location.path("#/list");
                }
            });
        };
    }]);
