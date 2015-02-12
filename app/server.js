/**
 * Module dependencies.
 */

require('node-jsx').install()

var express = require('express');
var http = require('http');
var path = require('path');
var bodyParser = require('body-parser');
var React = require('react');
var Router = require('react-router');
var Routes = require('../assets/js/Routes.jsx');

var app = express();

// all environments
app.set('port', process.env.PORT || 3000);
app.set('views', __dirname + '/views');
app.set('view engine', 'ejs');
//app.use(express.favicon());
//app.use(express.logger('dev'));
app.use(bodyParser());
//app.use(express.methodOverride());
//app.use(express.static(path.join(__dirname, '../public')));

// development only
if ('development' == app.get('env')) {
  //app.use(express.errorHandler());
}

app.post("/react/", function(req, res){

    var props = JSON.parse(req.body.props || "{}");

    Router.run(Routes, req.body.route || "", function(Handler, state) {
        var markup = React.renderToString(Handler(props));
        res.render("index", {
          markup: markup,
          data: req.body.props
        });
    });

});

http.createServer(app).listen(app.get('port'), function(){
  console.log('Express server listening on port ' + app.get('port'));
});
