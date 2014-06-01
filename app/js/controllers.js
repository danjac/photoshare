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

    ];

angular.module('photoshare.controllers', [])
    .controller('AppCtrl', ['$scope', function ($scope) {
        $scope.tabs = [
            {
                'url': '/list',
                'label': 'Latest',
                'active': true
            },
            {
                'url': '/upload',
                'label': 'Upload a photo',
                'active': false
            }
        ];
        $scope.$on('$locationChangeStart', function (next, current) {
            angular.forEach($scope.tabs, function (tab) {
                if (current.indexOf(tab.url, current.length - tab.url.length) !== -1) {
                    tab.active = true;
                } else {
                    tab.active = false;
                }
            });
        });
    }])
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
