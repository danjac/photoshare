'use strict';

/* jasmine specs for directives go here */

describe('directives', function() {
    beforeEach(module('photoshare.directives'));

    describe('navtab', function () {
        it('should show a tab', inject(function ($rootScope, $compile) {
            var element = $compile('<navtab url="/">hello</navtab>')($rootScope);
            expect(element.text()).toEqual("hello");

        }));
    });
});
