library 'magic-butler-catalogue'
def PROJECT_NAME = 'terraform-provider-mezmo'
def DEFAULT_BRANCH = 'main'
def CURRENT_BRANCH = [env.CHANGE_BRANCH, env.BRANCH_NAME]?.find{branch -> branch != null}

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

def slugify(str) {
  def s = str.toLowerCase()
  s = s.replaceAll(/[^a-z0-9\s-\/]/, "").replaceAll(/\s+/, " ").trim()
  s = s.replaceAll(/[\/\s]/, '-').replaceAll(/-{2,}/, '-')
  s
}

pipeline {
  agent {
    node {
      label 'ec2-fleet'
      customWorkspace("/tmp/workspace/${env.BUILD_TAG}")
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

  environment {
    FEATURE_TAG = slugify("${CURRENT_BRANCH}-${BUILD_NUMBER}")
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

  stages {
    stage('Format') {
      tools {
        nodejs 'NodeJS 16'
      }
      agent {
        node {
          label 'ec2-fleet'
          customWorkspace "${PROJECT_NAME}-${BUILD_NUMBER}-integration_tests"
        }
      }
      steps {
        script {
          currentBuild.description = "SANITY=${env.SANITY_BUILD}"
        }

        sh 'FILES_TO_FORMAT=$(gofmt -l .) && echo -e "Files with formatting errors: $FILES_TO_FORMAT" && [ -z "$FILES_TO_FORMAT" ]'
      }
    }

    stage('Test') {
      when {
        not {
          changelog '\\[skip ci\\]'
        }
      }

      parallel {
        stage('Integration Tests') {
          tools {
            nodejs 'NodeJS 16'
          }
          agent {
            node {
              label 'ec2-fleet'
              customWorkspace "${PROJECT_NAME}-${BUILD_NUMBER}-integration_tests"
            }
          }
          steps {
            sh 'make test'
          }
        }
        stage('Example Validation') {
          agent {
            node {
              label 'ec2-fleet'
              customWorkspace "${PROJECT_NAME}-${BUILD_NUMBER}-example_validation"
            }
          }
          steps {
            sh 'make -j8 -k -O examples'
          }
        }
      }
    }
  }
}
