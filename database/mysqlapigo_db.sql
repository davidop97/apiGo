-- DDL
DROP DATABASE IF EXISTS `mysqlapigo`;

CREATE DATABASE `mysqlapigo`;

USE `mysqlapigo`;

-- table `products`
CREATE TABLE `products` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `description` text NOT NULL,
    `expiration_rate` float NOT NULL,
    `freezing_rate` float NOT NULL,
    `height` float NOT NULL,
    `lenght` float NOT NULL,
    `netweight` float NOT NULL,
    `product_code` text NOT NULL,
    `recommended_freezing_temperature` float NOT NULL,
    `width` float NOT NULL,
    `id_product_type` int(11) NOT NULL,
    `id_seller` int(11) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- table `employees`
CREATE TABLE `employees` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `card_number_id` varchar(255) NOT NULL,
    `first_name` varchar(255) NOT NULL,
    `last_name` varchar(255) NOT NULL,
    `warehouse_id` int(11) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- table `warehouses`
CREATE TABLE `warehouses` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `address` varchar(255) NOT NULL,
    `telephone` varchar(255) NOT NULL,
    `warehouse_code` varchar(255) NOT NULL,
    `minimum_capacity` int(11) NOT NULL,
    `minimum_temperature` int(11) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- table `sections`
CREATE TABLE `sections` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `section_number` int(11) NOT NULL,
    `current_temperature` int(11) NOT NULL,
    `minimum_temperature` int(11) NOT NULL,
    `current_capacity` int(11) NOT NULL,
    `minimum_capacity` int(11) NOT NULL,
    `maximum_capacity` int(11) NOT NULL,
    `warehouse_id` int(11) NOT NULL,
    `id_product_type` int(11) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- table `sellers`
CREATE TABLE `sellers` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `cid` int(11) NOT NULL,
    `company_name` varchar(255) NOT NULL,
    `address` varchar(255) NOT NULL,
    `telephone` varchar(15) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- table `buyers`
CREATE TABLE `buyers` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `card_number_id` varchar(255) NOT NULL,
    `first_name` varchar(255) NOT NULL,
    `last_name` varchar(255) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- table `productBatches`
CREATE TABLE `productBatches` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `batch_number` int NOT NULL,
    `current_quantity` int NOT NULL,
    `current_temperature` int NOT NULL,
    `due_date` date NOT NULL,
    `initial_quantity` int NOT NULL,
    `manufacturing_date` date NOT NULL,
    `manufacturing_hour` int NOT NULL,
    `minimum_temperature` int NOT NULL,
    `product_id` int NOT NULL,
    `section_id` int NOT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT fk_product_id FOREIGN KEY (product_id) REFERENCES products(id),
    CONSTRAINT fk_section_id FOREIGN KEY (section_id) REFERENCES sections(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- table `inboudOrders`
CREATE TABLE `inboudOrders` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `order_date` date NOT NULL,
    `order_number` varchar(255) NOT NULL,
    `employee_id` int(11) NOT NULL,
    `product_batch_id` int(11) NOT NULL,
    `warehouse_id` int(11) NOT NULL,
    PRIMARY KEY (`id`),
    -- CONSTRAINT fk_employee_id FOREIGN KEY (employee_id) REFERENCES employees(id),
    CONSTRAINT fk_warehouse_id FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- table `productsRecord`
CREATE TABLE `productsRecord` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `last_update_date` date NOT NULL,
    `purchase_price` float NOT NULL,
    `sale_price` float NOT NULL,
    `product_id` int(11) NOT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT fk_product_id_productRecord FOREIGN KEY (`product_id`) REFERENCES `products`(`id`) ON DELETE CASCADE
);

-- table `purchase_orders`
CREATE TABLE purchase_orders(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `order_number` varchar(20) NOT NULL,
    `order_date` date default null,
    `tracking_code` varchar(8) NOT NULL,
    `buyer_id` int NOT NULL,
    `product_record_id` int NOT NULL,
    `order_status_id` int NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`buyer_id`) REFERENCES buyers(`id`)
    -- FOREIGN KEY (`product_record_id`) REFERENCES productsRecord(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

create table carries(
    `id` int not null primary key auto_increment,
    cid text not null,
    company_name text not null,
    `address` text not null,
    telephone varchar(15) not null,
    locality_id int not null
);

-- new table (Added in Sprint2)
create table locality(
    `id` int not null primary key auto_increment,
    postal_code int not null,
    locality_name text not null,
    province_name text not null,
    country_name text not null
);

-- Update table sellers: add column locality_id (FK)
ALTER TABLE `sellers` ADD locality_id int not null  DEFAULT 0;





