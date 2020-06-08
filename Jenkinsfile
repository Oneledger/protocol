pipeline {
    agent any

    environment {
        GO111MODULE="on"
        GOPATH="${WORKSPACE}/go"
        OLDATA="${GOPATH}/data"
        PATH="${GOPATH}/bin:${PATH}"
    }
    stages{       
       
        stage('performing tests'){
            steps {
                    sh 'pwd; ls -lrt; ./test_script.sh'
                }
            }
        
      stage('Results') {
          steps {
             publishHTML([allowMissing: false,
             alwaysLinkToLastBuild: true,
             keepAll: true,
             reportDir: 
                          '*',
             reportFiles: 'cover.html',
             reportName: 'Coverage report Dashboard'
])

             } 
           }
        }
}

