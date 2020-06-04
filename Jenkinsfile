// Pipeline script to invoke CI job

pipeline {
    
    environment {
        ANSIBLE_DIR = 'ansible-scripts'
    }

    stages {
        
        stage('clone infrastructure repo') {
            steps {
                 dir(env.ANSIBLE_DIR) {
                    checkout([
                    $class: 'GitSCM', 
        	        branches: [[name: '*/devops']], 
        	        doGenerateSubmoduleConfigurations: false, 
                    extensions: [[$class: 'CleanCheckout']], 
                    submoduleCfg: [], 
                    userRemoteConfigs: [[credentialsId: '9a3855d0-e5a5-4a47-acfd-96b75f917bbc', url: 'https://github.com/Oneledger/infrastructure.git']]
    ])                
            }
        }
    }
       stage('update binary') {
            steps {
                dir(env.ANSIBLE_DIR) {
                    sh 'ansible-playbook -i hosts_testnet.yml devnet2_deploy_script.yml --tags "update"'
                }
            }
        }
    
     }

  }
