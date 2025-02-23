FROM ghcr.io/cirruslabs/flutter:3.27.3

WORKDIR /app

COPY . .

RUN flutter config --enable-web \
    && flutter create . --platforms web \
    && flutter pub get \
    && flutter build web

EXPOSE 3000

CMD ["flutter", "run", "-d", "web-server", "--web-port", "3000", "--web-hostname", "0.0.0.0", "--release"] 