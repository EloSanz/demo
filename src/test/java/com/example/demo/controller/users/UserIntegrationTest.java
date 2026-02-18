package com.example.demo.controller.users;

import static org.assertj.core.api.Assertions.assertThat;

import com.example.demo.AbstractIntegrationTest;
import com.example.demo.dto.users.UserRequest;
import com.example.demo.util.UserTestUtil;
import com.example.demo.dto.users.UserResponse;
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
        void givenUserRequest_whenCreateUser_shouldReturnCreatedUser() {
                // Given
                UserRequest request = UserTestUtil.createDefaultUserRequest();

                // When: Create User
                ResponseEntity<UserResponse> response = restTemplate.postForEntity("/api/users", request,
                                UserResponse.class);

                // Then: Verify Creation
                assertThat(response.getStatusCode()).isEqualTo(HttpStatus.CREATED);

                UserResponse createdUser = Optional.ofNullable(response.getBody())
                                .orElseThrow(() -> new AssertionError("Response body is null"));

                assertThat(createdUser.getId()).isNotNull();
                assertThat(createdUser.getEmail()).isEqualTo("integration@test.com");
        }

        @Test
        void givenExistingUserId_whenGetUser_shouldReturnUser() {
                // Given
                UserRequest request = UserTestUtil.createSecondaryUserRequest();

                ResponseEntity<UserResponse> createResponse = restTemplate.postForEntity("/api/users", request,
                                UserResponse.class);
                assertThat(createResponse.getStatusCode()).isEqualTo(HttpStatus.CREATED);

                // Safely get user ID from response body using Optional
                Long userId = Optional.ofNullable(createResponse.getBody())
                                .map(UserResponse::getId)
                                .orElseThrow(
                                                () -> new AssertionError("Failed to create user or retrieve ID"));

                // When: Get User
                ResponseEntity<UserResponse> getResponse = restTemplate.getForEntity("/api/users/" + userId,
                                UserResponse.class);

                // Then
                assertThat(getResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
                UserResponse user = Optional.ofNullable(getResponse.getBody())
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
