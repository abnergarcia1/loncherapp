CREATE TABLE `Menus` (
                         `ID` int(11) NOT NULL AUTO_INCREMENT,
                         `Lonchera_ID` int(11) DEFAULT NULL,
                         `Name` varchar(100) DEFAULT NULL,
                         `Description` varchar(250) DEFAULT NULL,
                         `Price` decimal(10,2) DEFAULT NULL,
                         `Currency` varchar(3) DEFAULT NULL,
                         `Image_URL` varchar(250) DEFAULT NULL,
                         `Created_At` datetime DEFAULT NULL,
                         `Updated_At` datetime DEFAULT NULL,
                         PRIMARY KEY (`ID`),
                         UNIQUE KEY `Menus_UN` (`Name`,`Lonchera_ID`,`Description`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=latin1