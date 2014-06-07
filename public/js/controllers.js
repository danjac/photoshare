'use strict';

/* Controllers */

var user = null;

angular.module('photoshare.controllers', ['photoshare.services'])
    .controller('AppCtrl', ['$scope',
                            '$location',
                            '$timeout',
                            'Authenticator',
                            'Alert',
                            function ($scope, $location, $timeout, Authenticator, Alert) {

            $scope.auth = Authenticator;
            $scope.alert = Alert;

            $scope.$watch('alert.message', function (newValue, oldValue) {
                if (newValue) {
                    $timeout(function () { Alert.dismiss(); }, 3000);
                }
            });

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
                    }
                });
            }
            page += 1;
        };
    }])

    .controller('DetailCtrl', ['$scope',
                               '$routeParams',
                               '$location',
                               'Photo',
                               'Authenticator',
                               'Alert',
                               function ($scope, $routeParams, $location, Photo, Authenticator, Alert) {
            $scope.photo = null;
            $scope.editTitle = false;
            Photo.get({id: $routeParams.id}).$promise.then(function (photo) {
                $scope.photo = photo;
                $scope.canDelete = (
                    Authenticator.loggedIn &&
                    ($scope.photo.ownerId === Authenticator.currentUser.id ||
                            Authenticator.currentUser.isAdmin)
                );
                $scope.canEdit = (
                    Authenticator.loggedIn && $scope.photo.ownerId ===
                    Authenticator.currentUser.id
                );
            });
            $scope.deletePhoto = function () {
                $scope.photo.$delete();
                Alert.warning('Your photo has been deleted');
                $location.path("/");
            };
            $scope.showEditForm = function () {
                if ($scope.canEdit) {
                    $scope.editTitle = true;
                }
            };
            $scope.hideEditForm = function () {
                $scope.editTitle = false;
            };
            $scope.updateTitle = function () {
                $scope.photo.$save(function () {
                    $scope.editTitle = false;
                });
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
                Alert.success('Your photo has been uploaded');
                $location.path("#/list");
            });
        };

    }])

    .controller('LoginCtrl', ['$scope',
                              '$location',
                              'Authenticator',
                              'Alert', function ($scope, $location, Authenticator, Alert) {
        $scope.loginCreds = new Authenticator.resource();
        $scope.login = function () {
            $scope.loginCreds.$save(function () {
                Authenticator.currentUser = $scope.loginCreds;
                Authenticator.loggedIn = Authenticator.currentUser !== null;
                //$scope.$emit("login", $scope.loginCreds);
                $scope.loginCreds = new Authenticator.resource();
                if (Authenticator.loggedIn) {
                    Alert.success("Welcome back, " + Authenticator.currentUser.name);
                    $location.path("#/list");
                }
            });
        };
    }])

    .controller('SignupCtrl', ['$scope',
                               '$location',
                               'User',
                               'Authenticator',
                               'Alert', function ($scope, $location, User, Authenticator, Alert) {

        $scope.newUser = new User();
        $scope.signup = function () {
            $scope.newUser.$save(function () {
                Authenticator.currentUser = $scope.newUser;
                Authenticator.loggedIn = true;
                Alert.success("Welcome, " + $scope.newUser.name);
                $location.path("#/list");
            });
        };
    }]);
