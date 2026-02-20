package com.example.demo.client;

import com.example.demo.dto.users.UserResponseDto;
import java.util.List;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.service.annotation.GetExchange;
import org.springframework.web.service.annotation.HttpExchange;

/**
 * Declarative HTTP Interface for consuming external User API.
 *
 * <p>
 * This is the NEW Spring Boot 4 way of creating REST clients. No need for
 * RestTemplate or
 * WebClient boilerplate!
 *
 * <p>
 * The base URL is configured in application.yml and injected via @Bean
 * configuration.
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
