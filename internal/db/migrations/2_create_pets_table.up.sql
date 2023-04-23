CREATE TYPE SPECIAL_REQUIREMENT AS ENUM ('NO_TREATS');

CREATE TABLE IF NOT EXISTS PETS (
    ID INT GENERATED ALWAYS AS IDENTITY,
    petname VARCHAR(16),
    species VARCHAR(16),   
    age INT,
    special_requirement SPECIAL_REQUIREMENT,
    pet_owner_id INT,
    PRIMARY KEY(ID),
    CONSTRAINT fk_pets
        FOREIGN KEY(pet_owner_id) 
        REFERENCES users(ID)
);