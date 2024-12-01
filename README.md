# Crypto Exchange

A robust and scalable crypto exchange application built with Go, leveraging PostgreSQL, Redis, Cassandra, and Kafka. The application is fully containerized using Docker and orchestrated with Docker Compose.

## **Features**

- **Transaction Management**: Create and retrieve financial transactions.
- **Caching**: Utilize Redis for efficient data retrieval.
- **Distributed Storage**: Use Cassandra for scalable data storage.
- **Event Streaming**: Implement Kafka for real-time event processing.
- **Structured Logging**: Employ zerolog for performant and structured logs.
- **Configuration Management**: Manage configurations with Viper and YAML.

## **Prerequisites**

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Go](https://golang.org/dl/) (optional, for local development outside Docker)

## **Setup Instructions**

1. **Clone the Repository**

   ```bash
   git clone https://github.com/yourusername/crypto-exchange.git
   cd crypto-exchange
   ```

2. **Configure Application**

   Modify `config.yaml` as needed, especially the `api_keys` and other sensitive information. Ensure that the service names in `config.yaml` match those in `docker-compose.yml`.

3. **Build and Run with Docker Compose**

   Ensure Docker is running on your system.

   ```bash
   docker-compose up --build
   ```

   This command builds the Docker image for the Go application and starts all services defined in `docker-compose.yml`.

4. **Accessing the Application**

   The application will be accessible at `http://localhost:8080`.

5. **Testing API Endpoints**

   - **Create Transaction**

     ```bash
     curl -X POST http://localhost:8080/transactions \
     -H "Content-Type: application/json" \
     -d '{
       "id": "tx123",
       "amount": 150.75,
       "type": "deposit",
       "status": "pending"
     }'
     ```

   - **Get Transaction**

     ```bash
     curl http://localhost:8080/transactions/tx123
     ```

## **Project Components**

### **1. Configuration (`config/config.go` & `config.yaml`)**

Handles loading, parsing, and validating configurations using Viper and Validator.

### **2. Logging (`logger/logger.go` & `middleware/logger.go`)**

Implements structured logging with zerolog and logs each HTTP request.

### **3. Models (`models/transaction.go`)**

Defines the `Transaction` model representing financial transactions.

### **4. Services (`services/`)**

- **Database Service**: Manages PostgreSQL connections and migrations.
- **Redis Service**: Handles caching operations.
- **Cassandra Service**: Manages Cassandra connections and data operations.
- **Kafka Service**: Handles event publishing to Kafka.
- **Transaction Service**: Orchestrates creation and retrieval of transactions across all services.
- **Mock Transaction Service**: Provides a mock implementation for testing purposes.

### **5. Controllers (`controllers/transaction_controller.go`)**

Manages HTTP requests related to transactions, utilizing the Transaction Service.

### **6. Routes (`routes/routes.go`)**

Defines the API endpoints and associates them with controller handlers.

### **7. Docker Configuration (`Dockerfile` & `docker-compose.yml`)**

Containers for the Go application, PostgreSQL, Redis, Cassandra, and Kafka, managed via Docker Compose.

## **Best Practices Implemented**

- **Dependency Injection**: Services are injected into controllers to enhance testability and modularity.
- **Structured Logging**: Facilitates easier debugging and monitoring.
- **Configuration Management**: Centralizes configurations, making it easier to manage different environments.
- **Containerization**: Ensures consistent environments across development, testing, and production.
- **Data Consistency**: Uses transactions in PostgreSQL to ensure data integrity across services.

## **Extending the Application**

1. **Authentication & Authorization**: Implement JWT-based authentication to secure endpoints.
2. **User Management**: Develop endpoints and services for user registration, login, and profile management.
3. **Order Matching Engine**: Create a real-time order matching system for buy/sell orders.
4. **Wallet Integration**: Integrate cryptocurrency wallets for handling deposits and withdrawals.
5. **Monitoring & Alerts**: Incorporate monitoring tools like Prometheus and Grafana for system health insights.
6. **API Documentation**: Utilize Swagger or similar tools to document API endpoints.

## **Stopping the Application**

To stop all running containers, press `Ctrl+C` in the terminal where Docker Compose is running or execute:

```bash
docker-compose down
```

This command stops and removes all containers, networks, and volumes defined in `docker-compose.yml`.

## **Cleaning Up Docker Resources**

To remove all Docker volumes (this deletes all persistent data):

```bash
docker-compose down -v
```

## **Conclusion**

This comprehensive setup provides a solid foundation for a scalable and maintainable crypto exchange application. By integrating essential technologies like PostgreSQL, Redis, Cassandra, and Kafka, the application is well-equipped to handle high volumes of transactions and real-time data processing. The use of Docker ensures consistent environments, facilitating seamless development and deployment processes.
