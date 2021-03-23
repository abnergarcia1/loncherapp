CREATE TABLE `Routes` (
  `ID` int(11) NOT NULL AUTO_INCREMENT,
  `Lonchera_ID` int(11) DEFAULT NULL,
  `Location` varchar(250) DEFAULT NULL,
  `Address` varchar(250) DEFAULT NULL,
  `Name` varchar(100) DEFAULT NULL,
  `Description` varchar(250) DEFAULT NULL,
  `Order` int(11) DEFAULT NULL,
  `Created_At` datetime DEFAULT NULL,
  `Updated_At` datetime DEFAULT NULL,
  `Latitude` decimal(11,7) DEFAULT NULL,
  `Longitude` decimal(11,7) DEFAULT NULL,
  `Google_Place_ID` varchar(250) DEFAULT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `Routes_UN` (`Lonchera_ID`,`Location`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=latin1