'use strict';

/* Services */

angular.module('photoshare.services', [])
    .service('Authenticator', ['$resource', 'urls', function ($resource, urls) {

        function Authenticator() {
            this.loggedIn = false;
            this.currentUser = null;
            this.resource = $resource(urls.auth);
        }

        Authenticator.prototype.canDelete = function (photo) {

            if (!this.currentUser) {
                return false;
            }
            return this.canEdit(photo) || this.currentUser.isAdmin;
        };
        
        Authenticator.prototype.canEdit = function (photo) {
            if (!this.currentUser) {
                return false;
            }
            return photo.ownerId === this.currentUser.id;
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

    
