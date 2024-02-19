USE `mysqlapigo`;

-- DML
-- products data
INSERT INTO `products` (`id`, `description`, `expiration_rate`, `freezing_rate`, `height`, `lenght`, `netweight`, `product_code`, `recommended_freezing_temperature`, `width`, `id_product_type`, `id_seller`) VALUES
(1, 'Fresh Milk', 0.1, 0.05, 25.0, 10.0, 1.0, 'MILK1001', -4.0, 10.0, 1, 101),
(2, 'Frozen Peas', 0.02, 0.2, 5.0, 15.0, 0.5, 'PEAS2002', -18.0, 8.0, 2, 102),
(3, 'Whole Wheat Bread', 0.15, 0, 10.0, 20.0, 0.75, 'BRED3003', 0.0, 15.0, 3, 103),
(4, 'Canned Tuna', 0.01, 0, 8.0, 8.0, 0.3, 'TUNA4004', 0.0, 8.0, 4, 104),
(5, 'Apple Juice', 0.07, 0.03, 30.0, 10.0, 2.0, 'AJUC5005', -5.0, 10.0, 5, 105),
(6, 'Frozen Pizza', 0.04, 0.25, 2.0, 30.0, 0.8, 'PIZZ6006', -15.0, 30.0, 6, 106),
(7, 'Fresh Strawberries', 0.2, 0.1, 10.0, 10.0, 0.25, 'STRW7007', -3.0, 10.0, 7, 107),
(8, 'Organic Eggs', 0.12, 0, 15.0, 12.0, 0.6, 'EGGS8008', 0.0, 12.0, 8, 108),
(9, 'Dark Chocolate', 0.05, 0, 2.0, 15.0, 0.1, 'CHOC9009', 0.0, 8.0, 9, 109),
(10, 'Almond Milk', 0.08, 0.04, 25.0, 10.0, 1.0, 'AMLK0010', -4.0, 10.0, 10, 110);

-- employees data
INSERT INTO `employees` (`id`, `card_number_id`, `first_name`, `last_name`, `warehouse_id`) VALUES
(1, 'A123B456C', 'John', 'Doe', 1),
(2, 'D789E012F', 'Jane', 'Smith', 2),
(3, 'G345H678I', 'Michael', 'Johnson', 1),
(4, 'J901K234L', 'Emily', 'Davis', 3),
(5, 'M567N890O', 'David', 'Wilson', 2),
(6, 'P123Q456R', 'Sarah', 'Brown', 1),
(7, 'S789T012U', 'James', 'Miller', 3),
(8, 'V345W678X', 'Jessica', 'Taylor', 2),
(9, 'Y901Z234A', 'William', 'Anderson', 1),
(10, 'B567C890D', 'Emma', 'Thomas', 3);

-- warehouses data
INSERT INTO `warehouses` (`id`, `address`, `telephone`, `warehouse_code`, `minimum_capacity`, `minimum_temperature`) VALUES
(1, '123 Main St, Cityville', '123-456-7890', 'WH001', 5000, -20),
(2, '456 Elm St, Townsburg', '234-567-8901', 'WH002', 10000, -15),
(3, '789 Oak St, Villageford', '345-678-9012', 'WH003', 3000, 0),
(4, '101 Pine St, Hamletville', '456-789-0123', 'WH004', 8000, -10),
(5, '202 Maple St, Boroughland', '567-890-1234', 'WH005', 7500, -5),
(6, '303 Birch St, Districtcity', '678-901-2345', 'WH006', 6000, -18),
(7, '404 Cedar St, Countytown', '789-012-3456', 'WH007', 4500, -20),
(8, '505 Cherry St, Regionburgh', '890-123-4567', 'WH008', 8500, -22),
(9, '606 Ash St, Provincefield', '901-234-5678', 'WH009', 5000, -10),
(10, '707 Willow St, Sectorville', '012-345-6789', 'WH010', 9500, -12);

-- sections data
INSERT INTO `sections` (`id`, `section_number`, `current_temperature`, `minimum_temperature`, `current_capacity`, `minimum_capacity`, `maximum_capacity`, `warehouse_id`, `id_product_type`) VALUES
(1, 1, -18, -20, 2000, 1500, 2500, 1, 1),
(2, 2, -15, -15, 3000, 2000, 4000, 1, 1),
(3, 3, 5, 0, 1000, 500, 1500, 2, 1),
(4, 4, -10, -10, 2500, 2000, 3000, 2, 2),
(5, 5, -5, -5, 3500, 3000, 4500, 3, 2),
(6, 6, -18, -20, 2200, 1800, 2700, 3, 2),
(7, 7, -20, -20, 1500, 1000, 2000, 4, 3),
(8, 8, -22, -22, 3200, 2500, 4000, 4, 3),
(9, 9, -10, -10, 2800, 2300, 3300, 5, 3),
(10, 10, -12, -15, 2300, 1800, 2900, 5, 4);

