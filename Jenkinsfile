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
                sh 'make install_c; pwd'
            }
        }

        stage('unit testing') {
            steps {
                script {
                    try {
                        sh 'make utest'
                    } catch (e) {
                        unstable('Unit testing stage failed!')
                        sh 'exit 0'
                    }
                }
            }
        }


        stage('validator test') {
            steps {
                script {
                    try {
                        sh 'make applytest'
                    } catch (e) {
                        unstable('validator testing stage failed!')
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
                        unstable('ons testing stage failed!')
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
                        unstable('withdraw testing stage failed!')
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
                        unstable('governance testing stage failed!')
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
                        unstable('all testing stage failed!')
                        sh 'exit 0'
                    }
                }
            }
        }
         stage('rpcAuthtest') {
            steps {
                script {
                    try {
                        sh 'make rpcAuthtest'
                    } catch (e) {
                        unstable('rpcAuthtesting stage failed!')
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
                        unstable('coverage testing stage failed!')
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
