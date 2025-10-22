# Configuration Directory

This directory contains all the YAML configuration files for Moogie synthetic checks.

## Structure

- `checks/` - Contains individual check configuration files
- Each YAML file defines a single synthetic check
- Files follow Kubernetes-style resource definitions

## Usage

All YAML files in the `checks/` directory will be automatically loaded by Moogie. The runner, API, and UI components all have access to these configuration files.

See the main README.md for detailed configuration examples and syntax.
