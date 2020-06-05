pipeline {
    agent any

    environment {
        GO111MODULE="on"
        GOPATH="${WORKSPACE}/go"
        OLDATA="${GOPATH}/data"
        PATH="${GOPATH}/bin:${PATH}"
    }
    stages{       
       
        stage('cloning protocol repo'){
            steps {
                    sh 'git clone https://github.com/Oneledger/protocol.git'
                }
            }
        
        stage('performing tests'){
            steps {
                    sh 'pwd; ls -lrt; ./test_script.sh'
                }
            }
        }
}

