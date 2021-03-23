CREATE TABLE loncherapp.App_Configuration (
                                              ID INT auto_increment NOT NULL,
                                              `Parameter` varchar(250) NOT NULL,
                                              Value varchar(100) NULL,
                                              CONSTRAINT NewTable_UN UNIQUE KEY (`Parameter`),
                                              CONSTRAINT NewTable_PK PRIMARY KEY (ID)
)
    ENGINE=InnoDB
    DEFAULT CHARSET=latin1
    COLLATE=latin1_swedish_ci;
