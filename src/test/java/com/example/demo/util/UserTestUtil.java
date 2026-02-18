package com.example.demo.util;

import com.example.demo.dto.users.UserRequest;

public class UserTestUtil {

    private UserTestUtil() {
        // Private constructor for utility class
    }

    public static UserRequest createDefaultUserRequest() {
        return UserRequest.builder()
                .name("Integration User")
                .email("integration@test.com")
                .phone("1234567890")
                .website("integration.com")
                .build();
    }

    public static UserRequest createSecondaryUserRequest() {
        return UserRequest.builder()
                .name("Integration User 2")
                .email("integration2@test.com")
                .phone("0987654321")
                .website("integration2.com")
                .build();
    }
}
