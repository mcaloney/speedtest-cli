node {
  try {
    docker.image("golang:1.16").inside {
      stage("init") {
        checkout scm
      }
      withEnv(['GOCACHE=/go/.cache', 'GOOS=linux']) {
        stage("build-arm") {
          sh "GOARCH=arm ./build.sh"
        }
        stage("build-amd64") {
          sh "GOARCH=amd64 ./build.sh"
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
