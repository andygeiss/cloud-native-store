set dotenv-load

cache-dir := ".cache"
cert-dir := ".tls"

# Set up "Cloud Build" according to https://cloud.google.com/build/docs/build-push-docker-image.
# Check if billing is enabled at: https://cloud.google.com/billing/docs/how-to/verify-billing-enabled#confirm_billing_is_enabled_on_a_project
cloud-build-setup:
    # Login and set up the project.
    @gcloud auth login
    @gcloud config set project $GCP_PROJECT_ID
    @gcloud config set disable_usage_reporting false
    @gcloud components update

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

# Run the service.
run:
    @go run cmd/main.go

# Set up the service.
# Create a local CA and sign a server certificate.
# This will only be used if domains = ["localhost"].
setup:
    @brew install mkcert
    @rm -rf {{cache-dir}} ; mkdir {{cache-dir}}
    @rm -rf {{cert-dir}} ; mkdir {{cert-dir}}
    @mkcert -install
    @mkcert -cert-file {{cert-dir}}/server.crt \
        -key-file {{cert-dir}}/server.key \
        localhost 127.0.0.1 ::1

# Test the Go sources (Units).
test:
    @go test -v -coverprofile=.coverprofile.out ./internal/app/core/services/...
