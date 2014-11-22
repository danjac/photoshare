var React = require('react'),
    Router = require('react-router'),
    _ = require('lodash'),
    Route = Router.Route,
    DefaultRoute = Router.DefaultRoute,
    RouteHandler = Router.RouteHandler;

  var Alerts = React.createClass({

    getInitialState: function () {
      return {
        messages: []
      };
    },

    render: function () {
      return (
        <div className="alerts">
        {_.map(this.state.messages, function(msg, index) {
            return(
              <div role="alert" className={"alert alert-" + msg.priority} key={index}>
                {msg.message}
              </div>
            )
        })}
        </div>
      )

    }
  });

  var NavBar = React.createClass({
    render: function() {
      return (
        <nav className="navbar navbar-inverse navbar-fixed-top" role="navigation">
          <div className="container-fluid">
              <div className="navbar-header">
                  <button type="button" className="navbar-toggle" data-toggle="collapse" data-target="#navbar-links">
                      <span className="sr-only">Toggle navigation</span>
                      <span className="icon-bar"></span>
                      <span className="icon-bar"></span>
                      <span className="icon-bar"></span>
                  </button>
                  <a className="navbar-brand" href="#"><i className="fa fa-camera"></i> Wallshare</a>
              </div>
              <div className="collapse navbar-collapse" id="navbar-links">
                  <ul className="nav navbar-nav navbar-left">
                  </ul>
                  <ul className="nav navbar-nav navbar-right">
                  </ul>
              </div>
          </div>
        </nav>
      )
    }

  });

  var Photos = React.createClass({
    render: function () {
      return (
        <div>
          Photos go here !!!
        </div>
      )
    }

  });

  var App = React.createClass({

    render: function () {
      return (
        <div>
          <NavBar />
          <div className="container-fluid">
            <RouteHandler />
          </div>
        </div>
      )
    }
  });

  var routes = (
    <Route name="app" path="/" handler={App}>
    {/* <Route name="" handler={} />  */}
      <DefaultRoute handler={Photos} />
    </Route>
  );

  Router.run(routes, function(Handler) {
      React.render(<Handler />, document.body);
  });

