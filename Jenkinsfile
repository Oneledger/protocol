pipeline {
    environment {
        DEPLOY_DIR = 'infrastructure/ansible-scripts'
    }
    stages {
        stage('Deploy') {
            steps {
                dir(env.DEPLOY_DIR) {
                    sh 'ansible-playbook -i hosts_devnet3.yml devnet_deploy_script.yml --tags "update"'
                }
            }
        }

    }
}
