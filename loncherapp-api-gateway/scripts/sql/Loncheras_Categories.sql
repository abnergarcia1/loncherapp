CREATE TABLE `Loncheras_Categories` (
                                        `ID` int(11) NOT NULL AUTO_INCREMENT,
                                        `Name` varchar(100) DEFAULT NULL,
                                        `Description` varchar(250) DEFAULT NULL,
                                        `Created_At` datetime DEFAULT NULL,
                                        `Updated_At` datetime DEFAULT NULL,
                                        PRIMARY KEY (`ID`),
                                        UNIQUE KEY `Loncheras_Categories_UN` (`Name`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1