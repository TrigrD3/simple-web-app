pipeline {
    agent any

    environment {
        DOCKER_HUB_CREDENTIALS = 'cab54931-41e6-4f0f-be37-aaaf1e94168e' // Replace with your Docker Hub credential ID in Jenkins
        IMAGE_NAME = 'trigrd/simple-web-app' // Replace with your Docker Hub username
        IMAGE_TAG = "${env.BUILD_NUMBER}"
        HELM_CHART_PATH = '/Users/arya/Downloads/Repo/simple-web-app/simple-web-app'
    }

    stages {
        stage('Build and Push Docker Image') {
            steps {
                script {
                    // Login to Docker Hub
                    withCredentials([usernamePassword(credentialsId: env.DOCKER_HUB_CREDENTIALS, passwordVariable: 'DOCKER_PASSWORD', usernameVariable: 'DOCKER_USERNAME')]) {
                        sh "echo \"$DOCKER_PASSWORD\" | docker login -u \"$DOCKER_USERNAME\" --password-stdin"
                    }

                    // Build Docker image
                    sh "docker build -t ${IMAGE_NAME}:${IMAGE_TAG} ."

                    // Push Docker image
                    sh "docker push ${IMAGE_NAME}:${IMAGE_TAG}"
                }
            }
        }

        stage('Deploy to Kubernetes') {
            steps {
                script {
                    // Update Helm values with the new image tag
                    sh "sed -i 's|tag: latest|tag: ${IMAGE_TAG}|g' ${HELM_CHART_PATH}/values.yaml"

                    // Deploy with Helm
                    sh "helm upgrade --install simple-web-app ${HELM_CHART_PATH} --set image.repository=${IMAGE_NAME} --set image.tag=${IMAGE_TAG}"
                }
            }
        }
    }

    post {
        always {
            // Clean up Docker images (optional)
            sh "docker rmi ${IMAGE_NAME}:${IMAGE_TAG} || true"
        }
    }
}
