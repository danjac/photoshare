var Dispatcher = require('flux').Dispatcher;
var Constants = require('./Constants').ActionSources;

var AppDispatcher = new Dispatcher();

AppDispatcher.handleViewAction = function(action) {

    action.source = Constants.VIEW_ACTION;
    this.dispatch(action);

};

AppDispatcher.handleServerAction = function(action) {

    action.source = Constants.SERVER_ACTION;
    this.dispatch(action);
};

module.exports = AppDispatcher;
