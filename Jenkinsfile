pipeline {
    agent any
    
    environment {
        // Jenkins credential ID for GHCR token
        GHCR_TOKEN      = credentials('gh-token-1')
        
        // GHCR variables
        GHCR_USERNAME   = "wg-controller"   // GitHub organization or username
        IMAGE_NAME      = "wg-controller"   // Docker image name
        REPO            = "ghcr.io/${GHCR_USERNAME}/${IMAGE_NAME}"
    }

    stages {
        stage('Checkout') {
            steps {
                // If this pipeline is triggered from a Multibranch Pipeline or
                // a Pipeline project configured to use this repo/branch, 
                // checkout scm will clone the relevant branch/tag
                checkout scm
            }
        }
        
        stage('Determine Tag') {
            steps {
                script {
                    // Capture the most recent git tag
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
                    // Docker login to GitHub Container Registry
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
                        # If the builder doesn't exist, create it; otherwise, just use it
                        docker buildx inspect container-builder || \
                            docker buildx create --name container-builder --driver docker-container --bootstrap

                        # Switch context to the container-builder
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
                            --platform linux/amd64 --platform linux/arm64 \
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
                    // Cleanup Docker Buildx
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
