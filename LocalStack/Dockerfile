FROM localstack/localstack-pro:latest

COPY ./aws/init_bucket.sh /etc/localstack/init/ready.d/init_bucket.sh

RUN sed -i 's/\r$//' /etc/localstack/init/ready.d/init_bucket.sh

RUN chmod +x /etc/localstack/init/ready.d/init_bucket.sh