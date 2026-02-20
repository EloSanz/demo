package com.example.demo.service;

import static org.assertj.core.api.Assertions.assertThat;

import com.example.demo.BaseIntegrationTest;
import com.example.demo.domain.User;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;

class UserServiceIntegrationTest extends BaseIntegrationTest {

    @Autowired private UserService userService;

    @Test
    @DisplayName("Should save and retrieve a user with pagination")
    void testCreateAndFindAll() {
        // Arrange
        User user =
                User.builder()
                        .name("Integration User")
                        .email("test-it@example.com")
                        .phone("123456")
                        .website("integration.com")
                        .build();

        // Act
        userService.createUser(user);
        Page<User> result = userService.getAllUsers(PageRequest.of(0, 10));

        // Assert
        assertThat(result.getContent()).hasSize(1);
        assertThat(result.getContent().getFirst().getName()).isEqualTo("Integration User");
        assertThat(result.getContent().getFirst().getEmail()).isEqualTo("test-it@example.com");
    }
}
