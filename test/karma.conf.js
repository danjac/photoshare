'use strict';

module.exports = function(config){
    config.set({

        basePath : '../',

        files : [
            'public/bower_components/angular/angular.js',
            'public/bower_components/angular-route/angular-route.js',
            'public/bower_components/angular-resource/angular-resource.js',
            'public/bower_components/nginfinitescroll/build/ng-infinite-scroll.js',
            'public/bower_components/angular-mocks/angular-mocks.js',
            'public/js/**/*.js',
            'test/unit/**/*.js'
        ],

        autoWatch : true,

        frameworks: ['jasmine'],

        browsers : ['Chrome'],

        plugins : [
            'karma-chrome-launcher',
            'karma-firefox-launcher',
            'karma-jasmine',
            'karma-junit-reporter'
        ],

        junitReporter : {
            outputFile: 'test_out/unit.xml',
            suite: 'unit'
        }

    });
};
