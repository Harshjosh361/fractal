# Fractal

**Fractal** is a flexible, configurable data processing tool built with **GoFr** and **Golang**. Fractal is designed to handle data ingestion from multiple sources, apply powerful transformations and validations, and deliver output to a wide range of destinations. With Fractal, you can automate complex data workflows without needing to manage low-level details.
Here's the documentation for setting up a new integration in your project:

## To Run RabbitMQ :  ```bash docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management```

# Adding a New Integration

The system is designed to make it simple to add new data integrations for both input and output. Each integration should define methods to read (input) and write (output) data, following a unified interface approach.

## Steps to Add a New Integration

1. **Create a New File**:  
   In the `integrations` directory, create a new file named after the integration. For example, if adding support for RabbitMQ, create `rabbitmq.go`.

2. **Implement the `Source` and `Destination` Interfaces**:  
   In this file, you need to define structs that implement the `Source` and `Destination` interfaces.

   - **Source Interface**: This interface should define the `Read` method, which reads data from the input source and returns it in a standardized format.
   - **Destination Interface**: This interface should define the `Write` method, which writes data to the output destination.

   Here's an example structure for `rabbitmq.go`:

   ```go
   package integrations

   import "fmt"

   // RabbitMQInput struct to handle input from RabbitMQ
   type RabbitMQInput struct {
       // Add any necessary fields, such as connection settings, queues, etc.
   }

   // Read method to read from RabbitMQ
   func (r *RabbitMQInput) Read() (interface{}, error) {
       // Implement logic to read from RabbitMQ
       fmt.Println("Reading from RabbitMQ...")
       return nil, nil
   }

   // RabbitMQOutput struct to handle output to RabbitMQ
   type RabbitMQOutput struct {
       // Add any necessary fields, such as connection settings, queues, etc.
   }

   // Write method to write to RabbitMQ
   func (r *RabbitMQOutput) Write(data interface{}) error {
       // Implement logic to write to RabbitMQ
       fmt.Println("Writing to RabbitMQ...")
       return nil
   }

   // Initialize the new integration
   func init() {
       RegisterSource("rabbitmq", &RabbitMQInput{})
       RegisterDestination("rabbitmq", &RabbitMQOutput{})
   }
   ```

3. **Register the Integration**:  
   In the `init()` function, use `RegisterSource` and `RegisterDestination` to add the integration to the system. This makes it available for both CLI and HTTP server modes.

4. **Configuration**:  
   If the integration requires additional configuration (like credentials or connection strings), make sure to add relevant fields to the struct and include a way to parse this information from the user-provided configuration.

5. **Testing the Integration**:  
   Run the application and select the new integration in either CLI or HTTP mode. Verify that data can be read from and written to the integration correctly.


With this setup, adding integrations is straightforward. Each integration can now be quickly defined and registered, keeping your system scalable and modular.



## Features
- **Multi-Source Data Ingestion**: Supports data ingestion from HTTP, CSV files, SQL databases, Pub-Sub systems, cloud storage, and more.
- **Customizable Data Transformations**: Apply data transformations, including data mapping, filtering, aggregation, and enrichment, with built-in or custom functions.
- **Validation Rules**: Define validation schemas to ensure incoming data meets quality standards before processing.
- **Flexible Output Options**: Output processed data to databases (SQL/NoSQL), CSV files, messaging queues, HTTP responses, or cloud storage.
- **YAML Configuration**: Configure data workflows and transformation rules through a YAML file for easy setup and customization.

## Getting Started

### Prerequisites
- Go 1.18+
- [GoFr Framework](https://gofrframework.com) installed

### Installation
Clone the repository and navigate to the Fractal directory:

```bash
git clone https://github.com/SkySingh04/fractal.git
cd fractal
```

Install the dependencies:

```bash
go mod tidy
```

### Configuration
Set up a `.yaml` configuration file in the root directory. Define inputs, transformations, validations, and outputs as per your workflow needs. Here's a basic example:

```yaml
pipeline:
  - name: "CSV to MongoDB Migration"
    input:
      type: "CSV"
      source: "./data/input.csv"
    validation:
      schema:
        fields:
          - name: "id"
            type: "integer"
            required: true
          - name: "name"
            type: "string"
            required: true
    transformation:
      - map:
          from: "old_field_name"
          to: "new_field_name"
    output:
      type: "MongoDB"
      destination: "mongodb://localhost:27017/fractal_db"
      collection: "data_output"
```

### Running Fractal
Start the pipeline using:

```bash
go run main.go -config=config.yaml
```

### Example Use Cases
- **Data Migration**: Migrate data from legacy systems to cloud databases or NoSQL databases.
- **Log Aggregation**: Aggregate logs from multiple sources and send them to a searchable data store.
- **Content Syndication**: Ingest and format content from RSS feeds or APIs, and distribute it across platforms.
- **Data Quality Checker**: Validate incoming data streams to ensure data quality before storing.

## Contributing
Contributions are welcome! Feel free to submit pull requests for new features, bug fixes, or documentation improvements.

## License
This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
