'use strict';

/* Controllers */

var user = null;

angular.module('photoshare.controllers', ['photoshare.services'])
    .controller('AppCtrl', ['$scope',
                            '$location',
                            'Authenticator',
                            'Alert',
                            function ($scope, $location, Authenticator, Alert) {

            $scope.auth = Authenticator;
            $scope.alert = Alert;

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

    .controller('DetailCtrl', ['$scope',
                               '$routeParams',
                               '$location',
                               'Photo',
                               'Authenticator',
                               function ($scope, $routeParams, $location, Photo, Authenticator) {
            $scope.photo = null;
            Photo.get({id: $routeParams.id}).$promise.then(function (photo) {
                $scope.photo = photo;
                $scope.canDelete = (
                    Authenticator.loggedIn &&
                    ($scope.photo.ownerId === Authenticator.currentUser.id ||
                            Authenticator.currentUser.isAdmin)
                );
            });
            $scope.deletePhoto = function () {
                $scope.photo.$delete();
                $location.path("/");
            };

        }])

    .controller('UploadCtrl', ['$scope',
                               '$location',
                               '$window',
                               'Authenticator',
                               'Alert',
                               'Photo', function ($scope, $location, $window, Authenticator, Alert, Photo) {
        if (!Authenticator.currentUser) {
            $location.path("#/list");
            return;
        }
        $scope.newPhoto = new Photo();
        $scope.upload = null;
        $scope.uploadPhoto = function () {
            $scope.newPhoto.$save(function () {
                $scope.newPhoto = new Photo();
                Alert.addMessage('Your photo has been uploaded', 'success');
                $location.path("#/list");
            });
        };

    }])

    .controller('LoginCtrl', ['$scope', '$location', 'Authenticator', 'Alert', function ($scope, $location, Authenticator, Alert) {
        $scope.loginCreds = new Authenticator.resource();
        $scope.login = function () {
            $scope.loginCreds.$save(function () {
                Authenticator.currentUser = $scope.loginCreds;
                Authenticator.loggedIn = Authenticator.currentUser !== null;
                //$scope.$emit("login", $scope.loginCreds);
                $scope.loginCreds = new Authenticator.resource();
                if (Authenticator.loggedIn) {
                    Alert.addMessage("Welcome back, " + Authenticator.currentUser.name, "success");
                    $location.path("#/list");
                }
            });
        };
    }]);
