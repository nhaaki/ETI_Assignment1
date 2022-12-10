CREATE USER 'user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL ON *.* TO 'user'@'localhost';

CREATE DATABASE DriveUserDB;
USE DriveUserDB;
CREATE TABLE Passengers (UserID INT auto_increment, Username varchar(50), Password varchar(50), FirstName varchar(20), LastName varchar(20), MobileNo varchar(8), EmailAddress varchar(50), PRIMARY KEY (UserID));
CREATE TABLE Drivers (UserID INT auto_increment, Username varchar(50), Password varchar(50), FirstName varchar(20), LastName varchar(20), MobileNo varchar(8), EmailAddress varchar(50), IdNo varchar(20), CarLicenseNo varchar(20), PRIMARY KEY (UserID));
CREATE DATABASE DriveFunctionDB;
USE DriveFunctionDB;
CREATE TABLE LiveRides (driverUID INT, passengerUID INT, pcPickUp varchar(10), pcDropOff varchar(10), status varchar(15), primary key (DriverUID));
CREATE DATABASE DriveDataDB;
USE DriveDataDB;
CREATE TABLE RideHistory (passengerUID INT, driverUID INT, pcPickUp varchar(10), pcDropOff varchar(10), rideDate DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP);