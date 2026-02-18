package com.example.demo.exception.rickandmorty;

import com.example.demo.exception.ResourceNotFoundException;

public class CharacterNotFoundException extends ResourceNotFoundException {

    public CharacterNotFoundException(Long id) {
        super("Character with ID " + id + " not found");
    }

    public CharacterNotFoundException(Long id, Throwable cause) {
        super("Character with ID " + id + " not found", cause);
    }
}
