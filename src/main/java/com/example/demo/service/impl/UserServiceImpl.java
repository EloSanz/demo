package com.example.demo.service.impl;

import com.example.demo.client.ExternalUserClient;
import com.example.demo.domain.User;
import com.example.demo.dto.users.UserRequest;
import com.example.demo.dto.users.UserResponse;
import com.example.demo.exception.users.UserNotFoundException;
import com.example.demo.repository.UserRepository;
import com.example.demo.service.UserService;
import java.util.List;
import java.util.stream.Collectors;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

/**
 * Implementation of UserService. Contains business logic for User operations.
 */
@Service
@RequiredArgsConstructor
@Slf4j
@Transactional
public class UserServiceImpl implements UserService {

    private final UserRepository userRepository;
    private final ExternalUserClient externalUserClient;

    @Override
    @Transactional(readOnly = true)
    public List<UserResponse> getAllUsers() {
        log.info("Fetching all users from database");
        return userRepository.findAll().stream()
                .map(this::mapToResponse)
                .collect(Collectors.toList());
    }

    @Override
    @Transactional(readOnly = true)
    public UserResponse getUserById(Long id) {
        log.info("Fetching user with id: {}", id);
        User user = userRepository
                .findById(id)
                .orElseThrow(
                        () -> new UserNotFoundException(id));
        return mapToResponse(user);
    }

    @Override
    public UserResponse createUser(UserRequest request) {
        log.info("Creating new user with email: {}", request.getEmail());

        if (userRepository.existsByEmail(request.getEmail())) {
            throw new RuntimeException("User already exists with email: " + request.getEmail());
        }

        User user = User.builder()
                .name(request.getName())
                .email(request.getEmail())
                .phone(request.getPhone())
                .website(request.getWebsite())
                .build();

        User savedUser = userRepository.save(user);
        log.info("User created successfully with id: {}", savedUser.getId());

        return mapToResponse(savedUser);
    }

    @Override
    public UserResponse updateUser(Long id, UserRequest request) {
        log.info("Updating user with id: {}", id);

        User user = userRepository
                .findById(id)
                .orElseThrow(
                        () -> new UserNotFoundException(id));

        user.setName(request.getName());
        user.setEmail(request.getEmail());
        user.setPhone(request.getPhone());
        user.setWebsite(request.getWebsite());

        User updatedUser = userRepository.save(user);
        log.info("User updated successfully with id: {}", updatedUser.getId());

        return mapToResponse(updatedUser);
    }

    @Override
    public void deleteUser(Long id) {
        log.info("Deleting user with id: {}", id);

        if (!userRepository.existsById(id)) {
            throw new UserNotFoundException(id);
        }

        userRepository.deleteById(id);
        log.info("User deleted successfully with id: {}", id);
    }

    @Override
    @Transactional(readOnly = true)
    public List<UserResponse> fetchUsersFromExternalApi() {
        log.info("Fetching users from external API");
        return externalUserClient.getAllUsers();
    }

    @Override
    public UserResponse syncUserFromExternalApi(Long externalUserId) {
        log.info("Syncing user from external API with id: {}", externalUserId);

        // Fetch from external API
        UserResponse externalUser = externalUserClient.getUserById(externalUserId);

        // Check if already exists in local DB
        User user = userRepository
                .findByEmail(externalUser.getEmail())
                .orElse(
                        User.builder()
                                .name(externalUser.getName())
                                .email(externalUser.getEmail())
                                .phone(externalUser.getPhone())
                                .website(externalUser.getWebsite())
                                .build());

        // Update with latest data
        user.setName(externalUser.getName());
        user.setPhone(externalUser.getPhone());
        user.setWebsite(externalUser.getWebsite());

        User savedUser = userRepository.save(user);
        log.info("User synced successfully with local id: {}", savedUser.getId());

        return mapToResponse(savedUser);
    }

    /** Helper method to map User entity to UserResponse DTO */
    private UserResponse mapToResponse(User user) {
        return UserResponse.builder()
                .id(user.getId())
                .name(user.getName())
                .email(user.getEmail())
                .phone(user.getPhone())
                .website(user.getWebsite())
                .build();
    }
}
