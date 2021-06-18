FROM ubuntu

ARG USER=model
ARG UID=1000
ARG GID=1000
ARG PASSWORD=cryptohacker321
RUN useradd --create-home ${USER} --uid=${UID} && echo "${USER}:${PASSWORD}" | chpasswd

RUN set -xe \
	&& apt-get upgrade \
 	&& apt-get update \
	&& apt-get install -y python3 \
	&& apt-get install -y python3-pip


ENV WDIR=/home/${USER}
COPY requirements.txt ${WDIR}/requirements.txt
RUN python3 -m pip install --upgrade pip
RUN pip install -r ${WDIR}/requirements.txt

WORKDIR $WDIR
COPY code/ .

# Specific to my setup
ENV WDIR=${WDIR}/dataset
RUN tar -xf ${WDIR}/nmnist.tar.gz -C ${WDIR} \
	&& rm ${WDIR}/nmnist.tar.gz

USER ${UID}:${GID}
