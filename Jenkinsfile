pipeline {
    agent any
    options {
        skipStagesAfterUnstable()
    }
    // Ensure the desired Go version is installed for all stages,
    // using the name defined in the Global Tool Configuration
    tools {
        go '1.19'
        dockerTool "docker"
    }
    environment{
        registryCredential=credentials("lisportalprod-docker-pwd")
        AZURE_PIPLINE_STATUS    = 'start'
        GIT_SHA = sh(script: "git rev-parse HEAD", returnStdout: true).trim()
    }
    stages{
        stage('Test') {
            steps {
                println('start testing')
                script {
                    sh 'apt-get -y install redis-server'
                    sh 'apt-get -y install gcc'
                    sh 'go test ./... -v'
                }
            }
        }
        stage('Build') {
            when {
                expression {
                    ("${env.BRANCH_NAME}" == 'main' || "${env.BRANCH_NAME}" == 'staging')
                }
            }
            steps {
                println('start build')
                script{
                    if ("${env.BRANCH_NAME}" == 'staging') {
                        sh 'docker login -u 4693646 -p Yy8m7KJSm@Lm2'
                        if ("${env.JENKINS_URL}".startsWith('http://192.168.10')){
                            sh 'docker build -t 192.168.10.62:6004/vibrant/lis/coresamples-v2:staging .'
                        } else if ("${env.JENKINS_URL}".startsWith('http://192.168.60')){
                            sh 'docker build -t 192.168.60.10:6004/vibrant/lis/coresamples-v2:staging .'
                        } else {
                            error "jenkins url match error: ${env.JENKINS_URL}"
                        }

                        try {
                            sh 'docker login -u lisportalprod -p ${registryCredential} http://lisportalprod.azurecr.io'
                            sh 'docker build -t lisportalprod.azurecr.io/vibrant/lis/coresamples-v2:staging .'
                            AZURE_PIPLINE_STATUS = 'build'
                            }catch (Exception e) {
                                echo 'Exception occurred: ' + e.toString()
                            }
                    } else {
                        try {
                        sh 'docker login -u lisportalprod -p ${registryCredential} http://lisportalprod.azurecr.io'
                        sh 'docker build -t lisportalprod.azurecr.io/vibrant/lis/coresamples-v2:${GIT_SHA} .'
                        AZURE_PIPLINE_STATUS = 'build'
                        }catch (Exception e) {
                            echo 'Exception occurred: ' + e.toString()
                        }
                    }
                }
                println('build finish')
            }
        }
        stage ('Push'){
            when {
                expression {
                    ("${env.BRANCH_NAME}" == 'main' || "${env.BRANCH_NAME}" == 'staging')
                }
            }
            steps{
                println('start push')
                script{
                    if ("${env.BRANCH_NAME}" == 'staging') {
                        if ("${env.JENKINS_URL}".startsWith('http://192.168.10')){
                            sh 'docker push 192.168.10.62:6004/vibrant/lis/coresamples-v2:staging'
                        } else {
                            sh 'docker push 192.168.60.10:6004/vibrant/lis/coresamples-v2:staging'
                        }
                        //push cloud images
                        //sh 'docker login -u listestdocker -p ${registryCredential} http://listestdocker.azurecr.io'
                        try{
                            sh 'docker push lisportalprod.azurecr.io/vibrant/lis/coresamples-v2:staging'
                            AZURE_PIPLINE_STATUS = 'push'
                        }catch (Exception e) {
                            echo 'Exception occurred: ' + e.toString()
                        }
                    } else {
                        try{
                            sh 'docker push lisportalprod.azurecr.io/vibrant/lis/coresamples-v2:${GIT_SHA}'
                            AZURE_PIPLINE_STATUS = 'push'
                        }catch (Exception e) {
                            echo 'Exception occurred: ' + e.toString()
                        }
                    }
                }
                println('push finish')
            }
        }
        stage('Deploy'){
            when {
                expression {
                    ("${env.BRANCH_NAME}" == 'main' || "${env.BRANCH_NAME}" == 'staging')
                }
            }
            steps{
                println('deploy')
                script{
                    if ("${env.BRANCH_NAME}" == 'staging') {
                        if ("${env.JENKINS_URL}".startsWith('http://192.168.10')){
                            sh 'ssh lis_updater@192.168.10.212 kubectl rollout restart deployment/lis-coresamples-v2-deployment-staging'
                        } else {
                            sh 'ssh yuxuan@192.168.60.6 kubectl rollout restart deployment/lis-coresamples-v2-deployment-staging'
                        }

                        //rollout restart aks deployment
                        try{
                            withKubeConfig([credentialsId:'lisportalprod-kube-config',serverUrl:'https://lisportalprod-dns-nrpmbcaa.hcp.westus3.azmk8s.io:443']){
                                sh 'command -v kubectl || curl -LO "https://storage.googleapis.com/kubernetes-release/release/v1.20.5/bin/linux/amd64/kubectl"'
                                sh 'chmod u+x ./kubectl'
                                sh "./kubectl rollout restart deployment/lis-coresamples-v2-deployment-staging -n coresamplesv2"
                            }
                            AZURE_PIPLINE_STATUS = 'deploy'
                        }catch (Exception e) {
                            echo 'Exception occurred: ' + e.toString()
                        }
                    } else {
                        //rollout restart aks deployment
                        try{
                            withKubeConfig([credentialsId:'lisportalprod-kube-config',serverUrl:'https://lisportalprod-dns-nrpmbcaa.hcp.westus3.azmk8s.io:443']){
                                sh 'command -v kubectl || curl -LO "https://storage.googleapis.com/kubernetes-release/release/v1.20.5/bin/linux/amd64/kubectl"'
                                sh 'chmod u+x ./kubectl'
                                sh "./kubectl set image deployment/lis-coresamples-v2-deployment lis-coresamples-v2=lisportalprod.azurecr.io/vibrant/lis/coresamples-v2:${GIT_SHA} -n coresamplesv2"
                            }
                            AZURE_PIPLINE_STATUS = 'deploy'
                        }catch (Exception e) {
                            echo 'Exception occurred: ' + e.toString()
                        }
                    }

                    if ("${env.BRANCH_NAME}" == 'staging') {
                    jiraSendDeploymentInfo environmentId: 'us-staging-1', environmentName: 'us-staging-1', environmentType: 'staging', site: 'vibrantamerica.atlassian.net', serviceIds: ['lis-coresamples-v2']
                    } else {
                    jiraSendDeploymentInfo environmentId: 'us-prod-1', environmentName: 'us-prod-1', environmentType: 'production', site: 'vibrantamerica.atlassian.net', serviceIds: ['lis-coresamples-v2']
                    }
                }
                println('deploy finish')
            }
        }
    }
     post {
            success {
                script {
                    if ("${env.BRANCH_NAME}" == 'main' || "${env.BRANCH_NAME}" == 'staging') {
                    notifySlack('SUCCESS')
                    }
                }
            }
            failure {
                script {

                    if ("${env.BRANCH_NAME}" == 'main' || "${env.BRANCH_NAME}" == 'staging') {

                    notifySlack('FAILURE')
                    }

                }
            }
        }
}
def notifySlack(String buildStatus = 'STARTED') {
    // Build status of null means success.
    buildStatus = buildStatus ?: 'SUCCESS'

    def color
    def environmentType
    if ("${env.BRANCH_NAME}" == "main") {
        environmentType = 'production'
    } else if ("${env.BRANCH_NAME}" == "staging") {
        environmentType = 'staging'
    }
    env.GIT_COMMIT_MSG = sh (script: 'git log -1 --pretty=%B ${GIT_COMMIT}', returnStdout: true).trim()
    env.GIT_AUTHOR = sh (script:"git log -1 --pretty=format:'%an'", returnStdout: true).trim()
    if (buildStatus == 'STARTED') {
        color = '#D4DADF'
    } else if (buildStatus == 'SUCCESS') {
        color = '#BDFFC3'
    } else if (buildStatus == 'FAILURE') {
        color = '#FFFE89'
    } else {
        color = '#FF9FA1'
    }


    def azureBuildStatus
    if ("${AZURE_PIPLINE_STATUS}" == "deploy") {
        azureBuildStatus = 'SUCCESS'
    }else{
        azureBuildStatus = 'FAIL'
    }
    def msg = "${buildStatus}: `${env.JOB_NAME}` branch at `${GIT_COMMIT[0..7]}` by ${env.GIT_AUTHOR} on ${environmentType}.\n Azure-Deploy:${azureBuildStatus}. Last commit message: ${env.GIT_COMMIT_MSG} <${env.BUILD_URL}|#${env.BUILD_NUMBER}>"
//     def msg = "${buildStatus}: `${env.JOB_NAME}` branch at `${GIT_COMMIT[0..7]}` by ${env.GIT_AUTHOR} on ${environmentType}. Last commit message: ${env.GIT_COMMIT_MSG} <${env.BUILD_URL}|#${env.BUILD_NUMBER}>"

    slackSend(color: color, message: msg, channel: 'lis-bot')
}