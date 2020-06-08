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
             publishHTML([allowMissing: false,
             alwaysLinkToLastBuild: true,
             keepAll: true,
             reportDir: 
             '/var/lib/jenkins/workspace/pipeline-job_jenkins-test/*',
             reportFiles: 'cover.html',
             reportName: 'Docs Loadtest Dashboard'
])

             } 
        }
}

