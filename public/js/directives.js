'use strict';

/* Directives */

angular.module('photoshare.directives', []).
    directive('filesModel', function () {
        /* https://github.com/angular/angular.js/issues/1375#issuecomment-21933012 */
        return {
            controller: function ($parse, $element, $attrs, $scope, $window) {
                var exp = $parse($attrs.filesModel);
                $element.on('change', function () {
                    exp.assign($scope, this.files);
                    if ($window.FileReader !== null) {
                        var file = this.files[0],
                            reader = new $window.FileReader();
                        reader.onload = function () {
                            $scope.upload = {url: reader.result};
                            $scope.$apply();
                        };
                        reader.readAsDataURL(file);
                    }
                    
                    $scope.$apply();
                });
            }
        };
    }).
    directive('navtab', function () {

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
