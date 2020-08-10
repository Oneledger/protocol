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

        stage('clean up') {
            steps {
                sh 'ls'
            }
        }

    post {
        cleanup {
            echo 'One way or another, I have finished'
            deleteDir() /* clean up our workspace */
               }
           }
            
        }
   
}
