var bowerFiles = require('main-bower-files'),
    gulp = require('gulp'),
    path = require('path'),
    debug = require('gulp-debug'),
    _ = require('lodash'),
    util = require('gulp-util'),
    plumber = require('gulp-plumber'),
    shell = require('gulp-shell'),
    minifyCss = require('gulp-minify-css'),
    filter = require('gulp-filter'),
    replace = require('gulp-replace'),
    concat = require('gulp-concat'),
    uglify = require('gulp-uglify'),
    webpack = require('webpack'),
    WebpackDevServer = require('webpack-dev-server'),
    webpackConfig = require('./webpack.config.js');


var staticDir = path.join(__dirname, 'public');

var dest = {
  js: path.join(staticDir, 'js'),
  img: path.join(staticDir, 'img'),
  css: path.join(staticDir, 'css'),
  fonts: path.join(staticDir, 'fonts')
};

gulp.task('assets-dev', function() {
  return assets(true);
});

gulp.task('assets-build', function() {
  return assets(false);
});

function assets(useDevServer) {
  var assetsDir = path.join(__dirname, 'assets', '**/*'),
      cssFilter = filter('**/*.css', {restore: true}),
      imgFilter = filter(['**/*.png', '**/*.gif', '**/*.jpg'], {restore:true}),
      htmlFilter = filter('**/*.html', {restore:true});

  var appJs = useDevServer ? 'http://localhost:8090/js/app.js' : '/js/app.js';

  return gulp.src(assetsDir)
  .pipe(debug())
  .pipe(plumber())
  .pipe(cssFilter)
  .pipe(minifyCss())
  .pipe(concat('app.css'))
  .pipe(gulp.dest(dest.css))
  .pipe(cssFilter.restore)
  .pipe(imgFilter)
  .pipe(gulp.dest(staticDir))
  .pipe(imgFilter.restore)
  .pipe(htmlFilter)
  .pipe(replace('[APP_JS]', appJs))
  .pipe(gulp.dest(staticDir))
  .pipe(htmlFilter.restore);

}


gulp.task('bower', function() {

  var jsFilter = filter('**/*.js', {restore: true}),
      cssFilter = filter('**/*.css', {restore: true}),
      fontFilter = filter(['**/*.eot', '**/*.woff', '**/*.svg', '**/*.ttf'], {restore: true});

  // installs all the 3rd party stuffs.
  return gulp.src(bowerFiles({
    debugging: true,
    checkExistence: true,
    base: 'bower_components/**/*'
  }))
  .pipe(debug())
  .pipe(plumber())
  .pipe(jsFilter)
  .pipe(uglify())
  .pipe(concat('vendor.js'))
  .pipe(gulp.dest(dest.js))
  .pipe(jsFilter.restore)
  .pipe(cssFilter)
  .pipe(minifyCss())
  .pipe(concat('vendor.css'))
  .pipe(gulp.dest(dest.css))
  .pipe(cssFilter.restore)
  .pipe(fontFilter)
  .pipe(gulp.dest(dest.fonts))
  .pipe(fontFilter.restore);
});

gulp.task("build", ["install", "bower", "assets-build"], function(callback) {
  var webpackBuildOptions = _.create(webpackConfig, {
    debug: false,
    verbose: false,
    devServer: false,
    devtool: 'eval',
    entry: ['./ui/app.js'],
    plugins: [
        new webpack.optimize.UglifyJsPlugin({
            warnings: false
        })
    ]
  });

  webpack(webpackBuildOptions, function(err, stats) {
    if (err) {
        throw new util.PluginError("build", err);
    }
    util.log("build", stats.toString());
    callback();
  });
});


gulp.task("webpack-dev-server", function(callback) {
    new WebpackDevServer(webpack(webpackConfig), {
        publicPath: webpackConfig.output.publicPath,
        hot: true,
        quiet: false,
        lazy: false,
        watchOptions: {
          aggregateTimeout: 300
        },
        stats: { colors: true },
        historyApiFallback: true
    }).listen(8090, 'localhost', function (err, result) {
        if (err) {
            console.log(err);
        }
        console.log('Listening at localhost:8090');
    });
});

gulp.task('install', shell.task([
    'bower cache clean',
    'bower install'
]));

gulp.task('default', ['install', 'bower', 'assets-dev', 'webpack-dev-server']);
