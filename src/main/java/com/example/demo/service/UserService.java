package com.example.demo.service;

import com.example.demo.controller.dto.users.UserRequest;
import com.example.demo.controller.dto.users.UserResponse;
import java.util.List;

/** Service interface for User operations. Defines the contract for business logic. */
public interface UserService {

    /** Get all users from local database */
    List<UserResponse> getAllUsers();

    /** Get a user by ID from local database */
    UserResponse getUserById(Long id);

    /** Create a new user in local database */
    UserResponse createUser(UserRequest request);

    /** Update an existing user in local database */
    UserResponse updateUser(Long id, UserRequest request);

    /** Delete a user from local database */
    void deleteUser(Long id);

    /** Fetch users from external API and optionally sync to local DB */
    List<UserResponse> fetchUsersFromExternalApi();

    /** Sync a specific user from external API to local database */
    UserResponse syncUserFromExternalApi(Long externalUserId);
}
