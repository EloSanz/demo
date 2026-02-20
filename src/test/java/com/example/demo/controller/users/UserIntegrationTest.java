package com.example.demo.controller.users;

import static org.assertj.core.api.Assertions.assertThat;

import com.example.demo.AbstractIntegrationTest;
import com.example.demo.dto.users.UserRequestDto;
import com.example.demo.dto.users.UserResponseDto;
import com.example.demo.util.UserTestUtil;
import java.util.Optional;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;

class UserIntegrationTest extends AbstractIntegrationTest {

        @Autowired
        private TestRestTemplate restTemplate;

        @Test
        void givenUserRequestDto_whenCreateUser_shouldReturnCreatedUser() {
                // Given
                UserRequestDto request = UserTestUtil.createDefaultUserRequest();

                // When: Create User
                ResponseEntity<UserResponseDto> response = restTemplate.postForEntity("/api/users", request,
                                UserResponseDto.class);

                // Then: Verify Creation
                assertThat(response.getStatusCode()).isEqualTo(HttpStatus.CREATED);

                UserResponseDto createdUser = Optional.ofNullable(response.getBody())
                                .orElseThrow(() -> new AssertionError("Response body is null"));

                assertThat(createdUser.getId()).isNotNull();
                assertThat(createdUser.getEmail()).isEqualTo("integration@test.com");
        }

        @Test
        void givenExistingUserId_whenGetUser_shouldReturnUser() {
                // Given
                UserRequestDto request = UserTestUtil.createSecondaryUserRequest();

                ResponseEntity<UserResponseDto> createResponse = restTemplate.postForEntity("/api/users", request,
                                UserResponseDto.class);
                assertThat(createResponse.getStatusCode()).isEqualTo(HttpStatus.CREATED);

                // Safely get user ID from response body using Optional
                Long userId = Optional.ofNullable(createResponse.getBody())
                                .map(UserResponseDto::getId)
                                .orElseThrow(
                                                () -> new AssertionError("Failed to create user or retrieve ID"));

                // When: Get User
                ResponseEntity<UserResponseDto> getResponse = restTemplate.getForEntity("/api/users/" + userId,
                                UserResponseDto.class);

                // Then
                assertThat(getResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
                UserResponseDto user = Optional.ofNullable(getResponse.getBody())
                                .orElseThrow(() -> new AssertionError("Response body is null"));
                assertThat(user.getName()).isEqualTo("Integration User 2");
        }

        @Test
        void givenNonExistentUserId_whenGetUser_shouldReturnNotFound() {
                // Given
                long nonExistentId = 9999L;

                // When & Then
                ResponseEntity<String> response = restTemplate.getForEntity("/api/users/" + nonExistentId,
                                String.class);

                assertThat(response.getStatusCode()).isEqualTo(HttpStatus.NOT_FOUND);
        }
}
