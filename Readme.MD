# Fractal

**Fractal** is a flexible, configurable data processing tool built with **GoFr** and **Golang**. Fractal is designed to handle data ingestion from multiple sources, apply powerful transformations and validations, and deliver output to a wide range of destinations. With Fractal, you can automate complex data workflows without needing to manage low-level details.

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
