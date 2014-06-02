'use strict';

/* Directives */

angular.module('photoshare.directives', []).
    directive('tab', function () {

        function isActive(url, current) {
            return current.indexOf(url, current.length - url.length) !== -1;
        }

        return {
            restrict: 'E',
            replace: true,
            transclude: true,
            link: function ($scope, element, attrs) {
                $scope.url = attrs.url;
                $scope.$on('$locationChangeStart', function (next, current) {
                    $scope.active = isActive($scope.url, current);
                });
            },
            scope: {
                url: '@'
            },
            template: '<li ng-class="{\'active\': active}"><a ng-transclude href="{{url}}"></a></li>'
        };
    });
