'use strict';

/* Services */

angular.module('photoshare.services', [])
    .service('Authenticator', ['$resource', '$q', function ($resource, $q) {

        var AuthResource = $resource("/auth/");

        function AuthService() {
            this.currentUser = null;
            this.loggedIn = false;
        }

        AuthService.prototype.authenticate = function () {
            var deferred = $q.defer(), $this = this;
            AuthResource.get({}, function (user) {
                $this.currentUser = user;
                $this.loggedIn = true;
                deferred.resolve($this.currentUser);
            });
            return deferred.promise;
        };

        AuthService.prototype.isLoggedIn = function () {
            return this.loggedIn;
        };

        AuthService.prototype.login = function (email, password) {
            var deferred = $q.defer(), $this = this;
            this.currentUser = new AuthResource({
                email: email,
                password: password
            });
            this.currentUser.$save(function () {
                $this.loggedIn = true;
                deferred.resolve($this.currentUser);
            });
            return deferred.promise;
        };

        AuthService.prototype.logout = function () {
            this.currentUser.$delete();
            this.loggedIn = false;
        };

        return new AuthService();
    }])
    .service('Photo', ['$resource', function ($resource) {
        return $resource("/photos/");
    }]);
    
