#!groovy
@Library('github.com/cloudogu/ces-build-lib@4.0.1')
import com.cloudogu.ces.cesbuildlib.*

// Creating necessary git objects, object cannot be named 'git' as this conflicts with the method named 'git' from the library
gitWrapper = new Git(this, "cesmarvin")
gitWrapper.committerName = 'cesmarvin'
gitWrapper.committerEmail = 'cesmarvin@cloudogu.com'
gitflow = new GitFlow(this, gitWrapper)
github = new GitHub(this, gitWrapper)
changelog = new Changelog(this)
Docker docker = new Docker(this)
gpg = new Gpg(this, docker)
Makefile makefile = new Makefile(this)
goVersion = "1.25.1"

// Configuration of repository
repositoryOwner = "cloudogu"
repositoryName = "k8s-ces-control"
project = "github.com/${repositoryOwner}/${repositoryName}"
registry = "registry.cloudogu.com"
registry_namespace = "k8s"
helmTargetDir = "target/k8s"
helmChartDir = "${helmTargetDir}/helm"

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
                parameters([
                    choice(name: 'TrivySeverityLevels', choices: [TrivySeverityLevel.CRITICAL, TrivySeverityLevel.HIGH_AND_ABOVE, TrivySeverityLevel.MEDIUM_AND_ABOVE, TrivySeverityLevel.ALL], description: 'The levels to scan with trivy'),
                    choice(name: 'TrivyStrategy', choices: [TrivyScanStrategy.UNSTABLE, TrivyScanStrategy.FAIL, TrivyScanStrategy.IGNORE], description: 'Define whether the build should be unstable, fail or whether the error should be ignored if any vulnerability was found.'),
                ])
        ])

        stage('Checkout') {
            checkout scm
            make 'clean'
        }

        stage('Lint - Dockerfile') {
            lintDockerfile()
        }

        docker
                .image("golang:${goVersion}")
                .mountJenkinsUser()
                .inside("--volume ${WORKSPACE}:/go/src/${project} -w /go/src/${project}") {
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

                    stage('Generate k8s Resources') {
                        make 'helm-generate'
                        archiveArtifacts "${helmTargetDir}/**/*"
                    }

                    stage("Lint helm") {
                        make 'helm-lint'
                    }
                }

        stage('SonarQube') {
            stageStaticAnalysisSonarQube()
        }

        def k3d = new K3d(this, "${WORKSPACE}", "${WORKSPACE}/k3d", env.PATH)
        try {
            stage('Set up k3d cluster') {
                k3d.startK3d()
            }

            stage('Setup') {
                k3d.configureComponents(["k8s-minio":    ["version": "latest", "helmRepositoryNamespace": "k8s"],
                                         "k8s-loki":     ["version": "latest", "helmRepositoryNamespace": "k8s"],
                                         "k8s-prometheus": ["version": "latest", "helmRepositoryNamespace": "k8s", "valuesYamlOverwrite": "kube-prometheus-stack:\n  nodeExporter:\n    enabled: false"],
                                         "k8s-support-archive-operator-crd": ["version": "latest", "helmRepositoryNamespace": "k8s"],
                                         "k8s-support-archive-operator": ["version": "latest", "helmRepositoryNamespace": "k8s"]
                ])
                k3d.setup('4.2.0')
            }

            stage("Wait for Setup") {
                k3d.waitForDeploymentRollout("postfix", 300, 10)
            }

            String version = makefile.getVersion()
            String localImageName = "cloudogu/${repositoryName}:${version}"
            String imageName = ""
            stage('Build & Push Image') {
                imageName = k3d.buildAndPushToLocalRegistry("cloudogu/${repositoryName}", version)
            }

            stage('Update development resources') {
                def repository = imageName.substring(0, imageName.lastIndexOf(":"))
                docker.image("golang:${goVersion}")
                        .mountJenkinsUser()
                        .inside("--volume ${WORKSPACE}:/workdir -w /workdir") {
                            sh "STAGE=development IMAGE_DEV=${repository} make helm-values-replace-image-repo template-stage"
                        }
            }

            stage('Deploy k8s-ces-control') {
                k3d.helm("install ${repositoryName} ${helmChartDir}")
            }

            stage('Trivy scan') {
                Trivy trivy = new Trivy(this)
                trivy.scanImage(localImageName, params.TrivySeverityLevels, params.TrivyStrategy)
                trivy.saveFormattedTrivyReport(TrivyScanFormat.TABLE)
                trivy.saveFormattedTrivyReport(TrivyScanFormat.JSON)
                trivy.saveFormattedTrivyReport(TrivyScanFormat.HTML)
            }

            stage("Wait for k8s-ces-control") {
                k3d.waitForDeploymentRollout("k8s-ces-control", 300, 10)
            }

            stage('Test Grpc') {
                testK8sCesControl()
            }

            stageAutomaticRelease(makefile, docker)
        } catch (Exception e) {
            k3d.collectAndArchiveLogs()
            throw e as Throwable
        } finally {
            stage('Remove k3d cluster') {
                k3d.deleteK3d()
            }
        }
    }
}

String getCurrentCommit() {
    return sh(returnStdout: true, script: 'git rev-parse HEAD').trim()
}

private void testK8sCesControl() {
    sh "KUBECONFIG=${WORKSPACE}/k3d/.k3d/.kube/config make integration-test-bash"
    junit allowEmptyResults: true, testResults: 'target/bash-integration-test/*.xml'
}

void stageStaticAnalysisReviewDog() {
    String commitSha = getCurrentCommit()
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

void stageAutomaticRelease(Makefile makefile, Docker docker) {
    if (gitflow.isReleaseBranch()) {
        String releaseVersion = gitWrapper.getSimpleBranchName()
        String version = makefile.getVersion()

        stage('Build & Push Image') {
            def dockerImage = docker.build("cloudogu/${repositoryName}:${version}")

            docker.withRegistry('https://registry.hub.docker.com/', 'dockerHubCredentials') {
                dockerImage.push("${version}")
            }
        }

        stage('Sign after Release') {
            gpg.createSignature()
        }

        stage('Push Helm chart to Harbor') {
            docker
                    .image("golang:${goVersion}")
                    .mountJenkinsUser()
                    .inside("--volume ${WORKSPACE}:/go/src/${project} -w /go/src/${project}") {
                        // Package chart
                        make 'helm-package'
                        archiveArtifacts "${helmTargetDir}/**/*"

                        // Push chart
                        withCredentials([usernamePassword(credentialsId: 'harborhelmchartpush', usernameVariable: 'HARBOR_USERNAME', passwordVariable: 'HARBOR_PASSWORD')]) {
                            sh ".bin/helm registry login ${registry} --username '${HARBOR_USERNAME}' --password '${HARBOR_PASSWORD}'"

                            sh ".bin/helm push ${helmChartDir}/${repositoryName}-${version}.tgz oci://${registry}/${registry_namespace}/"
                        }
                    }
        }

        stage('Finish Release') {
            gitflow.finishRelease(releaseVersion, productionReleaseBranch)
        }

        stage('Add Github-Release') {
            releaseId = github.createReleaseWithChangelog(releaseVersion, changelog, productionReleaseBranch)
        }
    }
}

void make(String makeArgs) {
    sh "make ${makeArgs}"
}