## Architecture
<img width="1329" alt="Screenshot 2023-10-04 at 17 36 52" src="https://github.com/mohrizkimaulana/temukanjalan/assets/54239628/062e6c15-d0b2-4b11-9031-6a4c14b7b1d0">

# Overview:

Welcome to our CICD Pipeline with GitLab, DockerHub, ArgoCD, Kubernetes (Minikube), and Helm Charts. This CICD pipeline streamlines the development and deployment process for our two distinct backend services, one written in `Go` and the other in `Node.js`.

# Components
`GitLab Repository`: Our code resides here, organized into two service folders: app_go and app_node.

`GitLab Actions (CI)`: GitHub workflows are triggered automatically whenever changes are made in the service folders. This allows us to ensure that our applications are always up to date with the latest code.

`DockerHub Image Registry`: Docker images of our applications are stored here. Images are built, tagged, and pushed automatically for each service whenever changes occur.

`Kubernetes (Minikube)`: Our Kubernetes cluster, run locally using Minikube, provides the infrastructure for our applications. It's where our applications are deployed and managed.

`Helm Charts`: Helm charts in the helm-go and helm-node folders help manage our Kubernetes manifest files. They allow us to define and version our application configurations.

`ArgoCD (CD)`: ArgoCD is the glue that automates Continuous Deployment. It watches our Git repository, and when changes occur, it syncs our applications with the desired state defined in the Helm charts. It's our key to GitOps.

# Workflow

1. __`Code Changes`__: Whenever you make changes to the Go or Node.js code in the respective service folders (app_go or app_node), the corresponding GitHub workflow is triggered. This starts the CICD process.

2. __`Image Building`__: The workflow logs into DockerHub, builds Docker images, and tags them with the latest GitHub commit SHA.

3. __`Image Push`__: It then pushes the newly created images to DockerHub, ensuring they are readily available for deployment.

4. __`Configuration Update`__: The final step is to update the application configuration defined in the Helm charts. For example, the environment variable for the API version is updated to match the latest image tags.

5. __`ArgoCD Sync`__: After these updates, ArgoCD automatically detects changes in our Git repository and synchronizes the applications' Kubernetes manifests. This ensures our applications are always in the desired state, reflecting the latest code changes.


# Infrastructure, Code and Pipeline Provision

Requirements:

- Docker
  
- Minikube 
  
- ArgoCD
  
- Helm


Installation steps:

1. __Install docker__

   (In this case, we are using docker desktop for Mac), but if you have downloaded the docker, you could skip this step
  ```
  https://www.docker.com/products/docker-desktop/
  ```

2. __Install minikube__
  minikube is local Kubernetes, focusing on making it easy to learn and develop for Kubernetes.

  ```bash
  curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-darwin-amd64
  sudo install minikube-darwin-amd64 /usr/local/bin/minikube
  ```

3. __Start minikube__
   In this installation, we are using the docker driver. list of drivers can be read here: https://minikube.sigs.k8s.io/docs/drivers/)

   ```
   minikube start --driver=docker
   ```
   The successful execution output is similar to below.

   Output:
   ```
   ....
   🔎  Verifying ingress addon...
   🌟  Enabled addons: ingress-dns, storage-provisioner, default-storageclass, ingress
   🏄  Done! kubectl is now configured to use "minikube" cluster and "default" namespace by default
   ```
  
   Great!! , let's build new container within our new cluster :)

5. __Setting-up argocd__


   ```bash
   kubectl create namespace argocd
   kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
   ```

6. __Setting-up Helm__
   
   Helm simplifies deploying and managing applications on Kubernetes. Let’s configure Helm and add the necessary repositories:
   
   Step 1: Install Helm:
   ```
   curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
   chmod 700 get_helm.sh
   ./get_helm.sh
   ```
   Step 2: Add the official Helm repository:

   ```
   helm repo add stable https://charts.helm.sh/stable
   ```
   Step 3: Update the Helm repositories:
   ```
   helm repo update
   ```

