
#!/bin/bash

# Read the config.yml file
CONFIG_FILE="./../config.yml"

docker cp "$CONFIG_FILE" "spark-master:/opt/spark/conf/config.yml"