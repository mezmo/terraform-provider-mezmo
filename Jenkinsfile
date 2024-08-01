library 'magic-butler-catalogue'
def PROJECT_NAME = 'terraform-provider-mezmo'
def DEFAULT_BRANCH = 'main'
def WORKSPACE_PATH = "/tmp/workspace/${env.BUILD_TAG.replace('%2F', '/')}"
def CURRENT_BRANCH = currentBranch()

def CREDS = [
    aws(
      credentialsId: 'aws',
      accessKeyVariable: 'AWS_ACCESS_KEY_ID',
      secretKeyVariable: 'AWS_SECRET_ACCESS_KEY'
    ),
    string(
      credentialsId: 'github-api-token',
      variable: 'GITHUB_TOKEN'
    )
]
def NPMRC = [
    configFile(fileId: 'npmrc', variable: 'NPM_CONFIG_USERCONFIG')
]

pipeline {
  agent {
    node {
      label 'ec2-fleet'
      customWorkspace(WORKSPACE_PATH)
    }
  }

  parameters {
    string(name: 'SANITY_BUILD', defaultValue: '', description: 'Is this a scheduled sanity build that skips releasing?')
  }
  triggers {
    parameterizedCron(
      // Cron hours are in GMT, so this is roughly 12-3am EST, depending on DST
      env.BRANCH_NAME == DEFAULT_BRANCH ? 'H H(5-6) * * * % SANITY_BUILD=true' : ''
    )
  }

  options {
    timeout time: 1, unit: 'HOURS'
    timestamps()
    ansiColor 'xterm'
    withCredentials(CREDS)
    disableConcurrentBuilds()
  }

  stages {
    stage('Validate PR Author') {
      when {
        expression { env.CHANGE_FORK }
        not {
          triggeredBy 'issueCommentCause'
        }
      }
      steps {
        error("A maintainer needs to approve this PR for CI by commenting")
      }
    }

    stage('Format') {
      steps {
        script {
          if (env.SANITY_BUILD == 'true') {
            currentBuild.description = "SANITY=${env.SANITY_BUILD}"
          }
        }
        sh 'FILES_TO_FORMAT=$(gofmt -l .) && echo -e "Files with formatting errors: $FILES_TO_FORMAT" && [ -z "$FILES_TO_FORMAT" ]'
      }
    }

    stage('Lint'){
      tools {
        nodejs 'NodeJS 20'
      }
      environment {
        GIT_BRANCH = "${CURRENT_BRANCH}"
        // This is not populated on PR builds and is needed for the release dry runs
        BRANCH_NAME = "${CURRENT_BRANCH}"
        CHANGE_ID = ""
      }
      steps {
        script {
          configFileProvider(NPMRC) {
            sh 'npm ci --ignore-scripts'
            sh 'npm run commitlint'
            // goreleaser handles releases with `make build`
          }
        }
      }
    }

    stage('Test') {
      parallel {
        stage('Documentation') {
          steps {
            sh 'make test-docs'
          }
        }
        stage('Unit') {
          steps {
            sh 'command -v go-junit-report || go install github.com/jstemmer/go-junit-report/v2@latest'
            sh 'make test-unit 2>&1 | go-junit-report -iocopy -set-exit-code -out unit-results.xml'
          }
          post {
            always {
              junit testResults: 'unit-results.xml'
            }
          }
        }
        stage('Integration') {
          steps {
            sh 'make test-acceptance'
          }
          post {
            always {
              junit testResults: 'results.xml', allowEmptyResults: true
            }
          }
        }
        stage('Example Validation') {
          steps {
            sh 'make -j8 -k -O examples'
          }
        }
      }
    }
  }

  post {
    always {
      script {
        jiraSendBuildInfo()

        if (env.SANITY_BUILD == 'true') {
          notifySlack(
            currentBuild.currentResult,
            [channel: '#pipeline-bots'],
            "`${PROJECT_NAME}` sanity build took ${currentBuild.durationString.replaceFirst(' and counting', '')}."
          )
        }
      }
    }
  }
}
