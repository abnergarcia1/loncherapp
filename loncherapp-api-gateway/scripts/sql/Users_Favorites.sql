CREATE TABLE `Users_Favorites` (
                                   `ID` int(11) NOT NULL AUTO_INCREMENT,
                                   `User_ID` int(11) DEFAULT NULL,
                                   `Lonchera_ID` int(11) DEFAULT NULL,
                                   `Created_At` datetime DEFAULT NULL,
                                   PRIMARY KEY (`ID`),
                                   UNIQUE KEY `Users_Favorites_UN` (`User_ID`,`Lonchera_ID`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=latin1