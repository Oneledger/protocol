pipeline {
    agent any

    environment {
        GO111MODULE="on"
        GOPATH="${WORKSPACE}/go"
        OLDATA="${GOPATH}/data"
        PATH="${GOPATH}/bin:${PATH}"
    }
    stages{       
       
        stage('clone protocol repo') {
            steps {
                    checkout([
                    $class: 'GitSCM', 
        	        branches: [[name: '*/develop']], 
        	        doGenerateSubmoduleConfigurations: false, 
                    submoduleCfg: [], 
                    userRemoteConfigs: [[credentialsId: '9a3855d0-e5a5-4a47-acfd-96b75f917bbc', url: 'https://github.com/Oneledger/protocol.git']]
    ])   
            }
        }
        stage ('build binary'){
            steps{
                sh 'make install_c'
            }
        }

        stage('utest') {
          steps {
              sh 'make utest' || true
        }
    }

        stage ('validator test'){
            steps{
                sh 'make applytest'
            }
        }
        stage ('ons test'){
            steps{
                sh 'make onstest'
            }
        }
        stage ('withdraw test'){
            steps{
                sh 'make withdrawtest'
            }
        }
        stage ('governance test'){
            steps{
                sh 'make govtest'
            }
        }
        stage ('all test'){
            steps{
                sh 'make alltest'
            }
        }
        stage ('rpc Authtest'){
            steps{
                sh 'make rpcAuthtest'
            }
        }
        stage ('coverage test'){
            steps{
                sh 'make coverage'
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

