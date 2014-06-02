'use strict';

/* Services */
var photos = [
        {
            'title': 'test',
            'thumbnail': 'http://placeimg.com/200/200/people/sepia'
        },
        {
            'title': 'test',
            'thumbnail': 'http://placeimg.com/200/200/people'
        },
        {
            'title': 'test2',
            'thumbnail': 'http://placeimg.com/200/200/nature/grayscale'
        },
        {
            'title': 'test3',
            'thumbnail': 'http://placeimg.com/200/200/arch'
        },
        {
            'title': 'test3',
            'thumbnail': 'http://placeimg.com/200/200/nature'
        },
        {
            'title': 'test4',
            'thumbnail': 'http://placeimg.com/200/200/tech'
        },
        {
            'title': 'test5',
            'thumbnail': 'http://placeimg.com/200/200/animals'
        },
        {
            'title': 'test6',
            'thumbnail': 'http://placeimg.com/200/200/animals/sepia'
        }

    ];


angular.module('photoshare.services', [])
    .service('Authenticator', [function () {

        function AuthService() {
            this.currentUser = null;
        }

        AuthService.prototype.isLoggedIn = function () {
            return this.currentUser !== null;
        };

        AuthService.prototype.login = function (email, password) {
            this.currentUser = {email: email};
            return this.currentUser;
        };

        AuthService.prototype.logout = function () {
            this.currentUser = null;
        };

        return new AuthService();
    }])
    .service('Photo', [function () {

        var getPhotos = function () {
            return photos;
        };

        return {
            query: getPhotos
        };

    }]);
    
