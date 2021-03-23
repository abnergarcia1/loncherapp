CREATE TABLE `Routes_Schedule` (
                                   `ID` int(11) NOT NULL AUTO_INCREMENT,
                                   `Route_ID` int(11) DEFAULT NULL,
                                   `Weekday` int(11) DEFAULT NULL,
                                   `Arrive_At` datetime DEFAULT NULL,
                                   `Gone_At` datetime DEFAULT NULL,
                                   `Created_At` datetime DEFAULT NULL,
                                   `Active` tinyint(1) DEFAULT NULL,
                                   PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1