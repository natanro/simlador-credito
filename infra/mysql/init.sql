-- db/init.sql
CREATE DATABASE IF NOT EXISTS simlador_credito;
USE simlador_credito;

CREATE TABLE IF NOT EXISTS simulacoes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    requested_amount DECIMAL(10, 2) NOT NULL,
    installments INT NOT NULL,
    status VARCHAR(20) NOT NULL,
    age INT NOT NULL,
    annual_rate DECIMAL(10, 4) NOT NULL,
    monthly_rate DECIMAL(10, 4) NOT NULL,
    monthly_payment DECIMAL(10, 2) NOT NULL,
    total_amount DECIMAL(10, 2)
);