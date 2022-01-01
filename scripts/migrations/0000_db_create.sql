drop database osinexample_db;
create database osinexample_db;
\c osinexample_db;
CREATE USER golang WITH PASSWORD '123password';
GRANT ALL PRIVILEGES ON DATABASE osinexample_db to golang;
ALTER USER golang CREATEDB;
ALTER ROLE golang SUPERUSER;
