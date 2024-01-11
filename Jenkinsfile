
config {
	concurrentBuilds = false
	daysToKeep = 21
	cronTrigger = '@weekend'
}

node() {
	catchError {
		dir('work') {
			checkout scmGit(
				branches: [[name: '*/main']],
				extensions: [cleanAfterCheckout(deleteUntrackedNestedRepositories: true)],
				userRemoteConfigs: [[
					credentialsId: '358853c8-44a1-4a63-81c2-c89007ab2863',
					url: 'https://github.com/topicuskeyhub/sdk-go.git'
				], [
					credentialsId: '358853c8-44a1-4a63-81c2-c89007ab2863',
					url: 'https://github.com/topicuskeyhub/terraform-provider-keyhub-generator.git'
				]])
		}
	}
}
