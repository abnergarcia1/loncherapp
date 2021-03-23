CREATE TABLE `Loncheras` (
                             `ID` int(11) NOT NULL AUTO_INCREMENT,
                             `User_ID` int(11) DEFAULT NULL,
                             `Description` varchar(250) DEFAULT NULL,
                             `Category_ID` int(11) DEFAULT NULL,
                             `Cover_Image_URL` varchar(250) DEFAULT NULL,
                             `Website` varchar(250) DEFAULT NULL,
                             `Active` tinyint(1) DEFAULT NULL,
                             `Membership_Due_Date` datetime DEFAULT NULL,
                             `Created_At` datetime DEFAULT NULL,
                             `Updated_At` datetime DEFAULT NULL,
                             PRIMARY KEY (`ID`),
                             UNIQUE KEY `Loncheras_UN` (`User_ID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1