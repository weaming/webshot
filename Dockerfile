FROM ubuntu:16.04
RUN apt-get update -y && DEBIAN_FRONTEND=noninteractive apt-get -y install build-essential xorg libssl-dev libxrender-dev wget
RUN apt-get -y install fontconfig libjpeg-turbo8 xfonts-75dpi fonts-noto-cjk
RUN wget https://github.com/wkhtmltopdf/wkhtmltopdf/releases/download/0.12.5/wkhtmltox_0.12.5-1.xenial_amd64.deb && \
    dpkg -i wkhtmltox_0.12.5-1.xenial_amd64.deb; rm wkhtmltox_0.12.5-1.xenial_amd64.deb

COPY $PWD/webshot /usr/local/bin/
WORKDIR /app
EXPOSE 80
CMD ["/usr/local/bin/webshot"]
