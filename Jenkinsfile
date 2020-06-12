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

       
        stage('validator test') {
            steps {
                 sh 'make applytest'
            }
        }
        
        stage('ons test') {
            steps {
                 sh 'make onstest'
            }
        }
        
        stage('withdraw test') {
            steps {
                 sh 'make withdrawtest'
            }
        }
         
        stage('governance test') {
            steps {
                 sh 'make govtest'
            }
        }
        
        stage('all test') {
            steps {
                 sh 'make alltest'
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
