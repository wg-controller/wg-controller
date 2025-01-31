pipeline {
    agent any

    options {
        disableConcurrentBuilds(abortPrevious: true)
    }
    
    environment {
        GHCR_TOKEN      = credentials('gh-token-1')
        GHCR_USERNAME   = "wg-controller"   // GitHub organization or username
        IMAGE_NAME      = "wg-controller"   // Docker image name
        REPO            = "ghcr.io/${GHCR_USERNAME}/${IMAGE_NAME}"
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        
        stage('Get Tag') {
            steps {
                script {
                    def tag = sh(
                        script: "git describe --tags --abbrev=0",
                        returnStdout: true
                    ).trim()
                    
                    if (!tag) {
                        error "Error: No git tag found. Exiting."
                    }

                    env.IMAGE_TAG = tag
                    echo "Using image tag: ${env.IMAGE_TAG}"
                }
            }
        }

        stage('Login to GHCR') {
            steps {
                script {
                    sh """
                        echo ${GHCR_TOKEN_PSW} | docker login ghcr.io -u ${GHCR_USERNAME} --password-stdin
                    """
                }
            }
        }

        stage('Initialize Buildx') {
            steps {
                script {
                    sh """
                        docker buildx inspect container-builder || \
                            docker buildx create --name container-builder --driver docker-container --bootstrap
                        docker buildx use container-builder
                    """
                }
            }
        }


        stage('Build and Push Docker Image') {
            steps {
                script {
                    sh """
                        docker buildx build \
                            --platform linux/amd64 \
                            --build-arg GHCR_TOKEN=${GHCR_TOKEN_PSW} \
                            --build-arg IMAGE_TAG=${IMAGE_TAG} \
                            --build-arg IMAGE_NAME=${IMAGE_NAME} \
                            -t ${REPO}:${IMAGE_TAG} \
                            -t ${REPO}:latest \
                            . --push
                    """
                }
            }
        }

        stage('Cleanup') {
            steps {
                script {
                    sh """
                        docker system prune -f
                    """
                }
            }
        }
    }
    
    post {
        success {
            echo "Publish Script Complete"
        }
        failure {
            echo "Publish Script Failed"
        }
    }
}
