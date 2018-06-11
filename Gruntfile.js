module.exports = function(grunt) {
    require("matchdep").filterAll("grunt-*").forEach(grunt.loadNpmTasks);
    var webpack = require("webpack");
    var webpackConfig = require("./webpack.config.js")
    grunt.initConfig({
        webpack: {
            dist: webpackConfig
        },
    });

    grunt.registerTask('default', ['webpack'])
};
