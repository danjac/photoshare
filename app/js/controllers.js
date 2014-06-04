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
        Photo.query().then(function (photos) {
            $scope.photos = photos;
        });
    }])
    .controller('UploadCtrl', ['$scope', '$location', '$http', '$window', 'Authenticator', function ($scope, $location, $http, $window, Authenticator) {
        if (!Authenticator.isLoggedIn()) {
            $location.path("/login");
            return;
        }
        $scope.newPhoto = {};
        $scope.upload = null;
        $scope.uploadPhoto = function () {
            $http.post("/add", $scope.newPhoto, function (result) {
                console.log("RESULT", result);
            });
            $scope.newPhoto = {};
            $location.path("#/list");
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
