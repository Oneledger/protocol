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
        stage ('install'){
            steps{
                sh 'make update install'
            }
        }
        stage('unit test'){
            steps{
                sh 'make utest'
            }
        }
        stage ('full test'){
            steps{
                sh 'make fulltest'
            }
        }
        // stage('install_c'){
        //     steps{
        //         sh 'make install_c'
        //     }
        // }
    }
    post {
        always {
            archiveArtifacts artifacts: "go/bin/*", fingerprint: true, onlyIfSuccessful: true
        }
    }
}