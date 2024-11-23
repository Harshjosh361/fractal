from locust import HttpUser, task, between
import json

class MigrationUser(HttpUser):
    # Define the wait time between tasks (e.g., 1-2 seconds between requests)
    wait_time = between(1, 2)

    @task
    def run_migration(self):
        # Define the payload (this matches the API endpoint expected request structure)
        payload = {
            "input": "CSV",
            "output": "CSV",
            "csv_source_file_name": "sample.csv",
            "csv_destination_file_name": "destination_file.csv"
        }

        # Send the POST request to the /api/migration endpoint
        self.client.post("/api/migration", data=json.dumps(payload), headers={"Content-Type": "application/json"})

