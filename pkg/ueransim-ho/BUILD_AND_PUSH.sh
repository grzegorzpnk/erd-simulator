GNB_NAME=ueransim-gnb
UE_NAME=ueransim-ue

if [ -z "${DOCKER_REPO}" ] || [ -z "${UERANSIM_TAG}" ]; then
    echo "DOCKER_REPO and UERANSIM_TAG variables NOT set. Aborting..."
fi

echo "Compile code..."

make

echo "Building ${DOCKER_REPO}${GNB_NAME}:${UERANSIM_TAG}"
sudo docker build . -f ./Dockerfile.gnb -t ${DOCKER_REPO}${GNB_NAME}:${UERANSIM_TAG}

echo "Pushing ${DOCKER_REPO}${GNB_NAME}:${UERANSIM_TAG}"
sudo docker push ${DOCKER_REPO}${GNB_NAME}:${UERANSIM_TAG}

echo "Building ${DOCKER_REPO}${UE_NAME}:${UERANSIM_TAG}"
sudo docker build . -f ./Dockerfile.ue -t ${DOCKER_REPO}${GNB_NAME}:${UERANSIM_TAG}

echo "Pushing ${DOCKER_REPO}${UE_NAME}:${UERANSIM_TAG}"
sudo docker push ${DOCKER_REPO}${UE_NAME}:${UERANSIM_TAG}

