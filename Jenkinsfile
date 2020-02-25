@Library('devops-tools')
import org.netscale.jenkins.slack.SlackNotifier

pipeline {
  agent {
    label "jenkins-go"
  }
  environment {
    SLACK_CHANNEL = 'dkv-bots'
    SLACK_DOMAIN  = 'jaraxa'
    SLACK_CREDENTIALS = 'slack-netcomp-devops'
    CHANGE_LIST = true
    TEST_SUMMARY = true    
    ORG = 'netscale-technologies'
    APP_NAME = 'gorush'
    CHARTMUSEUM_CREDS = credentials('jenkins-x-chartmuseum')
    DOCKER_REGISTRY_ORG = 'netscale-technologies'
    PROMOTE_ENV_NAME = 'environment-dkv-preprod'
    CI_BRANCH_DEV = 'develop'
    CI_BRANCH_UAT = 'staging'
    REMOTE_ENV_NAME = 'jx-dkv-remote'

  }
  stages {
    stage('Build Preview for develop') {
      when {
        branch 'develop'
      }
      environment {
        DEVELOP_VERSION = "0.0.0-SNAPSHOT-$BRANCH_NAME-$BUILD_NUMBER"
        PREVIEW_VERSION = "$DEVELOP_VERSION"
        DEVELOP_NAMESPACE = "jx-dkv-develop"
        HELM_RELEASE = "$DEVELOP_NAMESPACE".toLowerCase()
      }
      steps {
        container('go') {     
          dir('/home/jenkins/agent/src/github.com/netscale-technologies/gorush') {
            checkout scm: [$class: 'GitSCM', branches: [[name: '*/develop']], userRemoteConfigs: [[credentialsId: 'jx-pipeline-git-github-github.com', url: 'https://github.com/netscale-technologies/gorush']]]
            sh 'make get'
            sh 'make build_linux_amd64'
            sh "export VERSION=$DEVELOP_VERSION && skaffold build -f skaffold.yaml"
            sh "jx step post build --image $DOCKER_REGISTRY/$ORG/$APP_NAME:$DEVELOP_VERSION"
          }
          dir('/home/jenkins/agent/src/github.com/netscale-technologies/gorush/charts/develop') {
            sh "make preview"
            sh "jx preview --app $APP_NAME --namespace $DEVELOP_NAMESPACE --name $PROMOTE_ENV_NAME --alias $APP_NAME --label $APP_NAME --release $APP_NAME --no-comment --no-poll --no-wait --dir  ../.."
          }          
        }
      }
    }
    stage('Build Release for Staging/UAT') {
      when {
        branch 'staging'
      }
      steps {
        container('go') {
          dir('/home/jenkins/agent/src/github.com/netscale-technologies/gorush') {
            checkout scm: [$class: 'GitSCM', branches: [[name: '*/staging']], userRemoteConfigs: [[credentialsId: 'jx-pipeline-git-github-github.com', url: 'https://github.com/netscale-technologies/gorush']]]
            // ensure we're not on a detached head
            sh "git checkout $CI_BRANCH_UAT"
            sh "git config --global credential.helper store"
            sh "jx step git credentials"

            // so we can retrieve the version in later steps
            sh "echo \$(jx-release-version) > VERSION"
            sh "jx step tag --version \$(cat VERSION) --charts-dir ./charts/gorush/"
            sh 'make get'
            sh 'make build_linux_amd64'
            sh "export VERSION=`cat VERSION` && skaffold build -f skaffold.yaml"
            sh "jx step post build --image $DOCKER_REGISTRY/$ORG/$APP_NAME:\$(cat VERSION)"
          }
        }
      }
    }
    stage('Promote to staging/UAT environment') {
      when {
        branch 'staging'
      }
      environment {
        STAGING_NAMESPACE = 'jx-dkv-staging'
      }         
      steps {
        container('go') {
          dir('/home/jenkins/agent/src/github.com/netscale-technologies/gorush/charts/gorush') {
            sh "jx step changelog --version v\$(cat ../../VERSION)"

            // release the helm chart
            sh "jx step helm release"

            // promote through promotion Environment
            sh "jx promote -b --timeout 1h --version \$(cat ../../VERSION) --env $REMOTE_ENV_NAME"

          // delete unnecessary staging namespace
            sh "jx delete namespace $STAGING_NAMESPACE --batch-mode -y"              
          }
        }
      }
    }
  }
  post {
        always {
          script {
            new SlackNotifier().notifyResultFull()
          }
          cleanWs()
        }
  }
}    



 

