'use strict';

/* Controllers */

var user = null;

angular.module('photoshare.controllers', ['photoshare.services'])
    .controller('AppCtrl', ['$scope', '$location', 'Authenticator', function ($scope, $location, Authenticator) {

        $scope.auth = Authenticator;

        Authenticator.resource.get({}, function (user) {
            $scope.auth.loggedIn = true;
            $scope.auth.currentUser = user;
        });

        $scope.logout = function () {
            $scope.auth.currentUser.$delete();
            $scope.auth.loggedIn = false;
            $scope.auth.currentUser = null;
            $location.path("#/list");
        };

        /*
        $scope.$on("login", function (event, newUser) {
            $scope.auth.loggedIn = true;
            $scope.auth.currentUser = newUser;
        });
        */
    }])
    .controller('ListCtrl', ['$scope', 'Photo', 'pageSize', function ($scope, Photo, pageSize) {
        var page = 1, stopScrolling = false;
        $scope.photos = [];
        $scope.nextPage = function () {
            if (!stopScrolling) {
                Photo.query({page: page}).$promise.then(function (photos) {
                    $scope.photos = $scope.photos.concat(photos);
                    if (photos.length < pageSize) {
                        stopScrolling = true;
                    } else {
                        page += 1;
                    }
                });
            }
        };
    }])
    .controller('DetailCtrl', ['$scope', '$routeParams', 'Photo', function ($scope, $routeParams, Photo) {
        $scope.photo = Photo.get({id: $routeParams.id});
    }])
    .controller('UploadCtrl', ['$scope', '$location', '$window', 'Authenticator', 'Photo', function ($scope, $location, $window, Authenticator, Photo) {
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
        $scope.loginCreds = new Authenticator.resource();
        $scope.login = function () {
            $scope.loginCreds.$save(function () {
                Authenticator.currentUser = $scope.loginCreds;
                Authenticator.loggedIn = Authenticator.currentUser !== null;
                //$scope.$emit("login", $scope.loginCreds);
                $scope.loginCreds = new Authenticator.resource();
                if (Authenticator.loggedIn) {
                    $location.path("#/list");
                }
            });
        };
    }]);
