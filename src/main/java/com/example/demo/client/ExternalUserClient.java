package com.example.demo.client;

import com.example.demo.dto.users.UserResponseDto;
import java.util.List;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.service.annotation.GetExchange;
import org.springframework.web.service.annotation.HttpExchange;

/**
 * Declarative HTTP Interface for consuming external User API. Leverages Spring Boot 4 declarative
 * REST clients to reduce boilerplate. The base URL is configured in application.yml and injected
 * via configuration.
 */
@HttpExchange
public interface ExternalUserClient {

    /** Fetch all users from external API */
    @GetExchange("/users")
    List<UserResponseDto> getAllUsers();

    /** Fetch a single user by ID from external API */
    @GetExchange("/users/{id}")
    UserResponseDto getUserById(@PathVariable Long id);
}
