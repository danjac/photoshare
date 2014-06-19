'use strict';

/* Services */

angular.module('photoshare.services', [])
    .service('Session', ['$location', 'Alert', function ($location, Alert) {

        function Session() {
            this.clear();
            this.lastLoginUrl = null;
        }

        Session.prototype.check = function () {
            if (!this.loggedIn) {
                Alert.danger("You must be logged in");
                this.setLastLoginUrl();
                $location.path("/login");
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

        Session.prototype.canDelete = function (photo) {

            if (!this.loggedIn) {
                return false;
            }
            return this.canEdit(photo) || this.isAdmin;
        };

        Session.prototype.canEdit = function (photo) {
            if (!this.loggedIn) {
                return false;
            }
            return photo.ownerId === this.id;
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
                $window.sessionStorage.token = token;
            }
        };

        Authenticator.prototype.logout = function () {
            var $this = this, d = $q.defer();
            delete $window.sessionStorage.token;
            $this.$delete(function (result) {
                Session.clear();
                d.resolve(result);
            });
            return d.promise;
        };

        return new Authenticator();

    }])
    .service('Photo', ['$resource', 'urls', function ($resource, urls) {
        return $resource(urls.photos, {id: '@id'}, { 'update': { method: 'PUT' } });
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


