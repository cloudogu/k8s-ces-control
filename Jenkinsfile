#!groovy
@Library('github.com/cloudogu/ces-build-lib@1.65.0')
import com.cloudogu.ces.cesbuildlib.*

// Creating necessary git objects, object cannot be named 'git' as this conflicts with the method named 'git' from the library
gitWrapper = new Git(this, "cesmarvin")
gitWrapper.committerName = 'cesmarvin'
gitWrapper.committerEmail = 'cesmarvin@cloudogu.com'
gitflow = new GitFlow(this, gitWrapper)
github = new GitHub(this, gitWrapper)
changelog = new Changelog(this)
Docker docker = new Docker(this)
goVersion = "1.21.4"

// Configuration of repository
repositoryOwner = "cloudogu"
repositoryName = "k8s-ces-control"
project = "github.com/${repositoryOwner}/${repositoryName}"
registry = "registry.cloudogu.com"
registry_namespace = "k8s"

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
                k3d.setup("v0.16.0", [
                        dependencies: ["official/postfix", "official/ldap", "k8s/nginx-static", "k8s/nginx-ingress"],
                        defaultDogu : ""
                ])
            }

            stage("Wait for Setup") {
                k3d.waitForDeploymentRollout("postfix", 300, 10)
            }

            stage('Install k8s-ces-control') {
                Makefile makefile = new Makefile(this)
                def localImageName = k3d.buildAndPushToLocalRegistry("cloudogu/${repositoryName}", makefile.getVersion())
                String pathToGeneratedFile = generateResources("IMAGE_DEV=${localImageName} STAGE=development LOG_LEVEL=DEBUG make k8s-generate")
                k3d.kubectl("apply -f ${pathToGeneratedFile}")
                make("clean")
            }

            stage("Wait for k8s-ces-control") {
                k3d.waitForDeploymentRollout("k8s-ces-control", 300, 10)
            }

            stage('Test Grpc') {
                testK8sCesControl(k3d)
            }

            stageAutomaticRelease()
        } catch (Exception e) {
            k3d.collectAndArchiveLogs()
            throw e
        } finally {
            stage('Remove k3d cluster') {
                k3d.deleteK3d()
            }
        }
    }
}

private void testK8sCesControl(K3d k3d) {
    sh "KUBECONFIG=${WORKSPACE}/k3d/.k3d/.kube/config make integration-test-bash"
    junit allowEmptyResults: true, testResults: 'target/bash-integration-test/*.xml'
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
            generateResources("make k8s-create-temporary-resource")
        }

        stage('Push to Registry') {
            GString targetSetupResourceYaml = "target/make/k8s/${repositoryName}_${version}.yaml"

            DoguRegistry registry = new DoguRegistry(this)
            registry.pushK8sYaml(targetSetupResourceYaml, repositoryName, "k8s", "${version}")
        }

        stage('Push Helm chart to Harbor') {
            new Docker(this)
                .image("golang:${goVersion}")
                .mountJenkinsUser()
                .inside("--volume ${WORKSPACE}:/go/src/${project} -w /go/src/${project}")
                        {
                            make 'helm-package-release'

                            withCredentials([usernamePassword(credentialsId: 'harborhelmchartpush', usernameVariable: 'HARBOR_USERNAME', passwordVariable: 'HARBOR_PASSWORD')]) {
                                sh ".bin/helm registry login ${registry} --username '${HARBOR_USERNAME}' --password '${HARBOR_PASSWORD}'"
                                sh ".bin/helm push target/make/k8s/helm/${repositoryName}-${version}.tgz oci://${registry}/${registry_namespace}/"
                            }
                        }
        }
    }
}

String generateResources(String makefileCommand = "") {
    Makefile makefile = new Makefile(this)
    String version = makefile.getVersion()
    String generatedFile = "target/make/k8s/k8s-ces-control_${version}.yaml".toString()
    new Docker(this).image("golang:${goVersion}")
            .mountJenkinsUser()
            .inside("--volume ${WORKSPACE}:/workdir -w /workdir") {
                sh "${makefileCommand}"
                archiveArtifacts "${generatedFile}"
            }

    return generatedFile
}

void make(String makeArgs) {
    sh "make ${makeArgs}"
}