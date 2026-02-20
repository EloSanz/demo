package com.example.demo.service;

import com.example.demo.domain.User;
import java.util.List;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;

/** Service interface for User operations. Defines the contract for business logic. */
public interface UserService {

    /** Get all users from local database with pagination */
    Page<User> getAllUsers(Pageable pageable);

    /** Get a user by ID from local database */
    User getUserById(Long id);

    /** Create a new user in local database */
    User createUser(User request);

    /** Update an existing user in local database */
    User updateUser(Long id, User request);

    /** Delete a user from local database */
    void deleteUser(Long id);

    /** Fetch users from external API and optionally sync to local DB */
    List<User> fetchUsersFromExternalApi();

    /** Sync a specific user from external API to local database */
    User syncUserFromExternalApi(Long externalUserId);
}
