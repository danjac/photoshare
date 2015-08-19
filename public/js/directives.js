/* Directives */

(function() {
    'use strict';
    angular.module('photoshare.directives', []).
    directive('errSrc', function($interval, $http) {
        return {
            link: function(scope, element, attrs) {
                element.bind('error', function() {
                    var imageFound = false,
                        originalSrc = attrs.src;
                    if (attrs.src != attrs.errSrc) {
                        attrs.$set('src', attrs.errSrc);
                    }
                    $interval(function() {
                        if (imageFound) return;
                        $http.get(originalSrc).success(function(response) {
                            attrs.$set('src', originalSrc);
                            imageFound = true;
                        });
                    }, 1000, 10);

                });
            }
        };
    }).
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
    });
})();
