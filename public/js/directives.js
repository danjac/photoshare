/* Directives */

(function() {
    'use strict';
    angular.module('photoshare.directives', []).
    directive('filesModel', function() {
        /* https://github.com/angular/angular.js/issues/1375#issuecomment-21933012 */
        return {
            controller: function($parse, $element, $attrs, $scope, $window) {
                var exp = $parse($attrs.filesModel);
                $element.on('change', function() {
                    exp.assign($scope, this.files);
                    if ($window.FileReader !== null) {
                        var file = this.files[0],
                            reader = new $window.FileReader();
                        reader.onload = function() {
                            $scope.upload = {
                                url: reader.result
                            };
                            $scope.$apply();
                        };
                        reader.readAsDataURL(file);
                    }

                    $scope.$apply();
                });
            }
        };
    }).
    directive('pagination', function() {

        return {
            restrict: 'E',
            replace: true,
            link: function(scope, element, attrs) {

                scope.isFirstPage = false;
                scope.isLastPage = true;

                scope.$watch('currentPage', function(page) {
                    scope.isFirstPage = (scope.currentPage == 1);
                    scope.isLastPage = (scope.currentPage == scope.numPages);
                    scope.pageRange = [];
                    for (var i = 0; i < scope.numPages; i++) {
                        scope.pageRange.push(i + 1);
                    }
                });

                scope.nextPage = function(page) {
                    scope.onNextPage(page);
                };
            },
            scope: {
                numPages: '=',
                currentPage: '=',
                onNextPage: '='
            },
            templateUrl: 'partials/pagination.html'

        }
    });
})();