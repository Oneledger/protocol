pipeline {
    agent {
   //kubernetes jenkins agent
      kubernetes {
        label 'slave'
        yaml """
  apiVersion: v1
  kind: Pod
  metadata:
    label:
      jenkins: slave
  spec:
    containers:
    - name: oneledger
      image: oneledgertech/olprotocol:latest
      command:
      tty: true
  """
      }
    }
    
    //testing stages
    stages {
      stage('build binary') {
        steps {
          container('oneledger') {
            sh 'make install_c'
          }
        } 
    }
    
        stage('unit test') {
            steps {
              container('oneledger') {
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
     }


        stage('coverage test') {
            steps {
              container('oneledger') {
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
     }

        stage('ons test') {
            steps {
              container('oneledger') {
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
     }

        stage('reward test') {
            steps {
              container('oneledger') {
                script {
                    try {
                        sh 'make rewardtest'
                    } catch (e) {
                        unstable('reward test stage failed!')
                        sh 'exit 0'
                    }
                }
            }
        }
     }
  
        stage('governance test') {
            steps {
              container('oneledger') {
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
     }

        stage('staking test') {
            steps {
              container('oneledger') {
                script {
                    try {
                        sh 'make stakingtest'
                    } catch (e) {
                        unstable('staking test stage failed!')
                        sh 'exit 0'
                    }
                }
            }
        }
     }

        stage('all test') {
            steps {
              container('oneledger') {
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
     }

        stage('rpcAuth test') {
            steps {
              container('oneledger') {
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
     }

        stage('Delegation test') {
            steps {
              container('oneledger') {
                script {
                    try {
                        sh 'make delegationtest'
                    } catch (e) {
                        unstable('delegation test stage failed!')
                        sh 'exit 0'
                    }
                }
            }
        }
     }

      stage('Result') {
        steps {
          container('oneledger') {
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
  } 

