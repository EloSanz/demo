package com.example.demo.controller.users;

import com.example.demo.service.UserService;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * Base controller for all User-related endpoints. Centralizes the base path mapping and common
 * dependencies.
 */
@RestController
@RequestMapping("/api/users")
@RequiredArgsConstructor
@Tag(name = "Users", description = "Operations related to User management")
public abstract class BaseUserController {

    protected final UserService userService;
}
