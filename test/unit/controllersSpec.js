'use strict';

/* jasmine specs for controllers go here */

describe('controllers', function (){
    beforeEach(module('photoshare.controllers'));

    it('should change active tab on route change', inject(function ($rootScope, $controller, $location) {
    }));

    it('should show a list of photos', inject(function ($rootScope, $controller, Photo) {
        var scope = $rootScope.$new();
        listCtrl = $controller('ListCtrl', { $scope: scope });
        expect(listCtrl).toBeDefined();
    }));

    it('should show upload form', inject(function ($location, $rootScope, $controller) {
        var scope = $rootScope.$new();
        $controller('UploadCtrl', { $scope: scope });

        scope.newPhoto.title = "this is a new photo";
        scope.newPhoto.file = "photo.jpg";
        scope.uploadPhoto();
        expect(scope.newPhoto.title).toBe(undefined);
        expect($location.path()).toBe("/list");

    }));
});
