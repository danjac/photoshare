'use strict';

/* Controllers */

angular.module('photoshare.controllers', ['photoshare.services'])
    .controller('AppCtrl', ['$scope',
                            '$location',
                            '$timeout',
                            'Session',
                            'Authenticator',
                            'Alert',
                            function ($scope,
                                      $location,
                                      $timeout,
                                      Session,
                                      Authenticator,
                                      Alert) {

            $scope.session = Session;
            $scope.alert = Alert;
            $scope.searchQuery = "";

            Authenticator.init();

            $scope.$watch('alert.message', function (newValue, oldValue) {
                if (newValue) {
                    $timeout(function () { Alert.dismiss(); }, 3000);
                }
            });

            $scope.logout = function () {
                Authenticator.logout().then(function () {
                    $location.path("/popular");
                });
            };

            $scope.login = function () {
                Session.setLastLoginUrl();
                $location.path("/login");
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
                             function ($scope,
                                       $location,
                                       $routeParams,
                                       Photo) {
            var pageLoaded = false,
                apiCall = null,
                q = $routeParams.q || "",
                ownerID = $routeParams.ownerID || "",
                ownerName = $routeParams.ownerName || "",
                orderBy = $location.path() == "/popular" ? "votes" : "";

            $scope.photos = [];
            $scope.searchQuery = q;
            $scope.ownerName = ownerName;
            $scope.searchComplete = false;
            $scope.total = 0;
            $scope.currentPage = 0;
            $scope.numPages = 0;

            if (q) {
                apiCall = function (page) { return Photo.search({ q: q, page: page })};
            } else if (ownerID) {
                apiCall = function (page) { return Photo.byOwner({ ownerID: ownerID, page: page })};
            } else {
                apiCall = function (page) { return Photo.query({ orderBy: orderBy, page: page })};
            }

            $scope.nextPage = function (page) {
                pageLoaded = false;
                apiCall(page).$promise.then(function (result) {
                    $scope.pageLoaded = true;
                    $scope.searchComplete = true;
                    $scope.photos = result.photos;
                    $scope.total = result.total;
                    $scope.numPages = result.numPages;
                    $scope.currentPage = page;

                });
            };
            $scope.nextPage(1);

            $scope.getDetail = function (photo) {
                $location.path("/detail/" + photo.id);
            };

        }])

    .controller('DetailCtrl', ['$scope',
                               '$routeParams',
                               '$location',
                               '$window',
                               'Photo',
                               'Tag',
                               'Session',
                               'Alert',
                               function ($scope,
                                         $routeParams,
                                         $location,
                                         $window,
                                         Photo,
                                         Tag,
                                         Session,
                                         Alert) {

            $scope.photo = null;
            $scope.editTitle = false;
            $scope.editTags = false;
            $scope.pageLoaded = false;

            Photo.get({id: $routeParams.id}).$promise.then(function (photo) {
                $scope.photo = photo;
                $scope.photo.taglist = $scope.photo.tags ? $scope.photo.tags.join(" ") : "";
                $scope.pageLoaded = true;
            });

            $scope.voteUp = function () {
                if (!$scope.photo.perms.vote) {
                    return;
                }
                $scope.photo.perms.vote = false;
                $scope.photo.upVotes += 1
                Photo.upvote({id: $scope.photo.id});
            }

            $scope.voteDown = function () {
                if (!$scope.photo.perms.vote) {
                    return;
                }
                $scope.photo.perms.vote = false;
                $scope.photo.downVotes += 1
                Photo.downvote({id: $scope.photo.id});
            }

            $scope.deletePhoto = function () {
                if (!$scope.photo.perms.delete || !$window.confirm('You sure you want to delete this?')) {
                    return;
                }
                $scope.photo.$delete(function () {
                    Alert.warning('Your photo has been deleted');
                    $location.path("/");
                });
            };
            $scope.showEditForm = function () {
                if ($scope.photo.perms.edit) {
                    $scope.editTitle = true;
                }
            };
            $scope.hideEditForm = function () {
                $scope.editTitle = false;
            };
            $scope.showEditTagsForm = function () {
                if ($scope.photo.perms.edit) {
                    $scope.editTags = true;
                }
            };
            $scope.hideEditTagsForm = function () {
                $scope.editTags = false;
            };
            $scope.updateTitle = function () {
                Photo.updateTitle(
                    {id: $scope.photo.id, title: $scope.photo.title }
                );
                $scope.editTitle = false;
            };
            $scope.updateTags = function () {
                var taglist = $scope.photo.taglist || "";
                if (!taglist) {
                    $scope.photo.tags = [];
                } else {
                    $scope.photo.tags = taglist.trim().split(" ");
                }
                Photo.updateTags({id: $scope.photo.id, tags: $scope.photo.tags });
                $scope.editTags = false;
            };

        }])

    .controller('TagsCtrl', ['$scope',
                             '$location',
                             'Tag', function ($scope, $location, Tag) {
        $scope.tags = [];
        $scope.orderField = '-numPhotos';
        $scope.pageLoaded = false;

        Tag.query().$promise.then(function (tags) {
            $scope.tags = tags;
            $scope.pageLoaded = true;
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
                               'Session',
                               'Alert',
                               'Photo', function ($scope,
                                                  $location,
                                                  Session,
                                                  Alert,
                                                  Photo) {
        Session.check();
        $scope.newPhoto = new Photo();
        $scope.upload = null;
        $scope.formDisabled = false;
        $scope.formErrors = {};
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
                    $location.path("/latest");
                },
                function (result) {
                    $scope.formErrors = result.data.errors;
                    $scope.formDisabled = false;
                }
            );
        };

    }])

    .controller('LoginCtrl', ['$scope',
                              '$location',
                              'Session',
                              'Authenticator',
                              'Alert',
                              'authToken', function ($scope,
                                                     $location,
                                                     Session,
                                                     Authenticator,
                                                     Alert,
                                                     authToken) {

        $scope.loginCreds = new Authenticator.resource();
        $scope.login = function () {
            $scope.loginCreds.$save(function (result, headers) {
                $scope.loginCreds = new Authenticator.resource();
                if (result.loggedIn) {
                    Authenticator.login(result, headers(authToken));
                    Alert.success("Welcome back, " + result.name);
                    var path = Session.getLastLoginUrl() || "/popular";
                    if (path == $location.path()) {
                        path = "/popular";
                    }
                    $location.path(path);
                }
            });
        };
    }])

    .controller('SignupCtrl', ['$scope',
                               '$location',
                               'User',
                               'Authenticator',
                               'Alert',
                               'authToken', function ($scope,
                                                      $location,
                                                      User,
                                                      Authenticator,
                                                      Alert,
                                                      authToken) {

        $scope.newUser = new User();
        $scope.formErrors = {};
        $scope.signup = function () {
            $scope.newUser.$save(function (result, headers) {
                Authenticator.login(result, headers(authToken));
                $scope.newUser = new User();
                Alert.success("Welcome, " + result.name);
                $location.path("/popular");
            }, function (result) {
                $scope.formErrors = result.data.errors;
            });
        };
    }]);
