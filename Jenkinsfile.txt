pipeline {
    agent any
    tools {
        go 'go1.14'
    }
    stages {        
        stage('Pre Test') {
            steps {
                echo 'Installing dependencies'
                sh 'ls -la'
                sh 'pwd'
                sh 'go version'
            }
        }
        
        stage('Execute') {
            steps {
                echo 'Compiling and building'
                sh 'go run main.go'
            }
        }  
}
}