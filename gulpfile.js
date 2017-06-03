const browserSync = require('browser-sync').create();
const concat = require('gulp-concat');
const gulp = require('gulp');
const imagemin = require('gulp-imagemin');
const postcss = require('gulp-postcss');
const sass = require('gulp-sass');
const uglify = require('gulp-uglify');

gulp.task('default', ['scripts', 'styles', 'images']);

gulp.task('scripts', function() {
    return gulp.src([
        './node_modules/jquery/dist/jquery.js',
        './node_modules/what-input/what-input.js',
        './node_modules/foundation-sites/dist/js/foundation.js',
        './node_modules/datatables.net/js/jquery.dataTables.js',
        './node_modules/datatables.net-zf/js/dataTables.foundation.js',
        './node_modules/chart.js/dist/Chart.bundle.js',
        './node_modules/blazy/blazy.js',
        './js/main.js'
    ])
    .pipe(concat('main.js'))
    .pipe(uglify())
    .pipe(gulp.dest('./public/js'))
});

gulp.task('styles', function() {
    return gulp.src([
        './css/main.scss',
        './node_modules/datatables.net-zf/css/dataTables.foundation.css'
    ])
    .pipe(sass({
        includePaths: [
            'node_modules/foundation-sites/scss'
        ]
    }).on('error', sass.logError))
    .pipe(concat('main.css'))
    .pipe(postcss())
    .pipe(gulp.dest('./public/css'))
    .pipe(browserSync.stream());
});

gulp.task('images', function() {
    return gulp.src([
        './node_modules/datatables.net-zf/images/*'
    ])
    .pipe(imagemin())
    .pipe(gulp.dest('./public/images'))
});

gulp.task('watch', ['default'], function() {
    browserSync.init({
        proxy: 'http://localhost:4000',
        open: false
    });

    gulp.watch('css/*.scss', ['styles']);
    gulp.watch('js/*.js', ['scripts'])
    gulp.watch('public/js/*.js').on('change', browserSync.reload);
    gulp.watch('templates/**/*.html').on('change', browserSync.reload);
});
