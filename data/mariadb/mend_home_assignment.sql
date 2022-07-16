CREATE DATABASE IF NOT EXISTS `mend_home_assignment`;
USE `mend_home_assignment`;

SET names utf8;

CREATE TABLE `books` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `author` varchar(255) NOT NULL DEFAULT '',
  `title` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;

insert into `books` (`id`, `author`, `title`) values
  (1, "Fedor Dostoevsky", "Crime and Punishment"),
  (2, "Ernest Hamingway", "For Whom the Bell Tolls"),
  (3, "Sholem Aleichem", "Wandering Stars"),
  (4, "Danilo Kiš", "A Tomb for Boris Davidovich");

CREATE TABLE `cars` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `brand` varchar(255) NOT NULL DEFAULT '',
  `model` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;

insert into `cars` (`id`, `brand`, `model`) values
  (1, "Mitsubishi", "Outlander"),
  (2, "Toyota", "Camry");
