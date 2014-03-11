module.exports = function(grunt) {
  grunt.initConfig({
    pkg: grunt.file.readJSON('package.json'),
    build_backend: {
      dev: {
        dest: 'build/backofficed'
      },
      release: {
        dest: 'build/backofficed-linux-amd64',
        opts: {
          env: {
            // For cross-compile to work on Mac OS do the following:
            // cd /usr/local/go/src && \
            // sudo GOOS=linux GOARCH=amd64 CGO_ENABLED=0 ./make.bash --no-clean
            'PATH': process.env.PATH,
            'GOPATH': process.env.GOPATH,
            'CGO_ENABLED': '0',
            'GOOS': 'linux',
            'GOARCH': 'amd64'
          }
        }
      }
    }
  });

  grunt.registerMultiTask('build_backend', "Build backend server", function() {
    var done = this.async();
    var proc = {
      cmd: "go",
      args: ["build", "-o", this.data.dest, "main.go"],
      opts: this.data.opts
    }
    return grunt.util.spawn(proc, function(error, result){
      if (error) {
        grunt.log.error(
          "Failed to build backoffice backend: " + error + " (" + result.code + ")\n\n" +
          result.stdout + "\n\n" + result.stderr);
      }
      done(!error);
    });
  });

  // grunt.registerTask('build', 'Build complete dev version', ['build:backend']);
}
