pipeline {
    // Run the entire pipeline inside a Docker container that has the docker CLI
    agent {
        docker {
            image 'docker:latest' // Use an official image that has docker tools
            args '-v /var/run/docker.sock:/var/run/docker.sock' // Mount the host's Docker socket
        }
    }

    environment {
        DOCKER_HUB_CREDENTIALS = 'cab54931-41e6-4f0f-be37-aaaf1e94168e' // Stays the same
        IMAGE_NAME             = 'trigrd/simple-web-app' // Stays the same
        IMAGE_TAG              = "${env.BUILD_NUMBER}"
        // IMPORTANT: The path inside the container will be different.
        // Update this path to where your Helm chart is relative to your git repo root.
        // Let's assume it's in a 'helm' directory in your repo.
        HELM_CHART_PATH        = 'simple-web-app/simple-web-app' // UPDATE THIS PATH
    }

    stages {
        stage('Build and Push Docker Image') {
            steps {
                script {
                    // Login to Docker Hub
                    withCredentials([usernamePassword(credentialsId: DOCKER_HUB_CREDENTIALS, passwordVariable: 'DOCKER_PASSWORD', usernameVariable: 'DOCKER_USERNAME')]) {
                        sh "echo \"$DOCKER_PASSWORD\" | docker login -u \"$DOCKER_USERNAME\" --password-stdin"
                    }

                    // Build and Push are now executed inside the container agent
                    sh "docker build -t ${IMAGE_NAME}:${IMAGE_TAG} ."
                    sh "docker push ${IMAGE_NAME}:${IMAGE_TAG}"
                }
            }
        }

        stage('Deploy to Kubernetes') {
            // This stage will fail unless you also add the 'helm' and 'sed' tools to your agent.
            // For now, let's focus on fixing the Docker part. To run this stage,
            // you would need to use a custom Docker image that has docker, helm, and sed installed.
            steps {
                echo "Skipping deploy for now to focus on Docker fix."
                // script {
                //     sh "sed -i 's|tag: latest|tag: ${IMAGE_TAG}|g' ${HELM_CHART_PATH}/values.yaml"
                //     sh "helm upgrade --install simple-web-app ${HELM_CHART_PATH} --set image.repository=${IMAGE_NAME} --set image.tag=${IMAGE_TAG}"
                // }
            }
        }
    }

    post {
        always {
            sh "docker rmi ${env.IMAGE_NAME}:${env.IMAGE_TAG} || true"
        }
    }
}