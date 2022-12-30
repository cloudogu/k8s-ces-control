#!groovy
@Library(['github.com/cloudogu/dogu-build-lib@v1.6.0', 'github.com/cloudogu/ces-build-lib@1.60.0'])
import com.cloudogu.ces.cesbuildlib.*
import com.cloudogu.ces.dogubuildlib.*

// Creating necessary git objects, object cannot be named 'git' as this conflicts with the method named 'git' from the library
gitWrapper = new Git(this, "cesmarvin")
gitWrapper.committerName = 'cesmarvin'
gitWrapper.committerEmail = 'cesmarvin@cloudogu.com'
gitflow = new GitFlow(this, gitWrapper)
github = new GitHub(this, gitWrapper)
changelog = new Changelog(this)
Docker docker = new Docker(this)
goVersion = "1.19"

// Configuration of repository
repositoryOwner = "cloudogu"
repositoryName = "k8s-ces-control"
project = "github.com/${repositoryOwner}/${repositoryName}"

// Configuration of branches
productionReleaseBranch = "main"
developmentBranch = "develop"
currentBranch = "${env.BRANCH_NAME}"

node('docker') {
    timestamps {
        properties([
                // Keep only the last x builds to preserve space
                buildDiscarder(logRotator(numToKeepStr: '10')),
                // Don't run concurrent builds for a branch, because they use the same workspace directory
                disableConcurrentBuilds(),
        ])

        stage('Checkout') {
            checkout scm
            make 'clean'
        }

        stage('Lint - Dockerfile') {
            lintDockerfile()
        }

        stage("Lint - k8s Resources") {
            stageLintK8SResources()
        }

        docker
                .image("golang:${goVersion}")
                .mountJenkinsUser()
                .inside("--volume ${WORKSPACE}:/go/src/${project} -w /go/src/${project}")
                        {
                            stage('Build') {
                                make 'compile'
                            }

                            stage('Unit Tests') {
                                make 'unit-test'
                                junit allowEmptyResults: true, testResults: 'target/unit-tests/*-tests.xml'
                            }

                            stage("Review dog analysis") {
                                stageStaticAnalysisReviewDog()
                            }
                        }

        stage('SonarQube') {
            stageStaticAnalysisSonarQube()
        }

        K3d k3d = new K3d(this, "${WORKSPACE}", "${WORKSPACE}/k3d", env.PATH)
        try {
            stage('Set up k3d cluster') {
                k3d.startK3d()
            }

            stage('Setup') {
                k3d.setup("v0.10.0", [
                        dependencies: ["official/postfix"],
                        defaultDogu : ""
                ])
            }

            stage("wait for setup") {
                k3d.waitForDeploymentRollout("postfix", 300, 10)
            }

            stage('Install locally into k3d') {
                Makefile makefile = new Makefile(this)
                def localImageName = k3d.buildAndPushToLocalRegistry("cloudogu/${repositoryName}", makefile.getVersion())
                String pathToGeneratedFile = generateResources(localImageName)
                k3d.kubectl("apply -f ${pathToGeneratedFile}")
                make("clean")
            }

            stage('Install grpcurl') {
                String grpcurlVersion = "1.8.7"
                sh "wget -O grpcurl.tar.gz https://github.com/fullstorydev/grpcurl/releases/download/v${grpcurlVersion}/grpcurl_${grpcurlVersion}_linux_x86_64.tar.gz"
                sh "tar -xf grpcurl.tar.gz"
                sh "rm -rf grpcurl.tar.gz"
            }

            stage("wait for setup") {
                k3d.waitForDeploymentRollout("k8s-ces-control", 300, 10)
            }

            stage('Test Grpc') {

            }

            stageAutomaticRelease()
        } catch (Exception e) {
            k3d.collectAndArchiveLogs()
            throw e
        } finally {
            // try to remove installed grpcurl files
            sh "rm -rf grpcurlDir"

            stage('Remove k3d cluster') {
                k3d.deleteK3d()
            }
        }
    }
}

