'use strict';

/* Controllers */

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

]

angular.module('photoshare.controllers', [])
    .controller('ListCtrl', ['$scope', function ($scope) {
        $scope.photos = photos;
    }])
    .controller('UploadCtrl', ['$scope', '$location', function ($scope, $location) {
        $scope.newPhoto = {};
        $scope.uploadPhoto = function (){
            $scope.newPhoto.thumbnail = 'http://placeimg.com/200/200/nature';
            photos.splice(0, 0, $scope.newPhoto);
            $scope.newPhoto = {};
            $location.path("/list");
        };

    }]);
