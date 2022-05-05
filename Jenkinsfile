node {
  try {
    docker.image("golang:1.16").inside {
      stage("init") {
        checkout scm
      }
      withEnv(['GOCACHE=/go/.cache', 'GOOS=linux']) {
        stage("build") {
          def archs=["arm", "arm64", "amd64", "mipsle"]
          archs.each {
            sh "GOARCH=${it} ./build.sh"
          }
        }
      }
      stage("package") {
        def archMap = [arm: "armv7", arm64: "aarch64", mipsle: "mips32", amd64: "amd64"]
        archMap.each(entry -> {
          sh "mv out/librespeed-cli-linux-${entry.key} out/librespeed-cli-${entry.value}"
        })
        archiveArtifacts artifacts: "out/librespeed-cli-*", followSymlinks: false
      }
    }
  } finally {
    cleanWs()
  }
}
