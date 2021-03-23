CREATE TABLE `App_Payments` (
                                `ID` varchar(250) NOT NULL,
                                `Lonchera_ID` int(11) DEFAULT NULL,
                                `Amount` decimal(7,2) DEFAULT NULL,
                                `Currency` varchar(100) DEFAULT NULL,
                                `Status` varchar(100) NOT NULL,
                                `Created_At` datetime DEFAULT NULL,
                                `Type` varchar(100) DEFAULT NULL,
                                PRIMARY KEY (`ID`,`Status`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1