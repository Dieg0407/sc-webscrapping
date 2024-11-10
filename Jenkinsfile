pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                echo 'Building webscrapper'
                sh 'go build -o sc-scrubber cmd/cli.go'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Building webscrapper'
                sh 'mv sc-scrubber /opt/custom/sc-scrubber'
            }
        }
    }
}
