module.exports = function(grunt) {
  grunt.initConfig({
    watch: {
      files: ['src/module1/static/js/*.js',
              'src/module1/static/js/*/*.js',
              'src/module1/static/*.html'],
      options: {
        livereload: true
      }
    }
  });

  grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.registerTask('default', ['watch']);

};