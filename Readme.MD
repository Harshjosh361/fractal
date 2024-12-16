# Fractal

**Fractal** is a flexible, configurable data processing tool built with **GoFr** and **Golang**. Fractal is designed to handle data ingestion from multiple sources, apply powerful transformations and validations, and deliver output to a wide range of destinations. With Fractal, you can automate complex data workflows without needing to manage low-level details.
Here's the documentation for setting up a new integration in your project:

### **Custom Syntax Documentation for Validation and Transformation Rules**

## **1. Overview**

The custom syntax enables users to:
1. Validate incoming data to ensure it meets predefined conditions.
2. Transform data fields to fit desired formats, structures, or requirements.
3. Define flexible error-handling strategies for data processing pipelines.

Rules can be written for any data source or destination, such as **JSON**, **YAML**, **CSV**, **SQL Databases**, **Message Brokers**, or **Cloud Services**.

---

## **2. Validation Rules**

Validation rules ensure that data meets specific quality and integrity requirements. 

### **Syntax**
```custom
FIELD(<field_name>) <validation_condition>
```

### **Validation Conditions**

| Condition         | Description                                                                                     | Example                                                                                                                                 |
|--------------------|-------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------|
| `TYPE(<data_type>)` | Ensures the field is of a specified type. Data types: `STRING`, `INT`, `FLOAT`, `BOOL`, `DATE`. | `FIELD("age") TYPE(INT)`                                                                                                               |
| `RANGE(<min>, <max>)` | Ensures the field's value is within a specified range.                                         | `FIELD("price") RANGE(0, 1000)`                                                                                                        |
| `MATCHES(<regex>)`  | Validates that the field's value matches a regular expression pattern.                          | `FIELD("email") MATCHES(EMAIL_REGEX)`                                                                                                  |
| `IN(<value_list>)`   | Validates that the field's value is one of the specified values.                                | `FIELD("status") IN ("active", "inactive")`                                                                                           |
| `REQUIRED`          | Ensures the field is present.                                                                  | `FIELD("name") REQUIRED`                                                                                                              |

### **Examples**
1. Validate that the field `age` is an integer and between 18 and 65:
   ```custom
   FIELD("age") TYPE(INT) RANGE(18, 65)
   ```

2. Ensure `email` matches a regex pattern for valid email addresses:
   ```custom
   FIELD("email") MATCHES(EMAIL_REGEX)
   ```

3. Check that `status` is either "active" or "inactive":
   ```custom
   FIELD("status") IN ("active", "inactive")
   ```

4. Make the `id` field mandatory:
   ```custom
   FIELD("id") REQUIRED
   ```

---

## **3. Transformation Rules**

Transformation rules modify or enrich data to meet specific requirements.

### **Syntax**
```custom
<operation>
```

### **Supported Operations**

| Operation                     | Description                                                                                                                                                       | Example                                                                                                                                         |
|--------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------|
| `RENAME(<old_field>, <new_field>)` | Renames a field in the data.                                                                                                                                      | `RENAME("old_field", "new_field")`                                                                                                             |
| `MAP(<field_name>, {<mapping>})`  | Maps values in a field to new values using a key-value pair mapping.                                                                                             | `MAP("status", {"0": "inactive", "1": "active"})`                                                                                              |
| `ADD_FIELD(<field_name>, <value>)`| Adds a new field with a specified value.                                                                                                                         | `ADD_FIELD("timestamp", CURRENT_TIME())`                                                                                                       |
| `IF <condition> THEN <operation>`| Applies a transformation based on a condition.                                                                                                                  | `IF FIELD("age") > 50 THEN ADD_FIELD("senior_discount", TRUE)`                                                                                 |

### **Examples**
1. Rename a field `old_field` to `new_field`:
   ```custom
   RENAME("old_field", "new_field")
   ```

2. Map values in the `status` field:
   ```custom
   MAP("status", {"0": "inactive", "1": "active"})
   ```

3. Add a field `processed_at` with the current timestamp:
   ```custom
   ADD_FIELD("processed_at", CURRENT_TIME())
   ```

4. Conditionally add a senior discount for users above 50:
   ```custom
   IF FIELD("age") > 50 THEN ADD_FIELD("senior_discount", TRUE)
   ```

---

## **4. Error Handling**

Error handling defines how the system reacts when a validation or transformation fails.

### **Syntax**
```custom
ON_ERROR(<action>)
```

### **Supported Actions**

| Action               | Description                                                                                   | Example                    |
|-----------------------|-----------------------------------------------------------------------------------------------|----------------------------|
| `LOG_AND_CONTINUE`     | Logs the error and continues processing the next record.                                      | `ON_ERROR(LOG_AND_CONTINUE)`|
| `STOP`                | Stops the entire pipeline on encountering an error.                                           | `ON_ERROR(STOP)`           |
| `RETRY`              | Attempts to retry processing the failed record.                                               | `ON_ERROR(RETRY)`          |
| `SEND_TO_QUARANTINE`   | Sends the failed record to a quarantine output for further analysis.                          | `ON_ERROR(SEND_TO_QUARANTINE)`|

### **Examples**
1. Log the error and continue processing:
   ```custom
   ON_ERROR(LOG_AND_CONTINUE)
   ```

2. Stop the pipeline if an error occurs:
   ```custom
   ON_ERROR(STOP)
   ```

---

## **5. Integration-Specific Features**

### **Field Naming**

Depending on the data source, fields can be named using:
1. **JSON Paths:** For JSON/YAML data.
   ```custom
   FIELD("$.user.age") TYPE(INT)
   ```
2. **Column Names:** For CSV or SQL tables.
   ```custom
   FIELD("Column1") TYPE(FLOAT)
   ```
3. **Message Keys/Values:** For Kafka or RabbitMQ.
   ```custom
   FIELD("$.message.key") MATCHES(KEY_REGEX)
   ```

---

## **6. Unified YAML Configuration**

The rules can be embedded into the YAML configuration for pipelines:

```yaml
pipeline:
  error-handling:
   strategy: LOG_AND_CONTINUE
inputconfig:
   csvsourcefilename: sample.csv
   inputmethod: CSV
outputconfig:
   csvdestinationfilename: test.csv
   outputmethod: CSV
cronjob:
   repetition_interval: "1h"
monitoring:
   job_status:"pending"
transformations:
   -ADD_FIELD("processed_at", CURRENT_TIME())
validations:
   -FIELD("age") RANGE(30,35)
```

---

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
error-handling:
   strategy: LOG_AND_CONTINUE
inputconfig:
   csvsourcefilename: sample.csv
   inputmethod: CSV
outputconfig:
   csvdestinationfilename: test.csv
   outputmethod: CSV
cronjob:
   repetition_interval: "1h"
monitoring:
   job_status:"pending"
transformations:
   -ADD_FIELD("processed_at", CURRENT_TIME())
validations:
   -FIELD("age") RANGE(30,35)
  
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
