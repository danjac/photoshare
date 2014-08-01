'use strict';

module.exports = function(config) {
    config.set({

        basePath: '../',

        files: [
            'public/bower_components/jquery/dist/jquery.min.js',
            'public/bower_components/sockjs/sockjs.min.js',
            'public/bower_components/bootstrap/dist/js/bootstrap.min.js',
            'public/bower_components/angular/angular.js',
            'public/bower_components/angular-cookies/angular-cookies.js',
            'public/bower_components/angular-animate/angular-animate.js',
            'public/bower_components/angular-ui-router/release/angular-ui-router.min.js',
            'public/bower_components/angular-bootstrap/ui-bootstrap.min.js',
            'public/bower_components/angular-bootstrap/ui-bootstrap-tpls.min.js',
            'public/bower_components/angular-sanitize/angular-sanitize.js',
            'public/bower_components/angular-resource/angular-resource.js',
            'public/bower_components/angular-gravatar/build/md5.js',
            'public/bower_components/angular-gravatar/build/angular-gravatar.js',
            'public/bower_components/angular-mocks/angular-mocks.js',
            'public/js/**/*.js',
            'test/unit/**/*.js'
        ],

        autoWatch: true,

        frameworks: ['jasmine'],

        browsers: ['Firefox'],

        plugins: [
            'karma-chrome-launcher',
            'karma-firefox-launcher',
            'karma-jasmine',
            'karma-junit-reporter'
        ],

        junitReporter: {
            outputFile: 'test_out/unit.xml',
            suite: 'unit'
        }

    });
};
