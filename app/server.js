/**
 * Module dependencies.
 */

var express = require('express');
var cors = require('cors');
var http = require('http');
var path = require('path');
var bodyParser = require('body-parser');
var React = require('react');
var API = require('../assets/js/API');

require('node-jsx').install()

var app = express();
app.use(cors());

// all environments
app.set('port', process.env.PORT || 3000);
app.set('views', __dirname + '/views');
app.set('view engine', 'ejs');
//app.use(express.favicon());
//app.use(express.logger('dev'));
app.use(bodyParser());
//app.use(express.methodOverride());
app.use(express.static(path.join(__dirname, '../public')));

// development only
if ('development' == app.get('env')) {
  //app.use(express.errorHandler());
}

// JSX components

var Popular = React.createFactory(require('../assets/js/components/Popular.jsx'));
// latest, search, detail
// we also want to render to JSON

app.get("/", function(req, res){
    API.getPhotos("votes", 1, function(data){
        var markup = React.renderToString(Popular({photos: data}));
        res.render("index", {
          markup: markup,
          //data: JSON.Stringify(data)
        });
    }, true);
});
//app.get('/', routes.index);
//app.get('/users', user.list);

http.createServer(app).listen(app.get('port'), function(){
  console.log('Express server listening on port ' + app.get('port'));
});
