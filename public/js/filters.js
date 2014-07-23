/* Filters */

(function() {
    'use strict';
    angular.module('photoshare.filters', []).
    filter('interpolate', ['version',
        function(version) {
            return function(text) {
                return String(text).replace(/\%VERSION\%/mg, version);
            };
        }
    ]);
})();
