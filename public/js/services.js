'use strict';

/* Services */

angular.module('photoshare.services', [])
    .service('MessageQueue', ['$window', 'Alert', 'Session', function ($window, Alert, Session) {

    // options usage example
        var options = {
            debug: true,
            devel: true,
            protocols_whitelist: ['websocket', 'xdr-streaming', 'xhr-streaming', 'iframe-eventsource', 'iframe-htmlfile', 'xdr-polling', 'xhr-polling', 'iframe-xhr-polling', 'jsonp-polling']
        };
        function Mq() {
            this.socket = null;
        }

        Mq.prototype.init = function () {
            this.socket = new $window.SockJS('/api/messages', undefined, options);
            this.socket.onopen = function (e) {
                console.log(e)
            };
            this.socket.onmessage = function (e) {
                var msg = JSON.parse(e.data), content=null, photoLink = null;
                if (msg.username === Session.name) {
                    return;
                }
                if (msg.photoID) {
                    photoLink = '<a href="/#/detail/' + msg.photoID + '">a photo</a>';
                }
                switch (msg.type) {
                    case 'login':
                    content = msg.username + " has logged in";
                    break;
                    case 'logout':
                    content = msg.username + " has logged out";
                    break
                    case 'photo_deleted':
                    content = msg.username + " has deleted a photo";
                    break;
                    case 'photo_updated':
                    content = msg.username + " has updated " + photoLink;
                    break;
                    case 'photo_uploaded':
                    content = msg.username + " has uploaded " + photoLink;
                    break;
                }
                if (content) {
                    Alert.success(content);
                }
            };
        };

        return new Mq();
    }])
    .service('Session', ['$location', 'Alert', function ($location, Alert) {

        function Session() {
            this.clear();
            this.lastLoginUrl = null;
        }

        Session.prototype.redirectToLogin = function () {
            this.clear();
            this.setLastLoginUrl();
            Alert.danger("You must be logged in");
            $location.path("/login");
        }

        Session.prototype.check = function () {
            if (!this.loggedIn) {
                this.redirectToLogin();
                return false;
                }
            return true;
        }

        Session.prototype.setLastLoginUrl = function () {
            this.lastLoginUrl = $location.path();
        };

        Session.prototype.getLastLoginUrl = function () {
            var url = this.lastLoginUrl;
            this.lastLoginUrl = null;
            return url;
        };

        Session.prototype.clear = function () {
            this.loggedIn = false;
            this.name = null;
            this.id = null;
            this.isAdmin = false;
        };

        Session.prototype.set = function (session) {
            this.loggedIn = session.loggedIn;
            this.name = session.name;
            this.id = session.id;
            this.isAdmin = session.isAdmin;
        };

        return new Session();

    }])
    .service('Authenticator', ['$resource',
                               '$q',
                               '$window',
                               'urls',
                               'Session', function ($resource, $q, $window, urls, Session) {

        function Authenticator() {
            this.resource = $resource(urls.auth);
        }

        Authenticator.prototype.init = function () {
            var $this = this;
            $this.resource.get({}, function (result) {
                $this.login(result);
            });
        };

        Authenticator.prototype.login = function (result, token) {
            Session.set(result);
            this.$delete = result.$delete;
            if (token) {
                $window.localStorage.setItem("authToken", token)
            }
        };

        Authenticator.prototype.logout = function () {
            var $this = this, d = $q.defer();
            $this.$delete(function (result) {
                Session.clear();
                d.resolve(result);
                $window.localStorage.removeItem("authToken")
            });
            return d.promise;
        };

        return new Authenticator();

    }])
    .service('Photo', ['$resource', 'urls', function ($resource, urls) {
        return $resource(urls.photos, {id: '@id'}, {
            'query': { method: 'GET', isArray: false },
            'search': { method: 'GET', isArray: false, url: urls.photos + '/search' },
            'byOwner': { method: 'GET', isArray: false, url: urls.photos + '/owner/:ownerID' }   ,
            'updateTitle': { method: 'PATCH', url: urls.photos + '/title' },
            'updateTags': { method: 'PATCH', url: urls.photos + '/tags' },
            'upvote': { method: 'PATCH', url: urls.photos + '/upvote'},
            'downvote': { method: 'PATCH', url: urls.photos + '/downvote'},
        });
    }])
    .service('User', ['$resource', 'urls', function ($resource, urls) {
        return $resource(urls.users);
    }])
    .service('Tag', ['$resource', 'urls', function ($resource, urls) {
        return $resource(urls.tags);
    }])
    .service('Alert', [function () {

        function Alert() {
            this.message = null;
            this.addMessage = function (message, level) {
                this.message = {message: message, level: level};
            };
            this.dismiss = function () { this.message = null; };

            this.success = function (message) { this.addMessage(message, "success"); };
            this.info = function (message) { this.addMessage(message, "info"); };
            this.warning = function (message) { this.addMessage(message, "warning"); };
            this.danger = function (message) { this.addMessage(message, "danger"); };
        }

        return new Alert();

    }]);


