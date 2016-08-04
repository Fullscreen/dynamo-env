## dynamo-env

Access DynamoDB key value pairs through a `/usr/bin/env` interface.

usage
=====
```shell
dynamo-env [-i] [--table dynamo_table] [name=value ...] [utility [argument ...]]
```

examples
========
```shell
# fetch all environment key pairs
dynamo-env -t table

# run a command with all environment key pairs
dynamo-env -t table command

# run a command with overloaded key pairs
dynamo-env -t table key=value command

# run a command with key pairs, ignoring the inherited environment
dynamo-env -i -t table command
```

dynamo
======

Your dynamo table needs to have a primary partition key named "Name" for this
tool to work properly. You can create a test table with the following command:

```shell
aws dynamodb create-table \
	--table-name mytable \
	--attribute-definitions AttributeName=Name,AttributeType=S \
	--key-schema AttributeName=Name,KeyType=HASH \
	--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5
```

install
=======
```shell
go get github.com/fullscreen/dynamo-env
```

Make sure your `PATH` includes your `$GOPATH` bin directory:

```shell
export PATH=$PATH:$GOPATH/bin
```
