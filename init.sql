CREATE TABLE orders(id varchar(37) PRIMARY KEY, title varchar(50),description varchar(255),dealine datetime);

-- CREATE TABLE requirements(requirementid integer PRIMARY KEY AUTO_INCREMENT,request varchar(50),expectedoutcome varchar(50),order_id integer,status bool,FOREIGN KEY(order_id) REFERENCES orders(id));
CREATE TABLE requirements(requirementid SERIAL PRIMARY KEY,request varchar(50),expectedoutcome varchar(50),order_id varchar(37),status bool,FOREIGN KEY(order_id) REFERENCES orders(id));