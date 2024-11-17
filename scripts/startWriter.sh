#!/bin/bash

# Obtenir le chemin absolu du script
CURRENT_PATH=$(dirname "$(realpath "$0")")

$CURRENT_PATH/load-test-rds --hostname localhost --port 5432 --database load_test --username load_test --password pg_password