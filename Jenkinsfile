pipeline {
    agent any
    
    tools {
      go 'GO.1.13.1'
    }

    environment {
        GO111MODULE="on"
        GOPATH="${WORKSPACE}/go"
        OLDATA="${GOPATH}/data"
        PATH="${GOPATH}/bin:${PATH}"
    }
    stages{       
        stage ('download apt dependency'){
            steps{
                sh 'apt-get update -y && apt-get install -y build-essential libleveldb-dev libsnappy-dev'
            }
        }
        
    stages {
        
        stage('cloning protocol repo'){
            steps {
                    sh 'git clone https://github.com/Oneledger/protocol.git'
                }
            }
        }
        
        stage('performing tests'){
            steps {
                    sh 'pwd; ls -lrt; ./test_script.sh'
                }
            }
        }
