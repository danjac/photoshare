'use strict';

/* jasmine specs for directives go here */

describe('directives', function() {
    beforeEach(module('photoshare.directives'));

    describe('tab', function () {
        it('should show a tab', function ($rootScope, $compile) {
            var element = $compile('<tab url="/">hello</tab>')($rootScope);
            console.log(element);
            expect(element.text()).toEqual("hello");

        });
    });
});
