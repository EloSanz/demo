package com.example.demo.exception.users;

import com.example.demo.exception.DomainException;

public class UserAlreadyExistsException extends DomainException {
    public UserAlreadyExistsException(String email) {
        super("User with email " + email + " already exists");
    }
}
