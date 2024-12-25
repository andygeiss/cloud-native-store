set dotenv-load

# Build
build:
    @podman build -t cloud-native-store .

# Generate an encryption key.
genkey:
    @go run cmd/genkey/main.go

# Run the service.
run:
    @podman run -p 8080:8080 \
        -e ENCRYPTION_KEY=$ENCRYPTION_KEY \
        -e PORT=8080 \
        cloud-native-store

# Test the Go sources (Units).
test:
    @go test -v -coverprofile=.coverprofile.out ./internal/app/core/services/...

# Set up "Cloud Build" according to https://cloud.google.com/build/docs/build-push-docker-image.
# Check if billing is enabled at: https://cloud.google.com/billing/docs/how-to/verify-billing-enabled#confirm_billing_is_enabled_on_a_project
cloud-build-setup:
    # Enable the Artifact Registry, Cloud Build and Compute Engine APIs.
    @gcloud services enable \
        artifactregistry.googleapis.com \
        cloudbuild.googleapis.com \
        compute.googleapis.com

    # Add the artifactregistry.writer role.
    @gcloud projects add-iam-policy-binding $GCP_PROJECT_ID \
        --member=serviceAccount:$(gcloud projects describe $GCP_PROJECT_ID \
        --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
        --role="roles/artifactregistry.writer"

    # Add the storage.admin role.
    @gcloud projects add-iam-policy-binding $GCP_PROJECT_ID \
        --member=serviceAccount:$(gcloud projects describe $GCP_PROJECT_ID \
        --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
        --role="roles/storage.admin"

    # Add the iam.serviceAccountUser role, which includes the actAspermission to deploy to the runtime.
    @gcloud iam service-accounts add-iam-policy-binding $(gcloud projects describe $GCP_PROJECT_ID \
        --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
        --member=serviceAccount:$(gcloud projects describe $GCP_PROJECT_ID \
        --format="value(projectNumber)")-compute@developer.gserviceaccount.com \
        --role="roles/iam.serviceAccountUser" \
        --project=$GCP_PROJECT_ID

    # Create a new Docker repository.
    @gcloud artifacts repositories create $GCP_DOCKER_REPOSITORY --repository-format=docker \
        --location=$GCP_REGION --description="Docker repository"
    @gcloud artifacts repositories list

# Build a docker image using Google Cloud Build.
cloud-build:
    @gcloud builds submit --region=$GCP_REGION \
        --tag $GCP_REGION-docker.pkg.dev/$GCP_PROJECT_ID/$GCP_DOCKER_REPOSITORY/$GCP_DOCKER_IMAGE

# Set up Google Cloud SDK.
cloud-cli-setup:
    # Login and set up the project.
    @gcloud auth login

    # Set the project and region.
    @gcloud config set project $GCP_PROJECT_ID
    @gcloud config set functions/region $GCP_REGION
    @gcloud config set run/region $GCP_REGION

    # Disable usage reporting.
    @gcloud config set disable_usage_reporting false

    # Update the components.
    @gcloud components update

# Set up Google Cloud Run.
cloud-run-setup:
    # Enable the Cloud Run APIs.
    @gcloud services enable \
        run.googleapis.com

# Run a docker container based on a Google Cloud Build.
cloud-run:
    # Deploy the service.
    @gcloud run deploy $GCP_SERVICE \
        --image $GCP_REGION-docker.pkg.dev/$GCP_PROJECT_ID/$GCP_DOCKER_REPOSITORY/$GCP_DOCKER_IMAGE \
        --allow-unauthenticated

    # Make service public accessible.
    @gcloud run services add-iam-policy-binding $GCP_SERVICE \
        --member="allUsers" \
        --role="roles/run.invoker"
