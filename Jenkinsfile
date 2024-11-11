pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                echo 'Building webscrapper'
                sh 'go build -o scrubber cmd/cli.go'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Building webscrapper'
                sh 'mv scrubber /opt/custom/scrubber'
                sh 'cp runner.sh /opt/custom/runner.sh'
            }
        }
    }
}
