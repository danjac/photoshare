'use strict';

/* jasmine specs for controllers go here */

describe('controllers', function (){


    beforeEach(function () {
        module('photoshare');
        module('photoshare.services');
        module('photoshare.controllers');
    });

    it('should show a list of photos', inject(function ($rootScope, $controller, _$httpBackend_) {

        var scope = $rootScope.$new(),
            httpBackend = _$httpBackend_;

        httpBackend.expectGET("/api/photos/?page=1").respond([
            {
                'title': 'this is a photo',
                'photo': 'test.jpg'
            }
        ]);
        var listCtrl = $controller('ListCtrl', { 
            $scope: scope,
        });
        scope.nextPage();
        httpBackend.flush();
        expect(scope.photos.length).toBe(1);
    }));

    it('should show upload form', inject(function ($location, $rootScope, $controller, _$httpBackend_, Authenticator) {
        var scope = $rootScope.$new(), 
            httpBackend = _$httpBackend_,
            data = {
                title: 'test',
                photo: 'test.png'
            };

        
        httpBackend.expectPOST("/api/photos/", data).respond({
            'id': 1,
            'title': 'test'
        });

        Authenticator.currentUser = {id : 1};
        Authenticator.loggedIn = true;

        $controller('UploadCtrl', { $scope: scope });

        scope.newPhoto.title = data.title;
        scope.newPhoto.photo = data.photo;
        scope.uploadPhoto();

        httpBackend.flush();

        expect(scope.newPhoto.title).toBe(undefined);
        expect($location.path()).toBe("/#/list");

    }));
});
