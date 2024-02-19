create table products(
    `id` int not null primary key auto_increment,
    `description` text not null,
    expiration_rate float not null,
    freezing_rate float not null,
    height float not null,
    lenght float not null,
    netweight float not null,
    product_code text not null,
    recommended_freezing_temperature float not null,
    width float not null,
    id_product_type int not null,
    id_seller int not null
);
create table employees(
    `id` int not null primary key auto_increment,
    card_number_id text not null,
    first_name text not null,
    last_name text not null,
    warehouse_id int not null
);
create table warehouses(
    `id` int not null primary key auto_increment,
    `address` text null,
    telephone text null,
    warehouse_code text null,
    minimum_capacity int null,
    minimum_temperature int null
);
create table sections(
    `id` int not null primary key auto_increment,
    section_number int not null,
    current_temperature int not null,
    minimum_temperature int not null,
    current_capacity int not null,
    minimum_capacity int not null,
    maximum_capacity int not null,
    warehouse_id int not null,
    id_product_type int not null
);
create table sellers(
    `id` int not null primary key auto_increment,
    cid int not null,
    company_name text not null,
    `address` text not null,
    telephone varchar(15) not null
);
create table buyers(
    `id` int not null primary key auto_increment,
    card_number_id text not null,
    first_name text not null,
    last_name text not null
);

create table inboudOrders(
    `id` int not null primary key auto_increment,
    order_date date not null,
    order_number text not null,
    employee_id int not null,
    product_batch_id int not null,
    warehouse_id int not null,
    CONSTRAINT fk_employee_id FOREIGN KEY (employee_id) REFERENCES employees(id),
    -- CONSTRAINT fk_product_batch_id FOREIGN KEY (product_batch_id) REFERENCES productBatches(id),
    CONSTRAINT fk_warehouse_id FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
);

create table productBatches(
    `id` int not null primary key auto_increment,
    batch_number int not null,
    current_quantity int not null,
    current_temperature int not null,
    due_date date not null,
    initial_quantity int not null,
    manufacturing_date date not null,
    manufacturing_hour int not null,
    minimum_temperature int not null,
    product_id int not null,
    section_id int not null,
    CONSTRAINT fk_product_id FOREIGN KEY (product_id) REFERENCES products(id),
    CONSTRAINT fk_section_id FOREIGN KEY (section_id) REFERENCES sections(id)
);

create table purchase_orders(
    `id` int not null primary key auto_increment,
    order_number varchar(20) not null,
    order_date date default null,
    tracking_code varchar(8) not null,
    buyer_id int not null,
    product_record_id int not null,
    order_status_id int not null,
    foreign key (buyer_id) references buyers(id),
    foreign key (product_record_id) references productsRecord(id)
)
