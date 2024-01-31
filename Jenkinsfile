
config {
	concurrentBuilds = false
	daysToKeep = 21
	cronTrigger = 'H 7 * * *'
}

node() {
	catchError {
		stage('Clone repos') {
			git.checkout { }

			dir('docker/work/sdk-go') {
				checkout scmGit(
					branches: [[name: '*/main']],
					extensions: [cleanAfterCheckout(deleteUntrackedNestedRepositories: true)],
					userRemoteConfigs: [[
						credentialsId: '358853c8-44a1-4a63-81c2-c89007ab2863',
						url: 'https://github.com/topicuskeyhub/sdk-go.git'
					]])
			}
			dir('docker/work/terraform-provider-keyhub-generator') {
				checkout scmGit(
					branches: [[name: '*/main']],
					extensions: [cleanAfterCheckout(deleteUntrackedNestedRepositories: true)],
					userRemoteConfigs: [[
						credentialsId: '358853c8-44a1-4a63-81c2-c89007ab2863',
						url: 'https://github.com/topicuskeyhub/terraform-provider-keyhub-generator.git'
					]])
			}
			dir('docker/work/terraform-provider-keyhub') {
				checkout scmGit(
					branches: [[name: '*/main']],
					extensions: [cleanAfterCheckout(deleteUntrackedNestedRepositories: true)],
					userRemoteConfigs: [[
						credentialsId: '358853c8-44a1-4a63-81c2-c89007ab2863',
						url: 'https://github.com/topicuskeyhub/terraform-provider-keyhub.git'
					]])
			}
			dir('docker/work/terraform-test-api') {
				checkout scmGit(
					branches: [[name: '*/main']],
					extensions: [cleanAfterCheckout(deleteUntrackedNestedRepositories: true)],
					userRemoteConfigs: [[
						credentialsId: '358853c8-44a1-4a63-81c2-c89007ab2863',
						url: 'https://github.com/topicuskeyhub/terraform-test-api.git'
					]])
			}
			dir('docker/work/terraform-plugin-framework') {
				checkout scmGit(
					branches: [[name: '*/debug']],
					extensions: [cleanAfterCheckout(deleteUntrackedNestedRepositories: true)],
					userRemoteConfigs: [[
						credentialsId: '358853c8-44a1-4a63-81c2-c89007ab2863',
						url: 'https://github.com/topicuskeyhub/terraform-plugin-framework.git'
					]])
			}
		}

		stage('Build docker container') {
			def img = dockerfile.build {
				root = 'docker'
				name = 'keyhub/terraform-tester'
			}
			dockerfile.publish {
				image = img
				tags = [ "latest" ]
			}
		}
	}
}
