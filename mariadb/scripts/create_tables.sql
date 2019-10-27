use mysql;


DROP TABLE IF EXISTS `clientes`;
CREATE TABLE `clientes` (
  `id` int(6) unsigned NOT NULL AUTO_INCREMENT,
  `first_name` varchar(30) NOT NULL,
  `last_name` varchar(30) NOT NULL,
  `birthday` varchar(10) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

DROP TABLE IF EXISTS `campanhas`;
CREATE TABLE `campanhas` (
  `id` int(6) unsigned NOT NULL AUTO_INCREMENT,
  `campanha` varchar(30) NOT NULL,
  `status` boolean NOT NULL,
  `pct` int(2) NOT NULL,  
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

Select * from campanhas;

DROP TABLE IF EXISTS `produtos`;
CREATE TABLE `produtos` (
  `id` int(6) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(30) NOT NULL,
  `description` varchar(30) NOT NULL,
  `priceincents` int(10) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;


INSERT INTO clientes(first_name, last_name, birthday) VALUES('Douglas','Marra','16/04/1989');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao2','abc','21/10/2019');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao3','abc','22/10/2019');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao4','abc','23/10/2019');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao5','abc','24/10/2019');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao6','abc','25/10/2019');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao7','abc','26/10/2019');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao8','abc','27/10/2019');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao9','abc','28/10/2019');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao10','abc','29/10/2019');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao11','abc','30/10/2019');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao12','abc','31/10/2019');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao13','abc','01/11/2019');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao14','abc','02/11/2019');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao15','abc','03/11/2019');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao16','abc','04/11/2019');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao17','abc','05/11/2019');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao18','abc','06/11/2019');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao19','abc','07/11/2019');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao20','abc','08/11/2019');
INSERT INTO clientes(first_name, last_name, birthday) VALUES('Joao21','abc','09/11/2019');
commit;


INSERT INTO campanhas(campanha, status, pct) VALUES('Blackfriday',false, 10);
INSERT INTO campanhas(campanha, status, pct) VALUES('Aniversario',true, 5);
commit;


select * from campanhas;

INSERT INTO produtos(title, description, priceincents) VALUES('IPHONE 11','64GB - green', 6500);
INSERT INTO produtos(title, description, priceincents) VALUES('MOTOROLA','64GB - green', 900);
INSERT INTO produtos(title, description, priceincents) VALUES('Notebook DELL','NOtebook Dell ABC', 4500);
commit;

Select * from campanhas;
select * from clientes;
select * from produtos;
