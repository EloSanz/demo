package com.example.demo.service.impl;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import com.example.demo.domain.User;
import com.example.demo.controller.dto.users.UserRequest;
import com.example.demo.controller.dto.users.UserResponse;
import com.example.demo.exception.ResourceNotFoundException;
import com.example.demo.repository.UserRepository;
import java.util.Optional;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

@ExtendWith(MockitoExtension.class)
class UserServiceTest {

    @Mock private UserRepository userRepository;

    @InjectMocks private UserServiceImpl userService;

    @Test
    void givenUserRequest_whenCreateUser_shouldReturnUserResponse() {
        // Given
        UserRequest request =
                UserRequest.builder()
                        .name("Test User")
                        .email("test@example.com")
                        .phone("1234567890")
                        .website("test.com")
                        .build();

        User savedUser = new User();
        savedUser.setId(1L);
        savedUser.setName("Test User");
        savedUser.setEmail("test@example.com");

        when(userRepository.save(any(User.class))).thenReturn(savedUser);

        // When
        UserResponse response = userService.createUser(request);

        // Then
        assertThat(response).isNotNull();
        assertThat(response.getId()).isEqualTo(1L);
        assertThat(response.getName()).isEqualTo("Test User");
        verify(userRepository).save(any(User.class));
    }

    @Test
    void givenExistingId_whenGetUser_shouldReturnUserResponse() {
        // Given
        Long id = 1L;
        User user = new User();
        user.setId(id);
        user.setName("Test User");
        user.setEmail("test@example.com");

        when(userRepository.findById(id)).thenReturn(Optional.of(user));

        // When
        UserResponse response = userService.getUserById(id);

        // Then
        assertThat(response.getId()).isEqualTo(id);
        assertThat(response.getName()).isEqualTo("Test User");
        assertThat(response.getEmail()).isEqualTo("test@example.com");
    }

    @Test
    void givenNonExistentId_whenGetUser_shouldThrowException() {
        Long id = 99L;
        when(userRepository.findById(id)).thenReturn(Optional.empty());

        assertThatThrownBy(() -> userService.getUserById(id))
                .isInstanceOf(ResourceNotFoundException.class)
                .hasMessageContaining("User not found");
    }
}