7. __Setting-Up Application__

    To configure and set up the application, please refer to the respective `app_go` and `app_node` folders located in the root of this repository. These folders contain all the necessary files and configurations for their respective services. You can find detailed instructions and documentation within each of these folders to help you get started with setting up and running the Go and Node.js applications.
   
9. __Setting Up the Helm Chart__

   As we have two separate backend services, we've created two distinct Helm charts, named helm-go and helm-node. Each of these charts is tailored to manage the deployment and configuration of their respective services, making it easier to handle and customize the settings for each backend application independently.

   Step 1: Initialize go & node Helm chart:
   ```
   helm create helm-go
   helm create helm-node
   ```

   This will create a new directory named helm-go and helm-node with the basic structure of a Helm chart.

   Step 2: Modify go Helm chart files according to our application’s requirements.
   - Update `helm-go/templates/deployment.yaml`:

     Replace the content of the file with the following code:
     
     ```
     apiVersion: apps/v1
     kind: Deployment
     metadata:
       name: go-app-deployment
       namespace: default
     spec:
       replicas: 3
       selector:
         matchLabels:
           app: go-backend
       template:
         metadata:
           labels:
             app: go-backend
         spec:
           containers:
           - name: go-backend
             image: mohrizkimaulana/app-go:{{ .Values.env.APP_VERSION }}
             resources:
               limits:
                 memory: "128Mi"
                 cpu: "500m"
             ports:
             - containerPort: 3000

     ```

     - Update `helm-go/templates/service.yaml`:

       Replace the content of the file with the following code:
    
     ```
      apiVersion: v1
      kind: Service
      metadata:
        name: go-backend-service
        namespace: default
      spec:
        selector:
          app: go-backend
        ports:
        - protocol: TCP
          port: 80
          targetPort: 3000

     ```
      - Update `helm-go/templates/ingress.yaml`:
        
        Replace the content of the file with the following code:

     ```
      apiVersion: networking.k8s.io/v1
      kind: Ingress
      metadata:
        name: go-app-ingress
        namespace: default
        annotations:
          nginx.ingress.kubernetes.io/rewrite-target: /$1
      spec:
        rules:
        - host: pintu.doraemon.local
          http:
            paths:
            - pathType: Prefix
              path: "/go"
              backend:
                service:
                  name: go-backend-service
                  port: 
                    number: 80

     ```

   Step 3: Modify node Helm chart files according to our application’s requirements.
   - Update `helm-node/templates/deployment.yaml`:

     Replace the content of the file with the following code:
     
     ```
      apiVersion: apps/v1
      kind: Deployment
      metadata:
        name: nodejs-app-deployment
        namespace: default
      spec:
        selector:
          matchLabels:
            app: nodejs-backend
        template:
          metadata:
            labels:
              app: nodejs-backend
          spec:
            containers:
            - name: nodejs-backend
              image: mohrizkimaulana/app-node:{{ .Values.env.APP_VERSION }}
              resources:
                limits:
                  memory: "128Mi"
                  cpu: "500m"
              ports:
              - containerPort: 3000


     ```

     - Update `helm-node/templates/service.yaml`:

       Replace the content of the file with the following code:
    
     ```
      apiVersion: v1
      kind: Service
      metadata:
        name: nodejs-backend-service
        namespace: default
      spec:
        selector:
          app: nodejs-backend
        ports:
        - protocol: TCP
          port: 80
          targetPort: 3000


     ```
      - Update `helm-go/templates/ingress.yaml`:
        
        Replace the content of the file with the following code:

     ```
      apiVersion: networking.k8s.io/v1
      kind: Ingress
      metadata:
        name: node-app-ingress
        namespace: default
        annotations:
          nginx.ingress.kubernetes.io/rewrite-target: /$1
      spec:
        rules:
        - host: pintu.doraemon.local
          http:
            paths:
            - pathType: Prefix
              path: "/node"
              backend:
                service:
                  name: nodejs-backend-service
                  port: 
                    number: 80

     ```

     Step 4: Change values.yaml in both of helm-go and helm-node with this value:

     ```
      env:
      APP_VERSION: <desired_application_version> 
     ```

