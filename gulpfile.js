var browserify = require('gulp-browserify'),
    bowerFiles = require('main-bower-files'),
    watchify = require('watchify'),
    gulp = require('gulp'),
    plumber = require('gulp-plumber'),
    shell = require('gulp-shell'),
    minifyCss = require('gulp-minify-css'),
    gulpFilter = require('gulp-filter'),
    concat = require('gulp-concat');

var staticDir = './public',
    assetsDir = './assets',
    watch = false,
    cssFilter = gulpFilter('*.css'),
    fontFilter = gulpFilter(['*.eot', '*.woff', '*.svg', '*.ttf']);


var src = {
    js: assetsDir + '/js',
    css: assetsDir + '/css',
    fonts: assetsDir + '/fonts'
};

var dest = {
    js: staticDir + '/js',
    css: staticDir + '/css',
    fonts: staticDir + '/fonts'
};


gulp.task('build-js', function() {

    gulp.src(src.js + '/main.js')
        .pipe(browserify({
            cache: {},
            packageCache: {},
            fullPaths: true,
            transform: [
                "reactify",
                "envify"
            ]
        }))
        .pipe(concat('/bundle.js'))
        .pipe(gulp.dest(dest.js));

});

gulp.task('pkg', function() {

    return gulp.src(bowerFiles({
            debugging: true,
            checkExistence: true,
            base: 'bower_components'
        }))
        .pipe(cssFilter)
        .pipe(concat('vendor.css'))
        .pipe(minifyCss())
        .pipe(gulp.dest(dest.css))
        .pipe(cssFilter.restore())
        .pipe(fontFilter)
        .pipe(gulp.dest(dest.fonts));
});

gulp.task('build-css', function() {
    return gulp.src(src.css + '/*.css')
        .pipe(plumber())
        .pipe(concat('app.css'))
        .pipe(minifyCss())
        .pipe(gulp.dest(dest.css));
});


gulp.task('install', shell.task([
    'bower cache clean',
    'bower install'
]));


gulp.task('default', function() {
    gulp.start('install', 'pkg');
    gulp.watch(src.js + '/**', {}, ['build-js']);
    gulp.watch(src.css + '/**', {}, ['build-css']);
});
