pipeline {
    // Use the Kubernetes plugin to spin up a pod based on the 'docker-builder' template
    agent {
        kubernetes {
            label 'docker-builder' // Must match the label you set in the Pod Template
            defaultContainer 'builder' // The container where 'sh' steps will run
        }
    }

    environment {
        // This remains the same
        DOCKER_HUB_CREDENTIALS = 'cab54931-41e6-4f0f-be37-aaaf1e94168e'
        IMAGE_NAME = 'trigrd/simple-web-app'
        IMAGE_TAG = "${env.BUILD_NUMBER}"
        // The DOCKER_HOST variable tells the docker client (in the 'builder' container)
        // to connect to the dind container over the local pod network.
        DOCKER_HOST = 'tcp://localhost:2376'
    }

    stages {
        stage('Wait for Docker Daemon') {
            steps {
                // The dind container can take a few seconds to start.
                // This step ensures the Docker daemon is ready before we try to use it.
                sh 'while ! docker info; do echo "Waiting for docker daemon..."; sleep 1; done'
            }
        }
        
        stage('Build and Push Docker Image') {
            steps {
                // These steps will now run inside the 'builder' container
                // and communicate with the 'dind' container.
                withCredentials([usernamePassword(credentialsId: DOCKER_HUB_CREDENTIALS, passwordVariable: 'DOCKER_PASSWORD', usernameVariable: 'DOCKER_USERNAME')]) {
                    sh "echo \"$DOCKER_PASSWORD\" | docker login -u \"$DOCKER_USERNAME\" --password-stdin"
                }
                sh "docker build -t ${env.IMAGE_NAME}:${env.IMAGE_TAG} ."
                sh "docker push ${env.IMAGE_NAME}:${env.IMAGE_TAG}"
            }
        }

        // The 'Deploy' stage will need 'helm' installed in the build container.
        // To fix that, you'd use a custom image like 'alpine/helm' instead of 'docker:latest'
        // or build your own image with all necessary tools.
        stage('Deploy to Kubernetes') {
            steps {
                echo "Skipping deploy for now, as helm is not in the 'docker:latest' image."
            }
        }
    }

    post {
        always {
            // This now works correctly within the pod's context.
            sh "docker rmi ${env.IMAGE_NAME}:${env.IMAGE_TAG} || true"
        }
    }
}