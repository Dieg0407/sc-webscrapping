pipeline {
    agent any

    stages {
        stage('Build CLI Tool') {
            steps {
                dir("cli") {
                    echo 'Building scrubber'
                    sh 'go build -o scrubber cmd/cli.go'

                    echo 'Updating scrubber'
                    sh 'mv scrubber /opt/custom/scrubber'
                    sh 'cp scripts/runner.sh /opt/custom/runner.sh'
                }
                
            }
        }
    }
}
