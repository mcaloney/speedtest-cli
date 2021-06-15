node {
  try {
    docker.image("golang:1.16").inside {
      stage("init") {
        checkout scm
      }
      stage("build-arm") {
        sh "GOOS=linux GOARCH=arm ./build.sh"
      }
      stage("build-amd64") {
        sh "GOOS=linux GOARCH=amd64 ./build.sh"
      }
      stage("package") {
        archiveArtifacts artifacts: "out/librespeed-cli-*", followSymlinks: false
      }
    }
  } finally {
    cleanWs()
  }
}
