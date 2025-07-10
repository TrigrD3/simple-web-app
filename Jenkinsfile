pipeline {
    // This agent block defines the pod where our pipeline will run.
    // It uses the template configured in the Jenkins UI.
    agent {
        kubernetes {
            label 'docker-builder'          // Must match the label in your Pod Template
            defaultContainer 'builder'      // The container where steps run. This now has docker & helm.
        }
    }

    // Environment variables are available throughout the pipeline via the 'env' object.
    environment {
        // User-defined variables
        DOCKER_HUB_CREDENTIALS = 'cab54931-41e6-4f0f-be37-aaaf1e94168e' // Your Jenkins credential ID for Docker Hub
        IMAGE_NAME             = 'trigrd/simple-web-app'                // Your Docker Hub repository name
        IMAGE_TAG              = "${env.BUILD_NUMBER}"                  // Use the Jenkins build number for a unique image tag
        HELM_CHART_PATH        = 'simple-web-app'                       // The relative path to your Helm chart directory in the repo
        HELM_RELEASE_NAME      = 'go-web-app'                       // The name for your Helm deployment release
        KUBE_NAMESPACE         = 'simple-web-app'                              // The Kubernetes namespace to deploy the application to

        // Docker-in-Docker (dind) configuration for the agent pod
        DOCKER_HOST        = 'tcp://localhost:2376'                     // Tell the Docker client to connect to the dind secure port
        DOCKER_CERT_PATH   = '/certs/client'                            // Tell the client where to find the TLS certs from the shared volume
        DOCKER_TLS_VERIFY  = '1'                                        // Enforce TLS communication with the dind sidecar
    }

    stages {
        // This initial stage is a safety check to ensure the dind sidecar is ready.
        stage('Wait for Docker Daemon') {
            steps {
                sh 'echo "Waiting for Docker daemon to be ready..."; while ! docker info; do sleep 1; done; echo "Docker daemon is ready."'
            }
        }

        // This stage builds the application's Docker image and pushes it to your registry.
        stage('Build and Push Docker Image') {
            steps {
                // Use withCredentials to securely inject Docker Hub username and password
                withCredentials([usernamePassword(credentialsId: env.DOCKER_HUB_CREDENTIALS, passwordVariable: 'DOCKER_PASSWORD', usernameVariable: 'DOCKER_USERNAME')]) {
                    sh "echo \"$DOCKER_PASSWORD\" | docker login -u \"$DOCKER_USERNAME\" --password-stdin"
                }

                // Build the image using the Dockerfile in the root of the repository
                sh "docker build -t ${env.IMAGE_NAME}:${env.IMAGE_TAG} ."

                // Push the newly built image to Docker Hub
                sh "docker push ${env.IMAGE_NAME}:${env.IMAGE_TAG}"
            }
        }

        // This stage uses Helm to deploy your application to the Kubernetes cluster.
        stage('Deploy to Kubernetes') {
            steps {
                script {
                    // Use 'helm upgrade --install' to either create the release or update it if it already exists.
                    // We use '--set' to dynamically inject the correct image repository and tag into the Helm chart,
                    // which is cleaner than modifying files with 'sed'.
                    sh """
                        helm upgrade --install ${env.HELM_RELEASE_NAME} ./${env.HELM_CHART_PATH} \
                             --namespace ${env.KUBE_NAMESPACE} \
                             --set image.repository=${env.IMAGE_NAME} \
                             --set image.tag=${env.IMAGE_TAG} \
                             --wait
                    """
                }
            }
        }
    }

    // The 'post' block runs after all stages are completed.
    post {
        // 'always' will run regardless of whether the pipeline succeeded or failed.
        always {
            // Good practice to clean up the build environment.
            // This removes the large Docker image from the agent pod to free up space.
            echo "Cleaning up Docker image from the build agent..."
            sh "docker rmi ${env.IMAGE_NAME}:${env.IMAGE_TAG} || true"
        }
    }
}