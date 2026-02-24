package com.example.demo.exception.storage;

import com.example.demo.exception.DomainException;

public class StorageException extends DomainException {
    public StorageException(String message) {
        super(message);
    }

    public StorageException(String message, Throwable cause) {
        super(message, cause);
    }
}
