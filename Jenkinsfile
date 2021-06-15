node {
  try {
    docker.image("golang:1.16").inside {
      stage("init") {
        checkout scm
      }
      withEnv(['GOCACHE=/go/.cache', 'GOOS=linux']) {
        stage("build") {
          def archs=["arm", "amd64", "mipsle"]
          archs.each {
            sh "GOARCH=${it} ./build.sh"
          }
        }
      }
      stage("package") {
        archiveArtifacts artifacts: "out/librespeed-cli-*", followSymlinks: false
      }
    }
  } finally {
    cleanWs()
  }
}
