
config {
	concurrentBuilds = false
	daysToKeep = 21
	cronTrigger = '@weekend'
}

node() {
	catchError {
		stage('Clone repos') {
			dir('work/sdk-go') {
				checkout scmGit(
					branches: [[name: '*/main']],
					extensions: [cleanAfterCheckout(deleteUntrackedNestedRepositories: true)],
					userRemoteConfigs: [[
						credentialsId: '358853c8-44a1-4a63-81c2-c89007ab2863',
						url: 'https://github.com/topicuskeyhub/sdk-go.git'
					]])
			}
			dir('work/terraform-provider-keyhub-generator') {
				checkout scmGit(
					branches: [[name: '*/main']],
					extensions: [cleanAfterCheckout(deleteUntrackedNestedRepositories: true)],
					userRemoteConfigs: [[
						credentialsId: '358853c8-44a1-4a63-81c2-c89007ab2863',
						url: 'https://github.com/topicuskeyhub/terraform-provider-keyhub-generator.git'
					]])
			}
			dir('work/terraform-provider-keyhub') {
				checkout scmGit(
					branches: [[name: '*/main']],
					extensions: [cleanAfterCheckout(deleteUntrackedNestedRepositories: true)],
					userRemoteConfigs: [[
						credentialsId: '358853c8-44a1-4a63-81c2-c89007ab2863',
						url: 'https://github.com/topicuskeyhub/terraform-provider-keyhub.git'
					]])
			}
			dir('work/terraform-test-api') {
				checkout scmGit(
					branches: [[name: '*/main']],
					extensions: [cleanAfterCheckout(deleteUntrackedNestedRepositories: true)],
					userRemoteConfigs: [[
						credentialsId: '358853c8-44a1-4a63-81c2-c89007ab2863',
						url: 'https://github.com/topicuskeyhub/terraform-test-api.git'
					]])
			}
		}
	}
}
