'use strict';

/* Services */

angular.module('photoshare.services', [])
    .service('Authenticator', ['$http', '$q', function ($http, $q) {

        function AuthService() {
            this.currentUser = null;
            this.loggedIn = false;
        }

        AuthService.prototype.authenticate = function () {
            var deferred = $q.defer(), $this = this;
            $http.get("/auth").then(function (response) {
                $this.currentUser = response.data;
                $this.loggedIn = true;
                deferred.resolve($this.currentUser);
            });
            return deferred.promise;
        };

        AuthService.prototype.isLoggedIn = function () {
            return this.currentUser !== null;
        };

        AuthService.prototype.login = function (email, password) {
            var deferred = $q.defer(), $this = this;
            $http.post("/login", {email: email, password: password}).then(function (response) {
                $this.currentUser = response.data;
                $this.loggedIn = true;
                deferred.resolve($this.currentUser);
            });
            return deferred.promise;
        };

        AuthService.prototype.logout = function () {
            this.currentUser = null;
            this.loggedIn = false;
            $http.post("/logout");
        };

        return new AuthService();
    }])
    .service('Photo', ['$q', '$http', function ($q, $http) {

        var getPhotos = function () {
            var deferred = $q.defer();
            $http.get("/photos").then(function (response) {
                deferred.resolve(response.data);
            });
            return deferred.promise;
        };

        return {
            query: getPhotos
        };

    }]);
    
