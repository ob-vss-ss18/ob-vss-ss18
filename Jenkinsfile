pipeline {
    agent any
    stages {
        stage('Build') {
            steps {
                sh 'cd stringutil && go build'
            }
        }
        stage('Test') {
            steps {
                timeout(20) {
                  sh 'go tool vet stringutil'
                  sh 'golint stringutil'
                  sh 'cd stringutil && go test'
                }
            }
        }
    }
}