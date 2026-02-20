package com.example.demo.util;

import com.example.demo.dto.users.UserRequestDto;

public class UserTestUtil {

    private UserTestUtil() {
        // Private constructor for utility class
    }

    public static UserRequestDto createDefaultUserRequest() {
        return UserRequestDto.builder()
                .name("Integration User")
                .email("integration@test.com")
                .phone("1234567890")
                .website("integration.com")
                .build();
    }

    public static UserRequestDto createSecondaryUserRequest() {
        return UserRequestDto.builder()
                .name("Integration User 2")
                .email("integration2@test.com")
                .phone("0987654321")
                .website("integration2.com")
                .build();
    }
}
