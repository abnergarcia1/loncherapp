-- loncherapp.Users definition


CREATE TABLE `Users` (
                         `ID` int(10) unsigned NOT NULL AUTO_INCREMENT,
                         `Type_ID` int(10) unsigned NOT NULL,
                         `FirstName` varchar(50) DEFAULT NULL,
                         `LastName` varchar(50) DEFAULT NULL,
                         `Email` varchar(100) DEFAULT NULL,
                         `CreationDate` datetime DEFAULT NULL,
                         `UpdatedDate` datetime DEFAULT NULL,
                         `Active` tinyint(1) DEFAULT NULL,
                         `Password` varchar(250) DEFAULT NULL,
                         PRIMARY KEY (`ID`),
                         UNIQUE KEY `Users_UN` (`Email`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=latin1;
