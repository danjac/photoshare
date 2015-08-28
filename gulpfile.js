var bowerFiles = require('main-bower-files'),
    gulp = require('gulp'),
    _ = require('lodash'),
    util = require('gulp-util'),
    plumber = require('gulp-plumber'),
    shell = require('gulp-shell'),
    minifyCss = require('gulp-minify-css'),
    filter = require('gulp-filter'),
    concat = require('gulp-concat'),
    uglify = require('gulp-uglify'),
    webpack = require('webpack'),
    WebpackDevServer = require('webpack-dev-server'),
    webpackConfig = require('./webpack.config.js');

var staticDir = './public',
    srcDir = './ui',
    cssFilter = filter('*.css', {restore: true}),
    jsFilter = filter('*.js', {restore: true}),
    fontFilter = filter(['*.eot', '*.woff', '*.svg', '*.ttf'], {restore: true});


var dest = {
  js: staticDir + '/js',
  css: staticDir + '/css',
  fonts: staticDir + '/fonts'
};

gulp.task('pkg', function() {
  // installs all the 3rd party stuffs.
  return gulp.src(bowerFiles({
          debugging: true,
          checkExistence: true,
          base: 'bower_components'
  }))
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

gulp.task("build", ["install", "pkg"], function(callback) {
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

gulp.task("watch", ["build"], function() {
  gulp.watch(["ui/**/*"], ["build"]);
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

gulp.task('default', ['install', 'pkg', 'webpack-dev-server']);
