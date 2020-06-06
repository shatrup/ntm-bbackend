node {
  def root = tool name: 'Go1.14.4', type: 'go'
  def version = ''
  def project_name = 'ntm-backend'

  ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src/ntm-backend") {
    withEnv(["GOROOT=${root}", "GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/", "PATH+GO=${root}/bin"]) {
      env.PATH="${GOPATH}/bin:$PATH:/var/lib/jenkins/go/bin"

      stage 'Checkout'
      slackSend (tokenCredentialId: 'slack-token', channel: '#builds', color: '#FFFF00', message: "STARTED: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})")
      git branch: 'master',
        credentialsId: 'jenkins-builder-gitlab',
        url: 'git@gitlab.com:bhaifi-dev/ntm-backend.git'
      script {
          version = sh (returnStdout:true, script: 'date +%d%m%Y_%H%M').trim()
          sh "echo '${version}' > version.txt"
      }

      stage 'Dependencies'
      sh 'go version'
      sh 'go get -u golang.org/x/tools/go/packages'
      sh "echo Generating mocks for unit tests..."
      sh '${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src/ntm-backend/generate-mock.sh'
      sh 'go mod vendor'

      stage 'Build'
      sh "go build"

      stage 'Test'
      sh "echo Running unit tests..."
      sh '${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src/ntm-backend/run-tests.sh'

      stage 'Package'
      script {
        def artifact_name = "${project_name}-${version}"
        def artifact_directory = "/tmp/${artifact_name}"
        def tar_artifact_path = "${artifact_directory}/${artifact_name}.tar.gz"

        sh "mkdir ${artifact_directory}"
        sh "tar -cvzf ${tar_artifact_path} ${project_name}"
      }

      stage 'Release artifact'
      sh "echo releasing the artifact to s3 bucket..."
      sh "ANSIBLE_CONFIG='/opt/playbooks/ansible.cfg' ansible-playbook /opt/playbooks/playbook-release-artifact-s3.yml -e 'type=webservers version=${version} project_name=${project_name}'"

      stage 'Deploy to Test environment'
      sh "ANSIBLE_CONFIG='/opt/playbooks/ansible.cfg' ansible-playbook /opt/playbooks/playbook-deploy-go.yml -e 'deploy_hosts=test type=webservers version=${version} project_name=${project_name}'"
    }
  }
}