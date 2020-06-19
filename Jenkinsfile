pipeline {
    agent any

    environment {
        GO111MODULE="on"
        GOPATH="${WORKSPACE}/go"
        OLDATA="${GOPATH}/data"
        PATH="${GOPATH}/bin:${PATH}"
        OLTEST="1"
    }
    stages{       
      
        stage ('build binary'){
            steps{
                sh 'make install_c'
            }
        }
        stage('unit test') {
            steps {
                script {
                    try {
                        sh 'make utest'
                    } catch (e) {
                        unstable('Unit test stage failed!')
                        sh 'exit 0'
                    }
                }
            }
        }

        stage('coverage test') {
            steps {
                script {
                    try {
                        sh 'make coverage'
                    } catch (e) {
                        unstable('coverage test stage failed!')
                        sh 'exit 0'
                    }
                }
            }
        }
        
       stage('ons test') {
            steps {
                script {
                    try {
                        sh 'make onstest'
                    } catch (e) {
                        unstable('ons test stage failed!')
                        sh 'exit 0'
                    }
                }
            }
        }
        
        stage('withdraw test') {
            steps {
                script {
                    try {
                        sh 'make withdrawtest'
                    } catch (e) {
                        unstable('withdraw test stage failed!')
                        sh 'exit 0'
                    }
                }
            }
        }
         
       stage('governance test') {
            steps {
                script {
                    try {
                        sh 'make govtest'
                    } catch (e) {
                        unstable('governance test stage failed!')
                        sh 'exit 0'
                    }
                }
            }
        }
        
        stage('all test') {
            steps {
                script {
                    try {
                        sh 'make alltest'
                    } catch (e) {
                        unstable('all test stage failed!')
                        sh 'exit 0'
                    }
                }
            }
        }
        
        stage('rpcAuth test') {
            steps {
                script {
                    try {
                        sh 'make rpcAuthtest'
                    } catch (e) {
                        unstable('rpcAuth test stage failed!')
                        sh 'exit 0'
                    }
                }
            }
        }
         
        stage('Results') {
          steps {
             publishHTML([allowMissing: false,
             alwaysLinkToLastBuild: true,
             keepAll: true,
             reportDir: 
                          '.',
             reportFiles: 'cover.html',
             reportName: 'Coverage report Dashboard'
])

             } 
           }
        }
   
}
