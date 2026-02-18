package com.example.demo.exception.users;

import com.example.demo.exception.ResourceNotFoundException;

public class UserNotFoundException extends ResourceNotFoundException {

    public UserNotFoundException(Long id) {
        super("User with ID " + id + " not found");
    }

    public UserNotFoundException(Long id, Throwable cause) {
        super("User with ID " + id + " not found", cause);
    }
}
