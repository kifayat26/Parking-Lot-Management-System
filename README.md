# Parking Management System API Documentation

This document provides an overview of the endpoints available in the Parking Management System API along with their respective request bodies and time complexities.

## Base URL

The base URL for all API endpoints is `https://example.com`.

## Endpoints

### 1. Create User

- **URL:** `/pms/createUser`
- **Method:** `POST`
- **Request Body:**
  ```json
  {
    "name": "string"
  }


### 2. Create Car

- **URL**: `/pms/createCar/{userID}`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "name": "string"
  }


### 3. Park Car

- **URL**: `/parkCar`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "parking_slot_id": "number"
  }

### 4. Unpark Car

- **URL**: `/unparkCar`
- **Method**: `POST`
- **Request Body**: 
  ```json
  {
    "parking_slot_id": "number"
  }


### 5. Create Parking Lot

- **URL**: `/createParking`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "location": "string",
    "slots": "number"
  }


### 6. Put Parking Slot in Maintenance

- **URL**: `/parking-slots/maintenance`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "parking_slot_id": "number"
  }


### 7. Put Parking Slot Out of Maintenance

- **URL**: `/parking-slots/out-of-maintenance`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "parking_slot_id": "number"
  }


### 8. Get Parking Lot Status

- **URL**: `/parking-lot/status`
- **Method**: `GET`
- **Request Body**: 
  ```json
  {
    "parking_lot_id": "number"
  }


### 9. Get History for Day

- **URL**: `/history`
- **Method**: `GET`
- **Request Body**:
  ```json
  {
    "date": "string (YYYY-MM-DD)"
  }
