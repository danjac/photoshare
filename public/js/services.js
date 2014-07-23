/* Services */

(function() {
    'use strict';
    angular.module('photoshare.services', [])
        .service('MessageQueue', ['$window',
            '$rootScope',
            'urls',
            'Alert',
            'Session',
            function($window, $rootScope, urls, Alert, Session) {

                var options = {
                    debug: true,
                    devel: true,
                    protocols_whitelist: ['websocket',
                        'xdr-streaming',
                        'xhr-streaming',
                        'iframe-eventsource',
                        'iframe-htmlfile',
                        'xdr-polling',
                        'xhr-polling',
                        'iframe-xhr-polling',
                        'jsonp-polling'
                    ]
                };

                function Mq() {
                    var $this = this,
                        newMessage = null;
                    $this.socket = new $window.SockJS(urls.messages, undefined, options);
                    $this.socket.onmessage = function(e) {
                        var msg = JSON.parse(e.data),
                            content = null,
                            photoLink = null;
                        if (msg.sender === Session.name) {
                            return;
                        }
                        if (msg.receiver && msg.receiver !== Session.name) {
                            return;
                        }
                        if (msg.photoID) {
                            photoLink = '<a href="/#/detail/' + msg.photoID + '">a photo</a>';
                        }
                        switch (msg.type) {
                            case 'login':
                                content = msg.sender + " has logged in";
                                break;
                            case 'logout':
                                content = msg.sender + " has logged out";
                                break;
                            case 'photo_deleted':
                                content = msg.sender + " has deleted a photo";
                                break;
                            case 'photo_updated':
                                content = msg.sender + " has updated " + photoLink;
                                break;
                            case 'photo_uploaded':
                                content = msg.sender + " has uploaded " + photoLink;
                                break;
                        }
                        $this.newMessage = content;
                        $rootScope.$digest();
                    };
                }

                return new Mq();
            }
        ])
        .service('Session', ['$location',
            '$window',
            '$q',
            'authTokenStorageKey',
            'Alert',
            function($location,
                $window,
                $q,
                authTokenStorageKey,
                Alert) {
                var noRedirectUrls = [
                    "/login",
                    "/changepass",
                    "/recoverpass",
                    "/signup"
                ];

                function isNoRedirectFromLogin(url) {
                    var result = false;
                    angular.forEach(noRedirectUrls, function(value) {
                        if (value == url) {
                            result = true;
                        }
                    });
                    return result;
                }

                function Session() {
                    this.clear();
                    this.lastLoginUrl = null;
                }

                Session.prototype.init = function(authResource) {
                    this.resource = authResource;
                    this.sync();
                };

                Session.prototype.sync = function() {
                    var $this = this,
                        d = $q.defer();
                    $this.resource.get({}, function(result) {
                        $this.login(result);
                        d.resolve(result);
                    });
                    return d.promise;
                };

                Session.prototype.redirectToLogin = function() {
                    this.clear();
                    this.setLastLoginUrl();
                    Alert.danger("You must be logged in");
                    $location.path("/login");
                };

                Session.prototype.check = function() {
                    var $this = this;
                    $this.sync().then(function() {
                        if (!$this.loggedIn) {
                            $this.redirectToLogin();
                        }
                    });
                };

                Session.prototype.setLastLoginUrl = function() {
                    this.lastLoginUrl = $location.path();
                };

                Session.prototype.getLastLoginUrl = function() {
                    var url = this.lastLoginUrl;
                    if (isNoRedirectFromLogin(url)) {
                        url = null;
                    }
                    this.lastLoginUrl = null;
                    return url;
                };

                Session.prototype.clear = function() {
                    this.loggedIn = false;
                    this.name = null;
                    this.id = null;
                    this.isAdmin = false;
                };

                Session.prototype.set = function(session) {
                    this.loggedIn = session.loggedIn;
                    this.name = session.name;
                    this.id = session.id;
                    this.isAdmin = session.isAdmin;
                };

                Session.prototype.login = function(result, token) {
                    this.set(result);
                    this.$delete = result.$delete;
                    if (token) {
                        $window.localStorage.setItem(authTokenStorageKey, token);
                    }
                };

                Session.prototype.logout = function() {
                    var $this = this,
                        d = $q.defer();
                    $this.$delete(function(result) {
                        $this.clear();
                        d.resolve(result);
                        $window.localStorage.removeItem(authTokenStorageKey);
                    });
                    return d.promise;
                };

                return new Session();

            }
        ])
        .service('Auth', ['$resource',
            'urls',
            function($resource, urls) {
                return $resource(urls.auth, {}, {
                    'signup': {
                        method: 'POST',
                        url: urls.auth + 'signup'
                    },
                    'recoverPassword': {
                        method: 'PUT',
                        url: urls.auth + 'recoverpass'
                    },
                    'changePassword': {
                        method: 'PUT',
                        url: urls.auth + 'changepass'
                    }
                });
            }
        ])
        .service('Photo', ['$resource', 'urls',
            function($resource, urls) {
                return $resource(urls.photos, {
                    id: '@id'
                }, {
                    'query': {
                        method: 'GET',
                        isArray: false
                    },
                    'search': {
                        method: 'GET',
                        isArray: false,
                        url: urls.photos + '/search'
                    },
                    'byOwner': {
                        method: 'GET',
                        isArray: false,
                        url: urls.photos + '/owner/:ownerID'
                    },
                    'updateTitle': {
                        method: 'PATCH',
                        url: urls.photos + '/title'
                    },
                    'updateTags': {
                        method: 'PATCH',
                        url: urls.photos + '/tags'
                    },
                    'upvote': {
                        method: 'PATCH',
                        url: urls.photos + '/upvote'
                    },
                    'downvote': {
                        method: 'PATCH',
                        url: urls.photos + '/downvote'
                    },
                });
            }
        ])
        .service('Tag', ['$resource', 'urls',
            function($resource, urls) {
                return $resource(urls.tags);
            }
        ])
        .service('Alert', [

            function() {

                function Alert() {
                    var $this = this;
                    $this.message = null;

                    var addMessage = function(level, message) {
                        $this.message = {
                            message: message,
                            level: level
                        };
                    };

                    $this.dismiss = function() {
                        $this.message = null;
                    };

                    $this.success = addMessage.bind(null, "success");
                    $this.info = addMessage.bind(null, "info");
                    $this.warning = addMessage.bind(null, "warning");
                    $this.danger = addMessage.bind(null, "danger");

                }
                return new Alert();

            }
        ]);
})();
