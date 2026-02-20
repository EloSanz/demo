package com.example.demo.service.impl;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import com.example.demo.domain.User;
import com.example.demo.entity.UserEntity;
import com.example.demo.exception.users.UserNotFoundException;
import com.example.demo.mapper.ExternalUserMapper;
import com.example.demo.mapper.UserEntityMapper;
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
    @Mock private UserEntityMapper entityMapper;
    @Mock private ExternalUserMapper externalMapper;

    @InjectMocks private UserServiceImpl userService;

    @Test
    void givenUserDomain_whenCreateUser_shouldReturnUserDomain() {
        // Given
        User domainUser = new User();
        domainUser.setName("Test User");
        domainUser.setEmail("test@example.com");

        UserEntity userEntity = new UserEntity();
        userEntity.setName("Test User");
        userEntity.setEmail("test@example.com");

        UserEntity savedEntity = new UserEntity();
        savedEntity.setId(1L);
        savedEntity.setName("Test User");
        savedEntity.setEmail("test@example.com");

        User savedDomain = new User();
        savedDomain.setId(1L);
        savedDomain.setName("Test User");
        savedDomain.setEmail("test@example.com");

        when(entityMapper.toEntity(domainUser)).thenReturn(userEntity);
        when(userRepository.save(any(UserEntity.class))).thenReturn(savedEntity);
        when(entityMapper.toDomain(savedEntity)).thenReturn(savedDomain);

        // When
        User response = userService.createUser(domainUser);

        // Then
        assertThat(response).isNotNull();
        assertThat(response.getId()).isEqualTo(1L);
        assertThat(response.getName()).isEqualTo("Test User");
        verify(userRepository).save(any(UserEntity.class));
    }

    @Test
    void givenExistingId_whenGetUser_shouldReturnUserDomain() {
        // Given
        Long id = 1L;
        UserEntity userEntity = new UserEntity();
        userEntity.setId(id);
        userEntity.setName("Test User");
        userEntity.setEmail("test@example.com");

        User domainUser = new User();
        domainUser.setId(id);
        domainUser.setName("Test User");
        domainUser.setEmail("test@example.com");

        when(userRepository.findById(id)).thenReturn(Optional.of(userEntity));
        when(entityMapper.toDomain(userEntity)).thenReturn(domainUser);

        // When
        User response = userService.getUserById(id);

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
                .isInstanceOf(UserNotFoundException.class)
                .hasMessageEndingWith("not found");
    }
}
