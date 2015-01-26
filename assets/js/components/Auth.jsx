var UserStore = require('../stores/UserStore');
var Login = require('./Login.jsx');

var Auth = {
  statics: {
    willTransitionTo: function (transition) {
      if (!UserStore.isLoggedIn()) {
        Login.attemptedTransition = transition;
        transition.redirect('/login');
      }
    }
  }
};

module.exports = Auth;
