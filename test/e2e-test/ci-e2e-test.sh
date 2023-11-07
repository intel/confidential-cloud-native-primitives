#!/bin/bash
set -o errexit
: '
This is an E2E test script.
If you want to run the script , please run the ci-setup.sh script in advance to set up the ci test environment.
'

WORK_DIR=$(cd "$(dirname "$0")"; pwd)

#Run E2E test cases
pip install pytest pytdxattest
pytest  "${WORK_DIR}/test_eventlog.py" "${WORK_DIR}/test_tdquote.py"  "${WORK_DIR}/test_tdreport.py"

