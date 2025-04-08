@Library('jc21') _

pipeline {
	options {
		buildDiscarder(logRotator(numToKeepStr: '5'))
		disableConcurrentBuilds()
		ansiColor('xterm')
	}
	agent {
		label 'docker'
	}
	stages {
		stage('Unit Test') {
			steps {
				sh 'docker run --rm -v "$(pwd):/workspace" -w /workspace jc21/gotools ./scripts/unit.sh'
				sh 'docker run --rm -v "$(pwd):/workspace" -w /workspace jc21/gotools chown -R "$(id -u):$(id -g)" /workspace'
			}
			post {
				always {
					dir('test/results') {
						archiveArtifacts(
							artifacts: 'html-reports/**.*',
							allowEmptyArchive: true
						)
					}
					junit allowEmptyResults: true, testResults: 'test/results/junit/*'
				}
			}
		}
		stage('Integration Test') {
			when {
				not {
					equals expected: 'UNSTABLE', actual: currentBuild.result
				}
			}
			steps {
				sh "./scripts/test.sh"
			}
		}
		stage('Publish') {
			when {
				not {
					equals expected: 'UNSTABLE', actual: currentBuild.result
				}
			}
			steps {
				incrementGithubRelease("jc21/go-docker-registry")
			}
		}
	}
	post {
		always {
			printResult(true)
		}
	}
}
