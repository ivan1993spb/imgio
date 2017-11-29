pipeline {
  agent {
    docker {
      image 'golang:alpine'
    }
    
  }
  stages {
    stage('1') {
      steps {
        sleep 1
      }
    }
    stage('2') {
      steps {
        sh 'echo ok'
      }
    }
    stage('3') {
      steps {
        sleep 1
      }
    }
  }
}