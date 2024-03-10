#!/usr/bin/env bash

# Local tests configuration file
# Configure to your Gatling bin directory
GATLING_BIN_DIR=$HOME/gatling/3.10.3/bin

RESULTS_WORKSPACE="$(pwd)/load-test/user-files/results"
GATLING_WORKSPACE="$(pwd)/load-test/user-files"

# Clean possible daemon conflicts and run the containers
runDockerCompose() {
	docker-compose down --volumes &&
		docker-compose up -d
}

# Run the load and validation tests
runGatling() {
	sh $GATLING_BIN_DIR/gatling.sh -rm local -s RinhaBackendCrebitosSimulation \
		-rd "Rinha de Backend - 2024/Q1: Crébito" \
		-rf $RESULTS_WORKSPACE \
		-sf "$GATLING_WORKSPACE/simulations"
}

# Try to connect before start the tests
startTest() {
	runDockerCompose

	# Try to connect for 60s while docker containers are running in background
	for i in {1..60}; do
		# 2 requests to wake the 2 api instances up :)
		curl --fail http://localhost:9999/clientes/1/extrato &&
			echo "" &&
			curl --fail http://localhost:9999/clientes/1/extrato &&
			echo "" &&
			runGatling &&
			break || sleep 1
	done
}

startTest