-- sellers data
INSERT INTO `sellers` (`id`, `cid`, `company_name`, `address`, `telephone`) VALUES
(1, 10001, 'FreshFoods Ltd.', '123 Green St, Veggieville', '555-0101'),
(2, 10002, 'DairyDelight Inc.', '456 Milk Rd, DairyTown', '555-0202'),
(3, 10003, 'GrainGuru Co.', '789 Wheat Way, Grainfield', '555-0303'),
(4, 10004, 'FruitFantasy LLC', '101 Apple Ave, Orchard City', '555-0404'),
(5, 10005, 'MeatMasters Ltd.', '202 Steak St, Carnivoreville', '555-0505'),
(6, 10006, 'SeafoodSolutions Inc.', '303 Fisherman Wharf, Oceanview', '555-0606'),
(7, 10007, 'BakeryBliss Co.', '404 Dough Ln, Breadtown', '555-0707'),
(8, 10008, 'FrozenFeasts LLC', '505 Chill St, Frozenville', '555-0808'),
(9, 10009, 'BeverageBarn Ltd.', '606 Drink Dr, Liquid City', '555-0909'),
(10, 10010, 'SnackSquad Inc.', '707 Crunch Cr, Snackburg', '555-1010');

-- buyers data
INSERT INTO `buyers` (`id`, `card_number_id`, `first_name`, `last_name`) VALUES
(1, 'BC12345678', 'Alice', 'Johnson'),
(2, 'BC23456789', 'Bob', 'Smith'),
(3, 'BC34567890', 'Carol', 'Davis'),
(4, 'BC45678901', 'David', 'Wilson'),
(5, 'BC56789012', 'Eve', 'Martinez'),
(6, 'BC67890123', 'Frank', 'Garcia'),
(7, 'BC78901234', 'Grace', 'Brown'),
(8, 'BC89012345', 'Henry', 'Miller'),
(9, 'BC90123456', 'Ivy', 'Lee'),
(10, 'BC01234567', 'Jack', 'Martin');

INSERT INTO `carries` (`id`, `cid`, `company_name`, `address`, `telephone`, `locality_id`) VALUES
(1, "C01", 'FreshFoods Ltd.', '123 Green St, Veggieville', '555-0101', 1),
(2, "C02", 'DairyDelight Inc.', '456 Milk Rd, DairyTown', '555-0202', 2),
(3, "C03", 'GrainGuru Co.', '789 Wheat Way, Grainfield', '555-0303', 3),
(4, "C04", 'FruitFantasy LLC', '101 Apple Ave, Orchard City', '555-0404', 4),
(5, "C05", 'MeatMasters Ltd.', '202 Steak St, Carnivoreville', '555-0505', 5),
(6, "C06", 'SeafoodSolutions Inc.', '303 Fisherman Wharf, Oceanview', '555-0606', 6),
(7, "C07", 'BakeryBliss Co.', '404 Dough Ln, Breadtown', '555-0707', 7),
(8, "C08", 'FrozenFeasts LLC', '505 Chill St, Frozenville', '555-0808', 8),
(9, "C09", 'BeverageBarn Ltd.', '606 Drink Dr, Liquid City', '555-0909', 9),
(10, "C10", 'SnackSquad Inc.', '707 Crunch Cr, Snackburg', '555-1010', 10);

INSERT INTO `locality` (`id`, `postal_code`, `locality_name`, `province_name`, `country_name`) 
VALUES
(1, 1, 'Veggieville', 'Buenos Aires', 'Argentina'),
(2, 2, 'DairyTown', 'Buenos Aires', 'Argentina'),
(3, 3, 'Grainfield', 'Buenos Aires', 'Argentina'),
(4, 4, 'Orchard City', 'Buenos Aires', 'Argentina'),
(5, 5, 'Carnivoreville', 'Buenos Aires', 'Argentina'),
(6, 6, 'Oceanview', 'Buenos Aires', 'Argentina'),
(7, 7, 'Breadtown', 'Buenos Aires', 'Argentina'),
(8, 8, 'Frozenville', 'Buenos Aires', 'Argentina'),
(9, 9, 'Liquid City', 'Buenos Aires', 'Argentina'),
(10, 10, 'Snackburg', 'Buenos Aires', 'Argentina');


