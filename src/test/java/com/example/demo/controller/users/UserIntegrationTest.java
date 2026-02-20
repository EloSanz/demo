package com.example.demo.controller.users;

import static org.assertj.core.api.Assertions.assertThat;

import com.example.demo.BaseIntegrationTest;
import com.example.demo.dto.users.UserRequestDto;
import com.example.demo.dto.users.UserResponseDto;
import com.example.demo.util.UserTestUtil;
import java.util.Map;
import java.util.Optional;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;

class UserIntegrationTest extends BaseIntegrationTest {

    @Autowired private TestRestTemplate restTemplate;

    @Test
    void givenUserRequestDto_whenCreateUser_shouldReturnCreatedUser() {
        // Given
        UserRequestDto request = UserTestUtil.createDefaultUserRequest();

        // When: Create User
        ResponseEntity<UserResponseDto> response =
                restTemplate.postForEntity("/api/users", request, UserResponseDto.class);

        // Then: Verify Creation
        assertThat(response.getStatusCode()).isEqualTo(HttpStatus.CREATED);

        UserResponseDto createdUser =
                Optional.ofNullable(response.getBody())
                        .orElseThrow(() -> new AssertionError("Response body is null"));

        assertThat(createdUser.getId()).isNotNull();
        assertThat(createdUser.getEmail()).isEqualTo("integration@test.com");
    }

    @Test
    void givenExistingUserId_whenGetUser_shouldReturnUser() {
        // Given
        UserRequestDto request = UserTestUtil.createSecondaryUserRequest();

        ResponseEntity<UserResponseDto> createResponse =
                restTemplate.postForEntity("/api/users", request, UserResponseDto.class);
        assertThat(createResponse.getStatusCode()).isEqualTo(HttpStatus.CREATED);

        // Safely get user ID from response body using Optional
        Long userId =
                Optional.ofNullable(createResponse.getBody())
                        .map(UserResponseDto::getId)
                        .orElseThrow(
                                () -> new AssertionError("Failed to create user or retrieve ID"));

        // When: Get User
        ResponseEntity<UserResponseDto> getResponse =
                restTemplate.getForEntity("/api/users/" + userId, UserResponseDto.class);

        // Then
        assertThat(getResponse.getStatusCode()).isEqualTo(HttpStatus.OK);
        UserResponseDto user =
                Optional.ofNullable(getResponse.getBody())
                        .orElseThrow(() -> new AssertionError("Response body is null"));
        assertThat(user.getName()).isEqualTo("Integration User 2");
    }

    @Test
    void givenNonExistentUserId_whenGetUser_shouldReturnNotFound() {
        // Given
        long nonExistentId = 9999L;

        // When & Then
        ResponseEntity<String> response =
                restTemplate.getForEntity("/api/users/" + nonExistentId, String.class);

        assertThat(response.getStatusCode()).isEqualTo(HttpStatus.NOT_FOUND);
    }

    @Test
    @DisplayName("Should control pagination via URL parameters")
    void testPaginationFromRouter() {
        // 1. Create 3 users
        createTestUser("User 1", "user1@test.com");
        createTestUser("User 2", "user2@test.com");
        createTestUser("User 3", "user3@test.com");

        // 2. Fetch first page with size 2
        ResponseEntity<Map<String, Object>> responsePage0 =
                restTemplate.exchange(
                        "/api/users?page=0&size=2&sort=name,asc",
                        HttpMethod.GET,
                        null,
                        new ParameterizedTypeReference<Map<String, Object>>() {});

        assertThat(responsePage0.getStatusCode()).isEqualTo(HttpStatus.OK);
        assertThat(responsePage0.getBody()).isNotNull();
        assertThat(responsePage0.getBody().get("numberOfElements")).isEqualTo(2);
        assertThat(responsePage0.getBody().get("totalElements")).isEqualTo(3);

        // 3. Fetch second page with size 2
        ResponseEntity<Map<String, Object>> responsePage1 =
                restTemplate.exchange(
                        "/api/users?page=1&size=2&sort=name,asc",
                        HttpMethod.GET,
                        null,
                        new ParameterizedTypeReference<Map<String, Object>>() {});

        assertThat(responsePage1.getStatusCode()).isEqualTo(HttpStatus.OK);
        assertThat(responsePage1.getBody()).isNotNull();
        assertThat(responsePage1.getBody().get("numberOfElements")).isEqualTo(1);
    }

    private void createTestUser(String name, String email) {
        UserRequestDto request =
                UserRequestDto.builder()
                        .name(name)
                        .email(email)
                        .phone("123")
                        .website("test.com")
                        .build();
        restTemplate.postForEntity("/api/users", request, UserResponseDto.class);
    }
}
