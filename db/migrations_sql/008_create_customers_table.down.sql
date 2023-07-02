ALTER TABLE `orders` DROP FOREIGN KEY `FK_order_customer_id`;
ALTER TABLE orders DROP COLUMN customer_id;
DROP TABLE  customers;