10. __Setting Up GitHub Repository and Linking to the Application__
    
    Step 1: Create a GitHub Repository
      - Visit https://github.com and log in to your account.
      - Click on "New" to create a new repository.
      - Provide a name for your repository, e.g., "ngetes"
      - Optionally, add a description and choose repository visibility settings.
      - Click "Create repository."

    Step 2: Initialize Git and Link the Repository
  
      - Open your terminal or command prompt.
      - Navigate to the root directory of your application.
      - Run `git init` to initialize Git in the directory.
      - Run `git remote add origin <repository_url>` to set the GitHub repository as your remote origin.
      - Replace <repository_url> with your GitHub repository's URL.
      - Run these commands:
          - `git add .` to stage the changes.
          - `git commit -m "Initial commit"` to commit with a meaningful message.
          - `git push -u origin main` to push the code to the repository.
          - 
        Note: If your default branch is named differently (e.g., "master"), replace "main" with your branch name.
        Your application code is now linked to the GitHub repository, and the initial commit is pushed.

11. __Setting Up Argo CD for Application Deployment__
    
    Step 1: Access the Argo CD UI
      - Run code below to access the Argo CD UI locally.
        
        ```
        kubectl -n argocd port-forward svc/argocd-server 8080:443
        ```
      - Visit http://localhost:8080 in your browser to access the Argo CD UI.
        
        ```
        http://localhost:8080
        ```
        
    Step 2: Log in to the Argo CD UI
    
      - Get the Argo CD admin password with
        
        ```
        kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d.
        ```
        
      - Copy the password and paste it in the Argo CD login page.
      - Log in with the username `admin` and the password you copied.
        
    Step 3: Create an Application in Argo CD
    
      - Go to the "Applications" tab.
      - Click "New Application."
      - Configure the application settings, including the `name`, `project`, `repository URL`, `path to the Helm chart directory`, `target cluster`, `namespace`.
      - Enable "Auto-Sync."
      - Click "Create" to create the application.
        
    Step 4: Verify the Deployment in Argo CD
      - After creating the application, wait for Argo CD to automatically synchronize it.
      - You'll see the deployment status and resources listed on the application details page.
      - Ensure that application resources are successfully deployed and running as expected.
        
