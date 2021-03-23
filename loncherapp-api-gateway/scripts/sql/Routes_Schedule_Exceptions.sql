CREATE TABLE `Routes_Schedule_Exceptions` (
                                              `ID` int(11) NOT NULL AUTO_INCREMENT,
                                              `Route_Schedule_ID` int(11) DEFAULT NULL,
                                              `Arrive_At` time DEFAULT NULL,
                                              `Gone_At` time DEFAULT NULL,
                                              `Discard_Arrival` tinyint(1) DEFAULT NULL,
                                              `Apply_Date` date DEFAULT NULL,
                                              PRIMARY KEY (`ID`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1