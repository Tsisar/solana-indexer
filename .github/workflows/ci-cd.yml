name: Solana Indexer CI/CD

on:
  pull_request:
    types: [closed]
    branches:
      - dev
      - master
  workflow_dispatch:

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    environment: ${{ github.ref_name == 'master' && 'prod' || 'dev' }}

    steps:
      - name: Checkout code
        uses: actions/checkout@v2

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Extract latest Git tag
        id: git_tag
        run: echo "tag=$(git describe --tags --abbrev=0)" >> "$GITHUB_OUTPUT"

      - name: Docker login
        run: echo "${{ secrets.DOCKER_PASSWORD }}" | docker login -u "${{ secrets.DOCKER_USERNAME }}" --password-stdin

      - name: Set image tag based on branch
        id: image_tag
        run: |
          RAW_TAG=${{ steps.git_tag.outputs.tag }}
          if [ "${{ github.ref_name }}" = "master" ]; then
            echo "tag=$RAW_TAG" >> "$GITHUB_OUTPUT"
          else
            echo "tag=${RAW_TAG}-dev" >> "$GITHUB_OUTPUT"
          fi

      - name: Build and push multi-platform Docker image
        run: |
          IMAGE="intothefathom/solana-indexer.vaults"
          TAG="${{ steps.image_tag.outputs.tag }}"
          docker buildx build \
            --platform linux/amd64,linux/arm64 \
            -t $IMAGE:$TAG \
            -t $IMAGE:latest \
            --push .

      - name: Sync Argo CD and wait for completion
        run: |
          APP_NAME="solana-indexer"
          IMAGE_TAG="${{ steps.image_tag.outputs.tag }}"
          MAX_RETRIES=3
          RETRY_DELAY=10
          RETRIES=0

          until [ $RETRIES -ge $MAX_RETRIES ]
          do
            docker run --rm \
              -e ARGOCD_AUTH_TOKEN=${{ secrets.ARGOCD_AUTH_TOKEN }} \
              argoproj/argocd:v2.6.15 \
              /bin/sh -c \
              "argocd app set $APP_NAME \
              --server ${{ secrets.ARGOCD_API_URL }} \
              --grpc-web \
              --parameter image.tag=$IMAGE_TAG && \
              argocd app wait $APP_NAME \
              --server ${{ secrets.ARGOCD_API_URL }} \
              --grpc-web \
              --operation && \
              argocd app sync $APP_NAME \
              --server ${{ secrets.ARGOCD_API_URL }} \
              --grpc-web \
              --timeout 180" && break

            RETRIES=$((RETRIES+1))
            echo "Retrying... ($RETRIES/$MAX_RETRIES)"
            sleep $RETRY_DELAY
          done

          if [ $RETRIES -eq $MAX_RETRIES ]; then
            echo "Failed to sync after $MAX_RETRIES attempts"
            exit 1
          fi