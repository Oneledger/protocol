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
                    extensions: [[$class: 'CleanCheckout']], 
                    submoduleCfg: [], 
                    userRemoteConfigs: [[credentialsId: '9a3855d0-e5a5-4a47-acfd-96b75f917bbc', url: 'https://github.com/Oneledger/protocol.git']]
    ])   
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

