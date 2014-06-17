'use strict';

/* Controllers */

angular.module('photoshare.controllers', ['photoshare.services'])
    .controller('AppCtrl', ['$scope',
                            '$location',
                            '$timeout',
                            'Authenticator',
                            'Alert',
                            function ($scope,
                                      $location,
                                      $timeout,
                                      Authenticator,
                                      Alert) {

            $scope.auth = Authenticator;
            $scope.alert = Alert;
            $scope.searchQuery = "";

            $scope.auth.init();

            $scope.$watch('alert.message', function (newValue, oldValue) {
                if (newValue) {
                    $timeout(function () { Alert.dismiss(); }, 3000);
                }
            });

            $scope.logout = function () {
                $scope.auth.logout().then(function () {
                    $location.path("/list");
                });
            };

            $scope.doSearch = function () {
                $location.path("/search/" + $scope.searchQuery);
                $scope.searchQuery = "";
            };
        }])

    .controller('ListCtrl', ['$scope',
                             '$location',
                             '$routeParams',
                             'Photo',
                             'pageSize',
                             function ($scope, $location, $routeParams, Photo, pageSize) {
            var page = 1,
                stopScrolling = false,
                q = $routeParams.q || "",
                ownerID = $routeParams.ownerID || "",
                ownerName = $routeParams.ownerName || "";

            $scope.photos = [];
            $scope.searchQuery = q;
            $scope.ownerName = ownerName;
            $scope.nextPage = function () {
                if (!stopScrolling) {
                    Photo.query({page: page, q: q, ownerID: ownerID}).$promise.then(function (photos) {
                        $scope.photos = $scope.photos.concat(photos);
                        if (photos.length < pageSize) {
                            stopScrolling = true;
                        }
                    });
                }
                page += 1;
            };

            $scope.getDetail = function (photo) {
                $location.path("/detail/" + photo.id);
            };

        }])

    .controller('DetailCtrl', ['$scope',
                               '$routeParams',
                               '$location',
                               'Photo',
                               'Tag',
                               'Authenticator',
                               'Alert',
                               function ($scope, $routeParams, $location, Photo, Tag, Authenticator, Alert) {

            function doUpdate(onSuccess) {
                var taglist = $scope.photo.taglist || "";
                if (!taglist) {
                    $scope.photo.tags = [];
                } else {
                    $scope.photo.tags = taglist.trim().split(" ");
                }
                Photo.update({id: $scope.photo.id,
                              title: $scope.photo.title,
                              tags: $scope.photo.tags}, function () {
                    onSuccess();
                });
            }

            $scope.photo = null;
            $scope.editTitle = false;
            $scope.editTags = false;

            Photo.get({id: $routeParams.id}).$promise.then(function (photo) {
                $scope.photo = photo;
                $scope.canDelete = Authenticator.canDelete($scope.photo);
                $scope.canEdit = Authenticator.canEdit($scope.photo);
                $scope.photo.taglist = $scope.photo.tags ? $scope.photo.tags.join(" ") : "";
            });
            $scope.deletePhoto = function () {
                $scope.photo.$delete(function () {
                    Alert.warning('Your photo has been deleted');
                    $location.path("/");
                });
            };
            $scope.showEditForm = function () {
                if ($scope.canEdit) {
                    $scope.editTitle = true;
                }
            };
            $scope.hideEditForm = function () {
                $scope.editTitle = false;
            };
            $scope.showEditTagsForm = function () {
                if ($scope.canEdit) {
                    $scope.editTags = true;
                }
            };
            $scope.hideEditTagsForm = function () {
                $scope.editTags = false;
            };
            $scope.updateTitle = function () {
                doUpdate(function () { $scope.editTitle = false; });
            };
            $scope.updateTags = function () {
                doUpdate(function () {
                    $scope.editTags = false;
                });
            };

        }])

    .controller('TagsCtrl', ['$scope',
                             '$location',
                             'Tag', function ($scope, $location, Tag) {
        $scope.tags = [];
        $scope.orderField = 'name';

        Tag.query().$promise.then(function (tags) {
            $scope.tags = tags;
        });

        $scope.doSearch = function (tag) {
            $location.path("/search/" + tag);
        };

        $scope.orderTags = function (field) {
            $scope.orderField = field;
        };

    }])

    .controller('UploadCtrl', ['$scope',
                               '$location',
                               '$window',
                               'Authenticator',
                               'Alert',
                               'Photo', function ($scope, $location, $window, Authenticator, Alert, Photo) {
        if (!Authenticator.session.loggedIn) {
            Alert.danger("You must be logged in");
            $location.path("/login");
            return;
        }
        $scope.newPhoto = new Photo();
        $scope.upload = null;
        $scope.formDisabled = false;
        $scope.uploadPhoto = function () {
            $scope.formDisabled = true;
            var taglist = $scope.newPhoto.taglist || "";
            if (!taglist) {
                $scope.newPhoto.tags = [];
            } else {
                $scope.newPhoto.tags = taglist.trim().split(" ");
            }
            $scope.newPhoto.$save(
                function () {
                    $scope.newPhoto = new Photo();
                    Alert.success('Your photo has been uploaded');
                    $location.path("/list");
                },
                function () {
                    $scope.formDisabled = false;
                }
            );
        };

    }])

    .controller('LoginCtrl', ['$scope',
                              '$location',
                              '$window',
                              'Authenticator',
                              'Alert', function ($scope, $location, $window, Authenticator, Alert) {
        $scope.loginCreds = new Authenticator.resource();
        $scope.login = function () {
            $scope.loginCreds.$save(function (result) {
                $scope.loginCreds = new Authenticator.resource();
                if (result.loggedIn) {
                    Authenticator.session = result;
                    $window.sessionStorage.token = result.token;
                    Alert.success("Welcome back, " + Authenticator.session.name);
                    $location.path("/list");
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
            $scope.newUser.$save(function (result) {
                Authenticator.session = result;
                $scope.newUser = new User();
                Alert.success("Welcome, " + result.name);
                $location.path("/list");
            });
        };
    }]);
