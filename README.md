# MySQL Lambda

A small MySQL UDF library for making AWS Lambda calls from a MySQL database. This adds `lambda_sync` and `lambda_async` functions to a regular MySQL database. The purpose of this is to help with migration to or from a managed AWS DB that has these functions natively, allowing for custom code executing directly from MySQL.

You can read more about these functions here: https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/AuroraMySQL.Integrating.Lambda.html.

## Usage

### `lambda_sync`

You invoke the lambda_sync function synchronously with the RequestResponse invocation type. The function returns the result of the Lambda invocation in a JSON payload. The function has the following syntax.

```sql
`lambda_sync` ( @lambda_function_ARN` , @JSON_payload` )
```

#### `@lambda_function_ARN` (string)

The Amazon Resource Name (ARN) of the Lambda function to invoke.

#### `@JSON_payload` (string)

The payload for the invoked Lambda function, in JSON format.

---

### `lambda_async`

You invoke the lambda_async function asynchronously with the Event invocation type. The function returns null immediately and the function is executed in the background. The function has the following syntax.

```sql
`lambda_async` ( @lambda_function_ARN , @JSON_payload )
```

#### `@lambda_function_ARN` (string)

The Amazon Resource Name (ARN) of the Lambda function to invoke.

#### `@JSON_payload` (string)

The payload for the invoked Lambda function, in JSON format.

---

## Docker

There is a Dockerfile included that when will generate `mysql_lambda.so` for you with Ubuntu 20.04. Compiling on Ubuntu on 20.10 or higher causes the extension not to work in older versions of Linux due to incompatible versions of glibc.

```shell
git clone https://github.com/StirlingMarketingGroup/mysql-lambda.git
cd mysql-lambda
docker build -o . .
sudo cp mysql_lambda.so /usr/lib/mysql/plugin/mysql_lambda.so # replace plugin dir here if needed
```

Continue below to see the MySQL commands needed to make the functions work in MySQL.

## Dependencies

You will need Golang, which you can get from here https://golang.org/doc/install.

You will also need to install the MySQL dev library:

### Debian / Ubuntu

```shell
sudo apt update
sudo apt install libmysqlclient-dev
```

## Installing

Know your MySQL plugin directory, which can be found by running this MySQL query:

```sql
select @@plugin_dir;
```

then replace `/usr/lib/mysql/plugin` below with your MySQL plugin directory.

```shell
cd ~ # or wherever you store your git projects
git clone https://github.com/StirlingMarketingGroup/mysql-lambda.git
cd mysql-lambda
go build -buildmode=c-shared -o mysql_lambda.so
sudo cp mysql_lambda.so /usr/lib/mysql/plugin/mysql_lambda.so # replace plugin dir here if needed
```

Enable the function in MySQL by running these MySQL commands

```sql
create function`_lambda_sync`returns string soname'mysql_lambda.so';
create function`lambda_async`returns string soname'mysql_lambda.so';

DROP function IF EXISTS `lambda_sync`;

DELIMITER $$
USE `sterling`$$
CREATE FUNCTION `lambda_sync` (`$arn` varchar (2048), `$payload` json)
RETURNS json
BEGIN
RETURN cast(`_lambda_sync`(`$arn`,`$payload`)as json);
END$$

DELIMITER ;
```