12. __Setting Up GitHub Actions for Automated Deployment__
    
    Step 1: Create the GitHub Actions Workflow File
      - In your GitHub repository, create the directory .github/workflows.
      - Create a new file named `go.yaml` and `node.yaml` for the workflow.
        
    Step 2: Add the Workflow Code to go.yaml
      

    ```
    name: GO CI

    on:
      push:
        paths:
          - 'app_go/**'

    env:
      DOCKERHUB_USERNAME: ${{ secrets.DOCKER_USERNAME }}
      DOCKERHUB_KEY: ${{ secrets.DOCKER_KEY }}
      IMAGE_NAME: app-go

    jobs:
      build-and-deploy:
        runs-on: ubuntu-latest

        steps:
          - name: Checkout code
            uses: actions/checkout@v2
    
          - name: Login to Docker Hub
            uses: docker/login-action@v1
            with:
              username: ${{ env.DOCKERHUB_USERNAME }}
              password: ${{ env.DOCKERHUB_KEY }}
    
          - name: Build Docker image
            run: docker build -t ${{ env.DOCKERHUB_USERNAME }}/${{ env.IMAGE_NAME }}:${{ github.sha }} --file ./app_go/Dockerfile .
    
          - name: Push Docker image
            run: docker push ${{ env.DOCKERHUB_USERNAME }}/${{ env.IMAGE_NAME }}:${{ github.sha }}
    
          - name: Update values.yaml
            run: |
              git pull
              cd helm/helm-go
              sed -i 's|APP_VERSION:.*|APP_VERSION: '${{ github.sha }}'|' values.yaml 
              git config --global user.name 'mohrizkimaulana'
              git config --global user.email 'maulana1507000@gmail.com'
              git add values.yaml
              git commit -m "Update values.yaml"
              git push

    ```

    Step 3: Add the Workflow Code to node.yaml

    ```
    name: NODE CI

    on:
      push:
        paths:
          - 'app_node/**'

    env:
      DOCKERHUB_USERNAME: ${{ secrets.DOCKER_USERNAME }}
      DOCKERHUB_KEY: ${{ secrets.DOCKER_KEY }}
      IMAGE_NAME: app-node

    jobs:
      build-and-deploy:
        runs-on: ubuntu-latest

        steps:
          - name: Checkout code
            uses: actions/checkout@v2
    
          - name: Login to Docker Hub
            uses: docker/login-action@v1
            with:
              username: ${{ env.DOCKERHUB_USERNAME }}
              password: ${{ env.DOCKERHUB_KEY }}
      
          - name: Build Docker image
            run: docker build -t ${{ env.DOCKERHUB_USERNAME }}/${{ env.IMAGE_NAME }}:${{ github.sha }} --file ./app_node/Dockerfile .
    
          - name: Push Docker image
            run: docker push ${{ env.DOCKERHUB_USERNAME }}/${{ env.IMAGE_NAME }}:${{ github.sha }}
    
          - name: Update values.yaml
            run: |
              git pull
              cd helm/helm-node
              sed -i 's|APP_VERSION:.*|APP_VERSION: '${{ github.sha }}'|' values.yaml 
              git config --global user.name 'GitHub Actions'
              git config --global user.email 'actions@github.com'
              git add values.yaml
              git commit -m "Update values.yaml"
              git push
    ```

13. __Configure Docker Hub and GitHub Secrets__
    1. In your GitHub repository, navigate to `"Settings"` and click on `"Secrets."` (maybe you need to create an environment first)
    2. Add the `DOCKER_USERNAME` and `DOCKER_KEY` secrets with your Docker Hub credentials (username & password docker hub)
    3. __Grant necessary permissions in the "Actions" settings under "Settings" for the workflow to run__. `Settings`>`Actions`>`General`>`Workflow Permissions`> `Enable Read and Write Permissions
    4. With these steps, your GitHub repository is now configured with the GitHub Actions workflow for automated deployment using Argo CD.

14. __Accessing Locally with Custom Domain and Minikube Tunnel (MacBook)__

    Step 1: Modify the Hosts File:

    To access a service locally using the custom domain "pintu.doraemon.local," you need to update your local machine's /etc/hosts file. Open the /etc/hosts file with administrative privileges (on Linux/Unix-based systems) and add the following line at the end of the file:

    ```
      127.0.0.1 pintu.doraemon.local

    ```

    This entry tells your computer to redirect requests for "pintu.doraemon.local" to your local machine (127.0.0.1).

    Step 2: Start Minikube Tunnel

    Minikube is a local Kubernetes cluster manager. To ensure that requests to your custom domain are correctly routed to services running in your Minikube cluster, you need to start Minikube's tunnel feature. Open a terminal or command prompt and run the following command:

    ```
      minikube tunnel

    ```

    The minikube tunnel command sets up a network route between your local machine and the Minikube cluster, allowing traffic to flow smoothly between them.

    Now, you can access services within your Minikube cluster using the custom domain `pintu.doraemon.local` as if they were running locally on your machine. This setup is particularly useful for local development and testing with custom domains.
    

15. When changes are pushed to the main branch, the workflow will automate Docker image building, values.yaml updating, and pushing changes to your repository.
   
    
