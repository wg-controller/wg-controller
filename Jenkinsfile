pipeline {
    agent any
    
    environment {
        // Jenkins credential ID for GHCR token
        GHCR_TOKEN      = credentials('a1063d11-4e4c-40bf-96b1-ce70f84a4542')
        
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
                        echo ${GHCR_TOKEN} | docker login ghcr.io -u ${GHCR_USERNAME} --password-stdin
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
                            --build-arg GHCR_TOKEN=${GHCR_TOKEN} \
                            --build-arg IMAGE_TAG=${IMAGE_TAG} \
                            --build-arg IMAGE_NAME=${IMAGE_NAME} \
                            -t ${REPO}:${IMAGE_TAG} \
                            -t ${REPO}:latest \
                            . --push
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
