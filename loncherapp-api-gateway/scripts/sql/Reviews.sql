CREATE TABLE `Reviews` (
                           `ID` int(11) NOT NULL AUTO_INCREMENT,
                           `Lonchera_ID` int(11) DEFAULT NULL,
                           `Comment` varchar(250) DEFAULT NULL,
                           `User_ID` int(11) DEFAULT NULL,
                           `User_Name` varchar(50) DEFAULT NULL,
                           `Rating` int(11) DEFAULT NULL,
                           `Created_At` datetime DEFAULT NULL,
                           PRIMARY KEY (`ID`),
                           UNIQUE KEY `Reviews_UN` (`Created_At`,`User_ID`,`Lonchera_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1