void stageLintK8SResources() {
    String kubevalImage = "cytopia/kubeval:0.13"
    docker
            .image(kubevalImage)
            .inside("-v ${WORKSPACE}/k8s:/data -t --entrypoint=")
                    {
                        sh "kubeval /data/${repositoryName}.yaml --ignore-missing-schemas"
                    }
}

String getCurrentCommit() {
    return sh(returnStdout: true, script: 'git rev-parse HEAD').trim()
}

void stageStaticAnalysisReviewDog() {
    def commitSha = getCurrentCommit()
    withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'sonarqube-gh', usernameVariable: 'USERNAME', passwordVariable: 'REVIEWDOG_GITHUB_API_TOKEN']]) {
        withEnv(["CI_PULL_REQUEST=${env.CHANGE_ID}", "CI_COMMIT=${commitSha}", "CI_REPO_OWNER=${repositoryOwner}", "CI_REPO_NAME=${repositoryName}"]) {
            make 'static-analysis'
        }
    }
}

void stageStaticAnalysisSonarQube() {
    def scannerHome = tool name: 'sonar-scanner', type: 'hudson.plugins.sonar.SonarRunnerInstallation'
    withSonarQubeEnv {
        gitWrapper.fetch()

        if (currentBranch == productionReleaseBranch) {
            echo "This branch has been detected as the production branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME}"
        } else if (currentBranch == developmentBranch) {
            echo "This branch has been detected as the development branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME}"
        } else if (env.CHANGE_TARGET) {
            echo "This branch has been detected as a pull request."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.pullrequest.key=${env.CHANGE_ID} -Dsonar.pullrequest.branch=${env.CHANGE_BRANCH} -Dsonar.pullrequest.base=${developmentBranch}"
        } else if (currentBranch.startsWith("feature/")) {
            echo "This branch has been detected as a feature branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME}"
        } else {
            echo "This branch has been detected as a miscellaneous branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME} "
        }
    }
    timeout(time: 2, unit: 'MINUTES') { // Needed when there is no webhook for example
        def qGate = waitForQualityGate()
        if (qGate.status != 'OK') {
            unstable("Pipeline unstable due to SonarQube quality gate failure")
        }
    }
}

void stageAutomaticRelease() {
    if (gitflow.isReleaseBranch()) {
        String releaseVersion = gitWrapper.getSimpleBranchName()
        Makefile makefile = new Makefile(this)
        String version = makefile.getVersion()

        stage('Build & Push Image') {
            def dockerImage = docker.build("cloudogu/${repositoryName}:${version}")

            docker.withRegistry('https://registry.hub.docker.com/', 'dockerHubCredentials') {
                dockerImage.push("${version}")
            }
        }

        stage('Finish Release') {
            gitflow.finishRelease(releaseVersion, productionReleaseBranch)
        }

        stage('Add Github-Release') {
            releaseId = github.createReleaseWithChangelog(releaseVersion, changelog, productionReleaseBranch)
        }

        stage('Regenerate resources for release') {
            make 'create-temporary-release-resource'
        }

        stage('Push to Registry') {
            GString targetSetupResourceYaml = "target/make/k8s/${repositoryName}_${version}.yaml"

            DoguRegistry registry = new DoguRegistry(this)
            registry.pushK8sYaml(targetSetupResourceYaml, repositoryName, "k8s", "${version}")
        }
    }
}

String generateResources(String image = "") {
    Makefile makefile = new Makefile(this)
    String version = makefile.getVersion()
    String generatedFile = "target/make/k8s/k8s-ces-control_${version}.yaml".toString()
    docker.image('mikefarah/yq:4.22.1')
            .mountJenkinsUser()
            .inside("--volume ${WORKSPACE}:/workdir -w /workdir") {
                if (image == "") {
                    sh "IMAGE=${image} LOG_LEVEL=DEBUG make k8s-generate"
                } else {
                    sh "make k8s-generate"
                }
                archiveArtifacts "${generatedFile}"
            }

    return generatedFile
}

void make(String makeArgs) {
    sh "make ${makeArgs}"
}