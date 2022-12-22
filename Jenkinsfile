//#!groovy
//@Library(['github.com/cloudogu/ces-build-lib@1.55.0', 'github.com/cloudogu/dogu-build-lib@v1.6.0', 'github.com/cloudogu/zalenium-build-lib@v2.1.1'])
//import com.cloudogu.ces.cesbuildlib.*
//import com.cloudogu.ces.dogubuildlib.*
//import com.cloudogu.ces.zaleniumbuildlib.*
//
//node('vagrant') {
//    Git git = new Git(this, "cesmarvin")
//    git.committerName = 'cesmarvin'
//    git.committerEmail = 'cesmarvin@cloudogu.com'
//    GitFlow gitflow = new GitFlow(this, git)
//    GitHub github = new GitHub(this, git)
//    Changelog changelog = new Changelog(this)
//    Docker docker = new Docker(this)
//    Gpg gpg = new Gpg(this, docker)
//    Markdown markdown = new Markdown(this)
//
//    project = 'github.com/cloudogu/cesappd'
//    projectName = 'cesappd'
//    branch = "${env.BRANCH_NAME}"
//    githubCredentialsId = 'sonarqube-gh'
//
//    timestamps {
//        stage('Checkout') {
//            checkout scm
//        }
//
//        stage('Check Markdown Links') {
//            markdown.check()
//        }
//
//        stage('Shell-Check') {
//            shellCheck("./deb/DEBIAN/postinst ./deb/DEBIAN/postrm ./deb/DEBIAN/prerm ./deb/usr/local/bin/ssl_cesappd.sh ./deb/usr/local/bin/ssl_cesappd_generate.sh")
//        }
//
//        stage('Shell tests') {
//            executeShellTests()
//        }
//
//        withBuildDependencies{
//            withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: "${githubCredentialsId}", usernameVariable: 'USERNAME', passwordVariable: 'REVIEWDOG_GITHUB_API_TOKEN']]) {
//                sh 'git config --global url."https://$USERNAME:$REVIEWDOG_GITHUB_API_TOKEN@github.com".insteadOf "https://github.com"'
//
//                stage('Build') {
//                   make 'clean package checksum'
//                   archiveArtifacts 'target/**/*.deb'
//                   archiveArtifacts 'target/**/*.sha256sum'
//                }
//
//                stage('Unit Test') {
//                   try {
//                       make 'unit-test'
//                   } finally {
//                       junit allowEmptyResults: true, testResults: 'target/unit-tests/*-tests.xml'
//                   }
//                }
//
//                stage('Static Analysis') {
//                   def commitSha = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()
//
//                   withEnv(["CI_PULL_REQUEST=${env.CHANGE_ID}", "CI_COMMIT=${commitSha}", "CI_REPO_OWNER=cloudogu", "CI_REPO_NAME=cesappd"]) {
//                       make 'static-analysis-ci'
//                   }
//                }
//            }
//        }
//
//        stage('Sign'){
//            gpg.createSignature()
//            archiveArtifacts 'target/**/*.sha256sum.asc'
//        }
//
//        stage('SonarQube') {
//            def scannerHome = tool name: 'sonar-scanner', type: 'hudson.plugins.sonar.SonarRunnerInstallation'
//            withSonarQubeEnv {
//                sh "git config 'remote.origin.fetch' '+refs/heads/*:refs/remotes/origin/*'"
//                gitWithCredentials("fetch --all")
//
//                if (branch == "main") {
//                    echo "This branch has been detected as the main branch."
//                    sh "${scannerHome}/bin/sonar-scanner"
//                } else if (branch == "develop") {
//                    echo "This branch has been detected as the develop branch."
//                    sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME} -Dsonar.branch.target=main"
//                } else if (env.CHANGE_TARGET) {
//                    echo "This branch has been detected as a pull request."
//                    sh "${scannerHome}/bin/sonar-scanner -Dsonar.pullrequest.key=${env.CHANGE_ID} -Dsonar.pullrequest.branch=${env.CHANGE_BRANCH} -Dsonar.pullrequest.base=develop"
//                } else if (branch.startsWith("feature/")) {
//                    echo "This branch has been detected as a feature branch."
//                    sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME} -Dsonar.branch.target=develop"
//                } else if (branch.startsWith("bugfix/")) {
//                    echo "This branch has been detected as a bugfix branch."
//                    sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME} -Dsonar.branch.target=develop"
//                }  else {
//                    echo "This branch has been detected as a miscellaneous branch."
//                    sh "${scannerHome}/bin/sonar-scanner -Dsonar.projectKey=${projectName} -Dsonar.projectName=${projectName} -Dsonar.branch.name=${env.BRANCH_NAME} -Dsonar.branch.target=develop"
//                }
//            }
//            timeout(time: 2, unit: 'MINUTES') { // Needed when there is no webhook for example
//                def qGate = waitForQualityGate()
//                if (qGate.status != 'OK') {
//                    unstable("Pipeline unstable due to SonarQube quality gate failure")
//                }
//            }
//        }
//
//        if (gitflow.isReleaseBranch()) {
//            String releaseVersion = git.getSimpleBranchName();
//
//            stage('Finish Release') {
//                gitflow.finishRelease(releaseVersion, "main")
//            }
//
//            withBuildDependencies {
//                stage('Build after Release') {
//                    git.checkout(releaseVersion)
//                    make 'clean package checksum'
//                }
//
//                stage('Push to apt') {
//                    withAptlyCredentials{
//                        make 'deploy'
//                    }
//                }
//            }
//
//            stage('Sign after Release'){
//                gpg.createSignature()
//            }
//
//            stage('Add Github-Release') {
//                releaseId=github.createReleaseWithChangelog(releaseVersion, changelog, "main")
//                github.addReleaseAsset("${releaseId}", "target/cesappd.sha256sum")
//                github.addReleaseAsset("${releaseId}", "target/cesappd.sha256sum.asc")
//            }
//        }
//    }
//}
//
//void make(goal) {
//    sh "cd /go/src/${project} && make ${goal}"
//}
//
//void gitWithCredentials(String command) {
//    withCredentials([usernamePassword(credentialsId: 'cesmarvin', usernameVariable: 'GIT_AUTH_USR', passwordVariable: 'GIT_AUTH_PSW')]) {
//        sh(
//                script: "git -c credential.helper=\"!f() { echo username='\$GIT_AUTH_USR'; echo password='\$GIT_AUTH_PSW'; }; f\" " + command,
//                returnStdout: true
//        )
//    }
//}
//
//def executeShellTests() {
//    def bats_base_image = "bats/bats"
//    def bats_custom_image = "cloudogu/bats"
//    def bats_tag = "1.2.1"
//
//    def batsImage = docker.build("${bats_custom_image}:${bats_tag}", "--build-arg=BATS_BASE_IMAGE=${bats_base_image} --build-arg=BATS_TAG=${bats_tag} ./build/make/bats")
//    try {
//        sh "mkdir -p target"
//
//        batsContainer = batsImage.inside("--entrypoint='' -v ${WORKSPACE}:/workspace") {
//            sh "make unit-test-shell-ci"
//        }
//    } finally {
//        junit allowEmptyResults: true, testResults: 'target/shell_test_reports/*.xml'
//    }
//}
//
//void withAptlyCredentials(Closure closure){
//    withCredentials([usernamePassword(credentialsId: 'websites_apt-api.cloudogu.com_aptly-admin', usernameVariable: 'APT_API_USERNAME', passwordVariable: 'APT_API_PASSWORD')]) {
//        withCredentials([string(credentialsId: 'misc_signphrase_apt-api.cloudogu.com', variable: 'APT_API_SIGNPHRASE')]) {
//            closure.call()
//        }
//    }
//}
//
//void withBuildDependencies(Closure closure){
//       def etcdImage = docker.image('quay.io/coreos/etcd:v3.2.5')
//        def etcdContainerName = "${JOB_BASE_NAME}-${BUILD_NUMBER}".replaceAll("\\/|%2[fF]", "-")
//        withDockerNetwork { buildnetwork ->
//            etcdImage.withRun("--network ${buildnetwork} --name ${etcdContainerName}", 'etcd --listen-client-urls http://0.0.0.0:4001 --advertise-client-urls http://0.0.0.0:4001')
//           {
//               sh "mkdir -p jstemmer"
//                new Docker(this)
//                    .image('golang:1.18.6')
//                    .mountJenkinsUser()
//                    .inside("--network ${buildnetwork} -e ETCD=${etcdContainerName} -v ${WORKSPACE}:/go/src/${project} -v ${WORKSPACE}/jstemmer:/go/src/github.com/jstemmer -w /go/src/${project} -v ${WORKSPACE}/resources/compileHeaders/usr/include/btrfs:/usr/include/btrfs")
//                {
//                    closure.call()
//                }
//            }
//        }
//}
