package interfaces

type DataSource interface {
	FetchData(req Request) (interface{}, error)
}

type DataDestination interface {
	SendData(data interface{}, req Request) error
}

// Request struct to hold migration request data
type Request struct {
	Input                   string `json:"input"`          // List of input types (Kafka, SQL, MongoDB, etc.)
	Output                  string `json:"output"`         // List of output types (CSV, MongoDB, etc.)
	ConsumerURL             string `json:"consumer_url"`   // URL for Kafka
	ConsumerTopic           string `json:"consumer_topic"` // Topic for Kafka
	ProducerURL             string `json:"producer_url"`
	ProducerTopic           string `json:"producer_topic"`
	SQLSourceConnString     string `json:"sql_source_conn_string"`     // Source SQL connection string
	SQLTargetConnString     string `json:"sql_target_conn_string"`     // Target SQL connection string
	SourceMongoDBConnString string `json:"source_mongodb_conn_string"` // MongoDB source connection string
	SourceMongoDBDatabase   string `json:"source_mongodb_database"`    // MongoDB source database
	SourceMongoDBCollection string `json:"source_mongodb_collection"`  // MongoDB source collection
	TargetMongoDBConnString string `json:"target_mongodb_conn_string"` // MongoDB target connection string
	TargetMongoDBDatabase   string `json:"target_mongodb_database"`    // MongoDB target database
	TargetMongoDBCollection string `json:"target_mongodb_collection"`  // MongoDB target collection
	OutputFileName          string `json:"output_file_name"`           // Output file name for CSVs or other formats
	// RabbitMQ
	RabbitMQInputURL        string `json:"rabbitmq_input_url"`         // URL for RabbitMQ (consumer)
	RabbitMQInputQueueName  string `json:"rabbitmq_input_queue_name"`  // Queue name for RabbitMQ input
	RabbitMQOutputURL       string `json:"rabbitmq_output_url"`        // URL for RabbitMQ (producer)
	RabbitMQOutputQueueName string `json:"rabbitmq_output_queue_name"` // Queue name for RabbitMQ output
	// JSON
	JSONSourceData     string `json:"json_source_data"`     // JSON source data (raw or file path)
	JSONOutputFilename string `json:"json_output_filename"` // JSON output data (raw or file path)
	// YAML
	YAMLSourceFilePath      string `json:"yaml_source_file_path"`      // Source YAML file path
	YAMLDestinationFilePath string `json:"yaml_destination_file_path"` // Destination YAML file path
	// CSV
	CSVSourceFileName      string `json:"csv_source_file_name"`      // Source CSV file name
	CSVDestinationFileName string `json:"csv_destination_file_name"` // Destination CSV file name
	// Dynamodb
	DynamoDBSourceTable  string `json:"dynamodb_source_table"`  // Source DynamoDB table
	DynamoDBTargetTable  string `json:"dynamodb_target_table"`  // Target DynamoDB table
	DynamoDBSourceRegion string `json:"dynamodb_source_region"` // DynamoDB source region
	DynamoDBTargetRegion string `json:"dynamodb_target_region"` // DynamoDB target region
	// FTP
	FTPFILEPATH        string `json:"ftp_file_path"`        // FTP file path
	FTPURL             string `json:"ftp_url"`              // FTP URL
	FTPUser            string `json:"ftp_user"`             // FTP user
	FTPPassword        string `json:"ftp_password"`         // FTP password
	SFTPFILEPATH       string `json:"sftp_file_path"`       // SFTP file path
	SFTPURL            string `json:"sftp_url"`             // SFTP URL
	SFTPUser           string `json:"sftp_user"`            // SFTP user
	SFTPPassword       string `json:"sftp_password"`        // SFTP password
	WebSocketSourceURL string `json:"websocket_source_url"` // WebSocket source URL
	WebSocketDestURL   string `json:"websocket_dest_url"`   // WebSocket destination URL
	// Firebase
	CredentialFileAddr string `json:"firebase_credential_file"`
	Collection         string `json:"firebase_collection"`
	Document           string `json:"firebase_document"`
}