-- inboudOrders data
INSERT INTO `inboudOrders` (`id`, `order_date`, `order_number`, `employee_id`, `product_batch_id`, `warehouse_id`) VALUES
(1, '2024-11-10', 'order#1', '10', '2', '1'),
(2, '2024-12-10', 'order#2', '1', '1', '2'),
(3, '2024-09-10', 'order#3', '4', '3', '2'),
(4, '2024-07-07', 'order#4', '4', '4', '3'),
(5, '2024-07-19', 'order#5', '1', '5', '1'),
(6, '2024-05-10', 'order#6', '2', '5', '2'),
(7, '2024-07-27', 'order#7', '11', '7', '3'),
(8, '2024-01-06', 'order#8', '1', '9', '3'),
(9, '2024-12-14', 'order#9', '5', '10', '2'),
(10, '2024-04-25', 'order#10', '2', '2', '1');

-- productBatch data
INSERT INTO `productBatches` (`id`, `batch_number`, `current_quantity`, `current_temperature`, `due_date`, `initial_quantity`, `manufacturing_date`, `manufacturing_hour`, `minimum_temperature`, `product_id`, `section_id`) VALUES
(1, 1, 5, -10, '2024-11-10', 10, '2024-01-04', 10, -20, 1, 1),
(2, 2, 7, -3, '2024-09-06', 15, '2024-01-01', 9, -10, 2, 2),
(3, 3, 6, -15, '2024-12-30', 15, '2024-01-03', 11, -15, 3, 3),
(4, 4, 23, -1, '2024-04-04', 20, '2023-12-29', 5, -25, 4, 4),
(5, 5, 14, 6, '2024-05-17', 12, '2023-11-25', 3, -15, 5, 5),
(6, 6, 7, 8, '2024-11-01', 24, '2023-12-01', 7, -10, 6, 6),
(7, 7, 11, 3, '2024-11-30', 15, '2023-11-30', 4, -5, 7, 7),
(8, 8, 10, -13, '2024-07-10', 10, '2024-01-02', 2, -5, 8, 8),
(9, 9, 6, -20, '2024-08-04', 12, '2024-01-03', 1, 0, 9, 9),
(10, 10, 1, -16, '2024-07-26', 12, '2023-11-28', 11, -20, 10, 10);

-- purchase orders data
INSERT INTO `purchase_orders` (`id`, `order_number`, `order_date`, `tracking_code`, `buyer_id`, `product_record_id`, `order_status_id`) VALUES
(1, 'order#1', STR_TO_DATE('2023-01-01', '%Y-%m-%d'), 'abscf121', 1, 1, 1),
(2, 'order#2', STR_TO_DATE('2023-02-01', '%Y-%m-%d'), 'abscf122', 2, 2, 1),
(3, 'order#3', STR_TO_DATE('2023-03-01', '%Y-%m-%d'), 'abscf123', 1, 3, 2),
(4, 'order#4', STR_TO_DATE('2023-04-01', '%Y-%m-%d'), 'abscf124', 4, 1, 2),
(5, 'order#5', STR_TO_DATE('2023-05-01', '%Y-%m-%d'), 'abscf125', 5, 2, 3),
(6, 'order#6', STR_TO_DATE('2023-06-01', '%Y-%m-%d'), 'abscf126', 6, 3, 3),
(7, 'order#7', STR_TO_DATE('2023-07-01', '%Y-%m-%d'), 'abscf127', 6, 1, 1),
(8, 'order#8', STR_TO_DATE('2023-08-01', '%Y-%m-%d'), 'abscf128', 6, 2, 2),
(9, 'order#9', STR_TO_DATE('2023-09-01', '%Y-%m-%d'), 'abscf129', 5, 3, 2),
(10,'order#10', STR_TO_DATE('2023-10-01', '%Y-%m-%d'), 'abscf120', 7, 1, 3);

-- productsRecord data
INSERT INTO productsRecord (last_update_date, purchase_price, sale_price, product_id)
VALUES
('2021-01-10', 100, 110, 1),
('2022-02-20', 200, 230, 2),
('2023-03-05', 300, 340, 3),
('2024-01-11', 400, 500, 3);