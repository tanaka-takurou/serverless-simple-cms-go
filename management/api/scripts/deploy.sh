#!/bin/bash
echo 'Creating function.zip...'
`dirname $0`/create_function.sh

echo 'Updating Lambda-Function...'
cd `dirname $0`/../
aws lambda update-function-code \
	--profile default \
	--function-name SampleCMSManagementApi \
	--zip-file fileb://`pwd`/function.zip \
	--cli-connect-timeout 6000 \
	--publish
