CREATE TABLE orders(id varchar(37) PRIMARY KEY, title varchar(50),description varchar(255),deadline timestamp );

-- CREATE TABLE requirements(requirementid integer PRIMARY KEY AUTO_INCREMENT,request varchar(50),expectedoutcome varchar(50),order_id integer,status bool,FOREIGN KEY(order_id) REFERENCES orders(id));
CREATE TABLE requirements(id SERIAL PRIMARY KEY,request varchar(50),expectedoutcome varchar(50),order_id varchar(37),status bool,FOREIGN KEY(order_id) REFERENCES orders(id));

CREATE TABLE users(id varchar(37) PRIMARY KEY, username varchar(50),email varchar(32), pswd varchar(32));

INSERT INTO orders (id, title, description, deadline) values('asdasd','asd','asdasfasdasd','2 jan 2022');
INSERT INTO requirements (request, expectedoutcome, order_id, status) values('asdasd','asdasd','asdasd','0');

INSERT INTO users (id, username, email, pswd) values('123','elloy',md5('elloy@gmail.com'),md5('100300'));