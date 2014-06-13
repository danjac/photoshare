'use strict';

/* Services */

angular.module('photoshare.services', [])
    .service('Authenticator', ['$resource', '$q', 'urls', function ($resource, $q, urls) {

        function Authenticator() {
            this.resource = $resource(urls.auth);
            this.session = null;
        }

        Authenticator.prototype.init = function () {
            var $this = this;
            $this.resource.get({}, function (result) {
                $this.session = result;
            });
        };
        
        Authenticator.prototype.logout = function () {
            var $this = this, d = $q.defer();
            $this.session.$delete(function (result) {
                $this.session = result;
                d.resolve($this.session);
            });
            return d.promise;
        };

        Authenticator.prototype.canDelete = function (photo) {

            if (!this.session) {
                return false;
            }
            return this.canEdit(photo) || this.session.isAdmin;
        };
        
        Authenticator.prototype.canEdit = function (photo) {
            if (!this.session) {
                return false;
            }
            return photo.ownerId === this.session.id;
        };

        return new Authenticator();

    }])
    .service('Photo', ['$resource', 'urls', function ($resource, urls) {
        return $resource(urls.photos, {id: '@id'}, { 'update': { method: 'PUT' } });
    }])
    .service('User', ['$resource', 'urls', function ($resource, urls) {
        return $resource(urls.users);
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

    
