node {
    checkout scm

    // For PRs Jenkins will give the source branch name
    if (env.CHANGE_BRANCH) {
        isPR = true
        // env.BRANCH_NAME=PR-1
        version = env.BRANCH_NAME
    } else {
        isPR = false
        // env.BRANCH_NAME=k8s-1.10
        // env.BRANCH_NAME=master
        if (env.BRANCH_NAME.startsWith('k8s-')) {
            version = env.BRANCH_NAME.substring('k8s-'.size())
        } else {
            version = env.BRANCH_NAME
        }
    }

    def namespace = 'deepomatic'
    def repository = 'shared-gpu-gcp-k8s-device-plugin'

    // mount repository in GOPATH subdir
    def src = '/go/src/github.com/Deepomatic/container-engine-accelerators'
    docker.image('golang:1.9').inside(" -v ${env.WORKSPACE}:${src}") {
        stage("Test") {
            sh "cd ${src} && make presubmit test"
        }
    }

    stage("Build Docker image") {
        def image = docker.build("${namespace}/${repository}:${version}-${env.BUILD_ID}")
        if (! isPR) {
            echo "Mainline branch, pushing to repository"
            image.push()
            if (env.BRANCH_NAME == 'master') {
                image.push('latest')
            }
        }
    }
}
