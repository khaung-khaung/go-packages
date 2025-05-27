# go_packages

make test TEST_FUNC=TestCSVDownload
make test TEST_FUNC=TestGet

make test TEST_FUNC=TestRabbitProduce
make test TEST_FUNC=TestRabbitConsume2
make test TEST_FUNC=TestRabbitConsume1
make test TEST_FUNC=TestGetCSVFile
make test TEST_FUNC=TestBulkUpdateFailures
make test TEST_FUNC=TestSMTPSend
make test TEST_FUNC=TestPlanSummary
make test TEST_FUNC=TestPlanDetail





















services:
  zookeeper:
    image: confluentinc/cp-zookeeper:7.3.2
    container_name: zookeeper
    ports:
      - "2181:2181"
    environment:
      - ALLOW_ANONYMOUS_LOGIN=yes
    networks:
      - app-network

  kafka:
    image: confluentinc/cp-kafka:7.3.2
    container_name: kafka
    ports:
      - "9092:9092"  # Host:Container mapping for external access
    environment:
      - KAFKA_CFG_LISTENERS=PLAINTEXT_INTERNAL://:9092,PLAINTEXT_EXTERNAL://:9092
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT_INTERNAL://kafka:9092,PLAINTEXT_EXTERNAL://localhost:9092
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=PLAINTEXT_INTERNAL:PLAINTEXT,PLAINTEXT_EXTERNAL:PLAINTEXT
      - KAFKA_CFG_INTER_BROKER_LISTENER_NAME=PLAINTEXT_INTERNAL  # Critical fix
      - KAFKA_CFG_ZOOKEEPER_CONNECT=zookeeper:2181
      
    depends_on:
      - zookeeper
    networks:
      - app-network

  kafka-ui:
    image: provectuslabs/kafka-ui:latest
    container_name: kafka-ui
    ports:
      - "9999:8080"
    environment:
      - KAFKA_CLUSTERS_0_NAME=local
      - KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS=kafka:9092  # Internal Docker address
 

  
    depends_on:
      - kafka
      - zookeeper
    networks:
      - app-network

networks:
  app-network:
    driver: bridge


