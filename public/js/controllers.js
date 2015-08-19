/* Controllers */

(function() {
    'use strict';
    angular.module('photoshare.controllers', ['photoshare.services'])
        .controller('AppCtrl', ['$scope',
            '$state',
            '$timeout',
            'Session',
            'Auth',
            'MessageQueue',
            'Alert',
            function($scope,
                $state,
                $timeout,
                Session,
                Auth,
                MessageQueue,
                Alert) {

                $scope.session = Session;
                $scope.alert = Alert;
                $scope.mq = MessageQueue;
                $scope.searchQuery = "";

                Session.init(Auth);

                $scope.$watchCollection('alert.messages', function(newValue, oldValue) {
                    $timeout(function() {
                        Alert.dismissLast();
                    }, 3000);
                });

                $scope.$watch('mq.newMessage', function(newValue, oldValue) {
                    if (newValue) {
                        Alert.info(newValue);
                    }
                });

                $scope.$on('$stateChangeStart', function(event, toState) {
                    if (toState.data && toState.data.loginRequired) {
                        Session.check();
                    }
                });

                $scope.logout = function() {
                    Session.logout().then(function() {
                        $state.go("popular");
                    });
                };

                $scope.login = function() {
                    Session.setLastLoginUrl();
                    $state.go("login");
                };

                $scope.doSearch = function(clear) {
                    $state.go('search', {
                        q: $scope.searchQuery
                    });
                    if (clear) {
                        $scope.searchQuery = "";
                    }
                };

            }
        ])

    .controller('FrontCtrl', ['$scope', '$state', 'Photo',
        function($scope, $state, Photo) {

            $scope.pageLoaded = false;
            $scope.photos = [];

            Photo.query({
                orderBy: "votes",
                page: 1
            }, function(response) {
                $scope.photos = response.photos;
            });

            $scope.go = function(photo) {
                $state.go('detail', {
                    id: photo.id
                });
            };

        }
    ])
        .controller('ListCtrl', ['$scope',
            '$state',
            '$stateParams',
            'Photo',
            function($scope,
                $state,
                $stateParams,
                Photo) {
                var pageLoaded = false,
                    apiCall = null,
                    q = $stateParams.q || "",
                    ownerID = $stateParams.ownerID || "",
                    ownerName = $stateParams.ownerName || "",
                    tag = $stateParams.tag || "",
                    orderBy = "";

                if ($state.current.data && $state.current.data.orderBy) {
                    orderBy = $state.current.data.orderBy;
                }
                $scope.photos = [];
                $scope.tag = tag;
                $scope.searchQuery = q;
                $scope.ownerName = ownerName;
                $scope.searchComplete = false;

                $scope.showHeader = false;

                $scope.total = 0;
                $scope.currentPage = 1;
                $scope.showPagination = false;

                if (q) {
                    apiCall = function(page) {
                        return Photo.search({
                            q: q,
                            page: page
                        });
                    };
                    $scope.showHeader = true;
                } else if (tag) {
                    apiCall = function(page) {
                        return Photo.search({
                            q: "#" + tag,
                            page: page
                        });
                    };
                    $scope.showHeader = true;
                } else if (ownerID) {
                    apiCall = function(page) {
                        return Photo.byOwner({
                            ownerID: ownerID,
                            page: page
                        });
                    };
                    $scope.showHeader = true;
                } else {
                    apiCall = function(page) {
                        return Photo.query({
                            orderBy: orderBy,
                            page: page
                        });
                    };
                }

                $scope.nextPage = function() {
                    pageLoaded = false;
                    apiCall($scope.currentPage)
                    .$promise
                    .then(function(result) {
                        if (result.total == 1) {
                            $scope.getDetail(result.photos[0]);
                        }
                        $scope.pageLoaded = true;
                        $scope.searchComplete = true;
                        $scope.photos = result.photos;
                        $scope.total = result.total;
                        $scope.showPagination = result.numPages > 1;
                    });
                };
                $scope.nextPage();

                $scope.getDetail = function(photo) {
                    $state.go('detail', {
                        id: photo.id
                    });
                };

            }
        ])

    .controller('DetailCtrl', ['$scope',
        '$stateParams',
        '$window',
        'Photo',
        'Tag',
        'Session',
        'Alert',
        function($scope,
            $stateParams,
            $window,
            Photo,
            Tag,
            Session,
            Alert) {

            $scope.photo = null;
            $scope.editTitle = false;
            $scope.editTags = false;
            $scope.pageLoaded = false;

            function calcScore() {
                if ($scope.photo) {
                    $scope.photo.score = $scope.photo.upVotes - $scope.photo.downVotes;
                }
            }

            $scope.$watch('photo.upVotes', function() {
                calcScore();
            });
            $scope.$watch('photo.downVotes', function() {
                calcScore();
            });

            Photo.get({
                id: $stateParams.id
            })
            .$promise
            .then(function(photo) {
                $scope.photo = photo;
                $scope.photo.taglist = $scope.photo.tags ? $scope.photo.tags.join(" ") : "";
                $scope.pageLoaded = true;
                calcScore();
            }, function() {
                Alert.danger("This photo no longer exists, or cannot be found");
                $window.history.back();
            });

            $scope.voteUp = function() {
                if (!$scope.photo.perms.vote) {
                    return;
                }
                $scope.photo.perms.vote = false;
                $scope.photo.upVotes += 1;
                Photo.upvote({
                    id: $scope.photo.id
                });
            };

            $scope.voteDown = function() {
                if (!$scope.photo.perms.vote) {
                    return;
                }
                $scope.photo.perms.vote = false;
                $scope.photo.downVotes += 1;
                Photo.downvote({
                    id: $scope.photo.id
                });
            };

            $scope.deletePhoto = function() {
                if (!$scope.photo.perms.delete || !$window.confirm('You sure you want to delete this?')) {
                    return;
                }
                $scope.photo.$delete(function() {
                    Alert.success('Your photo has been deleted');
                    $window.history.back();
                });
            };
            $scope.showEditForm = function() {
                if ($scope.photo.perms.edit) {
                    $scope.editTitle = true;
                }
            };
            $scope.hideEditForm = function() {
                $scope.editTitle = false;
            };
            $scope.showEditTagsForm = function() {
                if ($scope.photo.perms.edit) {
                    $scope.editTags = true;
                }
            };
            $scope.hideEditTagsForm = function() {
                $scope.editTags = false;
            };
            $scope.updateTitle = function() {
                Photo.updateTitle({
                    id: $scope.photo.id,
                    title: $scope.photo.title
                });
                $scope.editTitle = false;
            };
            $scope.updateTags = function() {
                var taglist = $scope.photo.taglist || "";
                if (!taglist) {
                    $scope.photo.tags = [];
                } else {
                    $scope.photo.tags = taglist.trim().split(" ");
                }
                Photo.updateTags({
                    id: $scope.photo.id,
                    tags: $scope.photo.tags
                });
                $scope.editTags = false;
            };

        }
    ])

    .controller('TagsCtrl', ['$scope',
        '$state',
        '$filter',
        'Tag',
        function($scope, $state, $filter, Tag) {
            $scope.tags = [];
            $scope.orderField = '-numPhotos';
            $scope.pageLoaded = false;

            Tag.query()
            .$promise
            .then(function(response) {
                $scope.tags = response;
                $scope.pageLoaded = true;
                $scope.filteredTags = $scope.tags;
            });

            $scope.filterTags = function() {
                $scope.filteredTags = $filter('filter')($scope.tags, {
                    name: $scope.tagFilter.name
                });

                if ($scope.filteredTags.length === 1) {
                    $scope.doSearch($scope.filteredTags[0].name);
                }
            };

            $scope.doSearch = function(tag) {
                $state.go('tag', {
                    tag: tag
                });
            };

            $scope.orderTags = function(field) {
                $scope.orderField = field;
            };

        }
    ])

    .controller('UploadCtrl', ['$scope',
        '$state',
        '$window',
        'Session',
        'Alert',
        'Photo',
        function($scope,
            $state,
            $window,
            Session,
            Alert,
            Photo) {
            //Session.check();
            $scope.newPhoto = new Photo();
            $scope.upload = null;
            $scope.formDisabled = false;
            $scope.formErrors = {};
            $scope.uploadPhoto = function(addAnother) {
                $scope.formDisabled = true;
                var taglist = $scope.newPhoto.taglist || "";
                if (!taglist) {
                    $scope.newPhoto.tags = [];
                } else {
                    $scope.newPhoto.tags = taglist.trim().split(" ");
                }
                $scope.newPhoto.$save(
                    function(response) {
                        $scope.newPhoto = new Photo();
                        Alert.success('Your photo has been uploaded');
                        if (addAnother) {
                            $scope.upload = null;
                            $scope.formDisabled = false;
                            $window.document.getElementById('photo_input').value = '';
                        } else {
                            $state.go('detail', {
                                id: response.id
                            });
                        }
                    },
                    function(result) {
                        if (result.data) {
                            $scope.formErrors = result.data.errors;
                        }
                        $scope.formDisabled = false;
                    }
                );
            };

        }
    ])

    .controller('LoginCtrl', ['$scope',
        '$location',
        '$window',
        '$http',
        'Session',
        'Auth',
        'Alert',
        'authTokenHeader',
        function($scope,
            $location,
            $window,
            $http,
            Session,
            Auth,
            Alert,
            authTokenHeader) {

            // tbd wrap in service
            $scope.oauth2Login = function(provider) {
                $http.get('/api/auth/oauth2/' + provider + '/url').success(function(response) {
                    $window.location.href = response;
                });
            };
            $scope.formData = new Auth();
            $scope.login = function() {
                $scope.formData.$save(function(result, headers) {
                    $scope.formData = new Auth();
                    if (result.loggedIn) {
                        Session.login(result, headers(authTokenHeader));
                        Alert.success("Welcome back, " + result.name);
                        var path = Session.getLastLoginUrl();
                        if (path) {
                            $location.path(path);
                        } else {
                            $window.history.back();
                        }
                    }
                });
            };
        }
    ])

    .controller('RecoverPassCtrl', ['$scope',
        '$window',
        'Auth',
        'Alert',
        function($scope, $window, Auth, Alert) {

            $scope.formData = {};

            $scope.recoverPassword = function() {
                Auth.recoverPassword({}, $scope.formData, function() {
                    Alert.success("Check your email for a link to change your password");
                    $window.history.back();
                }, function(result) {
                    Alert.danger(result.data);
                });
            };

        }
    ])

    .controller('ChangePassCtrl', ['$scope',
        '$location',
        '$window',
        '$state',
        'Auth',
        'Session',
        'Alert',
        function($scope, $location, $window, $state, Auth, Session, Alert) {

            var code = $location.search().code;
            $scope.formData = {};

            if (code) {
                $scope.formData.code = code;
            } else {
                Session.check();
            }

            $scope.changePassword = function() {
                Auth.changePassword({}, $scope.formData, function() {
                    Alert.success("Your password has been updated");
                    if (!Session.loggedIn) {
                        $state.go('login');
                    } else {
                        $window.history.back();
                    }
                }, function(result) {
                    Alert.danger(result.data || "Sorry, an error occurred");
                });
            };
        }
    ])

    .controller('SignupCtrl', ['$scope',
        '$state',
        'Auth',
        'Session',
        'Alert',
        'authTokenHeader',
        function($scope,
            $state,
            Auth,
            Session,
            Alert,
            authTokenHeader) {

            $scope.formData = {};
            $scope.formErrors = {};
            $scope.signup = function() {
                Auth.signup({}, $scope.formData, function(result, headers) {
                    Session.login(result, headers(authTokenHeader));
                    $scope.formData = {};
                    Alert.success("Welcome, " + result.name);
                    $state.go('upload');
                }, function(result) {
                    $scope.formErrors = result.data.errors;
                });
            };
        }
    ]);
})();