/*
(function() {
    'use strict';
    // Declare app level module which depends on filters, and services
    angular.module('photoshare', [
        'ngResource',
        'ngAnimate',
        'ngSanitize',
        'ngCookies',
        'ui.router',
        'ui.bootstrap',
        'ui.gravatar',
        'photoshare.filters',
        'photoshare.services',
        'photoshare.directives',
        'photoshare.controllers'
    ]).
    constant('urls', {
        auth: '/api/auth/',
        photos: '/api/photos/:id',
        tags: '/api/tags/',
        messages: '/api/messages'
    }).
    constant('authTokenHeader', 'X-Auth-Token').
    constant('authTokenStorageKey', 'authToken').
    config(['$stateProvider',
        '$urlRouterProvider',
        '$httpProvider',
        '$resourceProvider',
        function(
            $stateProvider,
            $urlRouterProvider,
            $httpProvider,
            $resourceProvider
        ) {

            $urlRouterProvider.otherwise("/");

            $stateProvider.
            state('front', {
                url: '/',
                templateUrl: 'partials/front.html',
                controller: 'FrontCtrl'
            }).

            state('popular', {
                url: '/popular',
                templateUrl: 'partials/list.html',
                controller: 'ListCtrl',
                data: {
                    orderBy: "votes"
                }
            }).

            state('latest', {
                url: '/latest',
                templateUrl: 'partials/list.html',
                controller: 'ListCtrl'
            }).

            state('tags', {
                url: '/tags',
                templateUrl: 'partials/tags.html',
                controller: 'TagsCtrl'
            }).

            state('tag', {
                url: '/tag/:tag',
                templateUrl: 'partials/list.html',
                controller: 'ListCtrl'
            }).

            state('search', {
                url: '/search/:q',
                templateUrl: 'partials/list.html',
                controller: 'ListCtrl'
            }).

            state('owner', {
                url: '/owner/:ownerID/:ownerName',
                templateUrl: 'partials/list.html',
                controller: 'ListCtrl'
            }).

            state('detail', {
                url: '/detail/:id',
                templateUrl: 'partials/detail.html',
                controller: 'DetailCtrl'
            }).

            state('upload', {
                url: '/upload',
                templateUrl: 'partials/upload.html',
                controller: 'UploadCtrl',
                data: {
                    loginRequired: true
                }
            }).

            state('login', {
                url: '/login',
                templateUrl: 'partials/login.html',
                controller: 'LoginCtrl'
            }).

            state('recoverpass', {
                url: '/recoverpass',
                templateUrl: 'partials/recover_pass.html',
                controller: 'RecoverPassCtrl'
            }).

            state('changepass', {
                url: '/changepass',
                templateUrl: 'partials/change_pass.html',
                controller: 'ChangePassCtrl'
            }).

            state('signup', {
                url: '/signup',
                templateUrl: 'partials/signup.html',
                controller: 'SignupCtrl'
            });

            //$locationProvider.html5Mode(true);
            //
            $resourceProvider.defaults.stripTrailingSlashes = false;

            //$httpProvider.defaults.xsrfCookieName = "csrf_token";
            //$httpProvider.defaults.xsrfHeaderName = "X-CSRF-Token";

            // handle file uploads

            $httpProvider.defaults.transformRequest = function(data, headersGetter) {

                if (data === undefined) {
                    return data;
                }

                var fd = new FormData(),
                    isFileUpload = false,
                    headers = headersGetter();

                angular.forEach(data, function(value, key) {
                    if (value instanceof FileList) {
                        isFileUpload = true;
                        if (value.length === 1) {
                            fd.append(key, value[0]);
                        } else {
                            angular.forEach(value, function(file, index) {
                                fd.append(key + "_" + index, file);
                            });
                        }
                    } else {
                        fd.append(key, value);
                    }
                });
                if (isFileUpload) {
                    headers["Content-Type"] = undefined;
                    return fd;
                }

                return JSON.stringify(data);
            };

            var interceptors = ['AuthInterceptor', 'ErrorInterceptor'];

            angular.forEach(interceptors, function(interceptor) {
                $httpProvider.interceptors.push([
                    '$injector',
                    function($injector) {
                        return $injector.get(interceptor);
                    }
                ]);
            });

        }
    ]).factory('AuthInterceptor', function($window, $cookies, authTokenHeader, authTokenStorageKey) {

        return {
            request: function(config) {

                // oauth2 authentication sets token in cookie before redirect
                // check the existence of cookie, add to local storage, and delete the cookie.
                if ($cookies.authToken) {
                    $window.localStorage.setItem(authTokenStorageKey, $cookies.authToken);
                    delete $cookies.authToken;
                }

                config.headers = config.headers || {};

                var token = $window.localStorage.getItem(authTokenStorageKey);

                if (token) {
                    config.headers[authTokenHeader] = token;
                }
                return config;
            }
        };

    }).factory('ErrorInterceptor', function($q, $location, Session, Alert) {
        return {

            response: function(response) {
                return response;
            },

            responseError: function(response) {
                var rejection = $q.reject(response),
                    status = response.status,
                    msg = 'Sorry, an error has occurred';

                if (status == 401) {
                    Session.redirectToLogin();
                    return;
                }
                if (status == 404) {
                    // handle locally
                    return;
                }
                if (status == 403) {
                    msg = "Sorry, you're not allowed to do this";
                }
                if (status == 400 && response.data.errors) {
                    msg = "Sorry, your form contains errors, please try again";
                }
                if (status == 413) {
                    msg = "The file was too large!";
                }
                if (response.data && typeof(response.data) === 'string') {
                    msg = response.data;
                }
                if (msg) {
                    Alert.danger(msg);
                }
                return rejection;
            }
        };
    });
})();
*/
