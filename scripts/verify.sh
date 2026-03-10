#!/usr/bin/env bash

set -euo pipefail

export PYTEST_DISABLE_PLUGIN_AUTOLOAD=1
: "${VERIFY_TIMEOUT_SECONDS:=900}"

run_with_timeout() {
  if command -v timeout >/dev/null 2>&1; then
    timeout "${VERIFY_TIMEOUT_SECONDS}" "$@"
  else
    "$@"
  fi
}

run_python_tests() {
  if [ -f pyproject.toml ] || [ -f pytest.ini ] || [ -f setup.cfg ] || [ -f tox.ini ] || find . -maxdepth 2 -type d \( -name test -o -name tests \) | grep -q .; then
    if command -v pytest >/dev/null 2>&1; then
      echo "Running Python tests with pytest"
      run_with_timeout pytest
      return 0
    fi
    if command -v python >/dev/null 2>&1; then
      echo "Python test configuration detected, but pytest is not installed" >&2
      return 1
    fi
  fi
  return 0
}

run_node_tests() {
  if [ -f package.json ]; then
    if command -v npm >/dev/null 2>&1; then
      echo "Running Node tests with npm test"
      run_with_timeout npm test
      return 0
    fi
    echo "Node test configuration detected, but npm is not installed" >&2
    return 1
  fi
  return 0
}

main() {
  local ran_any=0

  if [ -f pyproject.toml ] || [ -f pytest.ini ] || [ -f setup.cfg ] || [ -f tox.ini ] || find . -maxdepth 2 -type d \( -name test -o -name tests \) | grep -q .; then
    ran_any=1
    run_python_tests
  fi

  if [ -f package.json ]; then
    ran_any=1
    run_node_tests
  fi

  if [ "${ran_any}" -eq 0 ]; then
    echo "No Python or Node test configuration detected"
  fi
}

main "$@"
